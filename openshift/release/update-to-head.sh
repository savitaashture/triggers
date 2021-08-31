#!/usr/bin/env bash

# Synchs the release-next branch to master and then triggers CI
# Usage: update-to-head.sh

set -e
REPO_NAME=$(basename $(git rev-parse --show-toplevel))
[[ ${REPO_NAME} != tektoncd-* ]] && REPO_NAME=tektoncd-${REPO_NAME}
TODAY=`date "+%Y%m%d"`

# Reset release-next to upstream/main.
git fetch upstream main
git checkout upstream/main --no-track -B release-next

# Update openshift's master and take all needed files from there.
git fetch openshift master
git checkout openshift/master openshift Makefile OWNERS_ALIASES OWNERS
make generate-dockerfiles

git add openshift OWNERS_ALIASES OWNERS Makefile
git commit -m ":open_file_folder: Update openshift specific files."

if [[ -d openshift/patches ]];then
    for f in openshift/patches/*.patch;do
        [[ -f ${f} ]] || continue
        git am ${f}
    done
fi

# add release.yaml from previous successful nightly build to re-synced release-next as a backup
git fetch openshift release-next
git checkout FETCH_HEAD openshift/release/tektoncd-triggers-nightly.yaml

git add openshift/release/tektoncd-triggers-nightly.yaml
git commit -m ":robot: Add previous days release.yaml as back up"

git push -f openshift release-next

# Trigger CI
git checkout release-next -B release-next-ci

./openshift/release/generate-release.sh nightly

date > ci
git add ci openshift/release/tektoncd-triggers-nightly.yaml
git commit -m ":robot: Triggering CI on branch 'release-next' after synching to upstream/master"

git push -f openshift release-next-ci

if hash hub 2>/dev/null; then
   # Test if there is already a sync PR in 
   COUNT=$(hub api -H "Accept: application/vnd.github.v3+json" repos/openshift/${REPO_NAME}/pulls --flat \
    | grep -c ":robot: Triggering CI on branch 'release-next' after synching to upstream/[master|main]") || true
   if [ "$COUNT" = "0" ]; then
      hub pull-request --no-edit -l "kind/sync-fork-to-upstream" -b openshift/${REPO_NAME}:release-next -h openshift/${REPO_NAME}:release-next-ci
   fi
else
   echo "hub (https://github.com/github/hub) is not installed, so you'll need to create a PR manually."
fi
