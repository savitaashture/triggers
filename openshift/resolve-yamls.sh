#!/usr/bin/env bash

# Take all yamls in $dir and generate a all in one yaml file with resolved registry image/tag
function resolve_resources() {
  local dir=$1
  local resolved_file_name=$2
  local ignores=$3
  local registry_prefix=$4
  local image_tag=${5}

  # This would get only one set of truth from the Makefile for the image lists
  #
  # % grep "^CORE_IMAGES" Makefile
  # CORE_IMAGES=./cmd/controller ./cmd/webhook ./cmd/eventlistenersink ./cmd/interceptors
  # to:
  #  % grep '^CORE_IMAGES' Makefile|sed -e 's/.*=//' -e 's,./cmd/,,g'|tr -d '\n'|sed -e 's/ /|/g' -e 's/^/(/' -e 's/$/)\n/'
  # (controller|webhook|eventlistenersink/interceptors)
  local image_regexp=$(grep '^CORE_IMAGES' $(git rev-parse --show-toplevel)/Makefile| \
                           sed -e 's/.*=//' -e 's,./cmd/,,g'|tr '\n' ' '| \
                           sed -e 's/ /|/g' -e 's/^/(/' -e 's/|$/)\n/')

  >$resolved_file_name
  for yaml in $(find $dir -maxdepth 1 -name "*.yaml" | grep -vE $ignores); do
    echo "---" >>$resolved_file_name
    if [[ -n ${image_tag} ]];then
        # This is a release format the output would look like this :
        # quay.io/openshift-pipeline/tektoncd-triggers-controller:$image_tag
        sed -e "s,ko://,,g" -e "s%\(.* image:\)\(github.com\)\(.*\/\)\(.*\)%\1 ${registry_prefix}-\4:${image_tag}%" $yaml \
            -r -e "s,github.com/tektoncd/triggers/cmd/${image_regexp},${registry_prefix}-\1:${image_tag},g" \
            >>$resolved_file_name
     else
        # This is CI which get built directly to the user registry namespace i.e: $OPENSHIFT_BUILD_NAMESPACE
        # The output would look like this :
        # internal-registry:5000/usernamespace:tektoncd-triggers-controller
        #
        # note: test image are images only used for testing not for releases
        sed -e "s,ko://,,g" -e 's%\(.* image:\)\(github.com\)\(.*\/\)\(test\/\)\(.*\)%\1\2 \3\4test-\5%' $yaml \
            -e "s%\(.* image:\)\(github.com\)\(.*\/\)\(.*\)%\1 ""$registry_prefix"'\:tektoncd-triggers-\4%'  \
            -re "s,github.com/tektoncd/triggers/cmd/${image_regexp},${registry_prefix}:tektoncd-triggers-\1,g" \
            >>$resolved_file_name
    fi

    # Remove runAsUser: id and runAsGroup: id, openshift takes care of randoming them and we dont need a fixed uid for that
	  sed -i '/runAsUser: [0-9]*/d' ${resolved_file_name}

	  sed -i '/runAsGroup: [0-9]*/d' ${resolved_file_name}

    echo >>$resolved_file_name
  done
}
