#!/usr/bin/env bash

NIGHTLY_RELEASE='https://raw.githubusercontent.com/openshift/tektoncd-pipeline/release-next/openshift/release/tektoncd-pipeline-nightly.yaml'

function tryurl {
    curl -s -o /dev/null -f ${1} || return 1
}

# check for pipeline nightly release.
if tryurl ${NIGHTLY_RELEASE};then
    export RELEASE_YAML=${NIGHTLY_RELEASE}
fi
