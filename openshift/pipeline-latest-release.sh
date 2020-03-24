#!/usr/bin/env bash

MAX_SHIFT=2
STABLE_RELEASE_URL='https://raw.githubusercontent.com/openshift/tektoncd-pipeline/${version}/openshift/release/tektoncd-pipeline-${version}.yaml'

function get_version {
    local shift=${1} # 0 is latest, increase is the version before etc...
    local version=$(curl -s https://api.github.com/repos/tektoncd/pipeline/releases | python -c "import sys, json;x=json.load(sys.stdin);print(x[${shift}]['tag_name'])")
    echo $(eval echo ${STABLE_RELEASE_URL})
}

function tryurl {
    curl -s -o /dev/null -f ${1} || return 1
}

for shifted in `seq 0 ${MAX_SHIFT}`;do
    versionyaml=$(get_version ${shifted})
    if tryurl ${versionyaml};then
        export RELEASE_YAML=${versionyaml}
    fi
done
