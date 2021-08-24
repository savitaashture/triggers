#!/usr/bin/env bash

MAX_SHIFT=2
STABLE_RELEASE_URL='https://raw.githubusercontent.com/openshift/tektoncd-pipeline/${version}/openshift/release/tektoncd-pipeline-${version}.yaml'

function get_version {
    local shift=${1} # 0 is latest, increase is the version before etc...
    local version=$(curl -s https://api.github.com/repos/tektoncd/pipeline/releases | python -c "from pkg_resources import parse_version;import sys, json;jeez=json.load(sys.stdin);print(sorted([x['tag_name'] for x in jeez], key=parse_version, reverse=True)[${shift}])")
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
