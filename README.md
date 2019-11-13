# Openshift Tektoncd Pipeline

This repository holds Openshift's fork of
[`tektoncd/triggers`](https://github.com/tektoncd/trigggers) with additions and
fixes needed only for the OpenShift side of things.

## List of releases

- [release-v0.1.0](https://github.com/openshift/tektoncd-triggers/tree/release-v0.1.0)

## How this repository works ?

The `master` branch holds up-to-date specific [openshift files](./openshift)
that are necessary for CI setups and maintaining it. This includes:

- Scripts to create a new release branch from `upstream`
- CI setup files
  - operator configuration (for Openshift's CI setup)
  - tests scripts
- Operator's base configurations

Each release branch holds the upstream code for that release and our
openshift's specific files.

## CI Setup

For the CI setup, two repositories are of importance:

- This repository
- [openshift/release](https://github.com/openshift/release) which
  contains the configuration of CI jobs that are run on this
  repository

All of the following is based on OpenShift’s CI operator
configs. General understanding of that mechanism is assumed in the
following documentation.

The job manifests for the CI jobs are generated automatically. The
basic configuration lives in the
[/ci-operator/config/openshift/tektoncd-pipeline](https://github.com/openshift/release/tree/master/ci-operator/config/openshift/tektoncd-pipeline) folder of the
[openshift/release](https://github.com/openshift/release) repository. These files include which version to
build against (OCP 4.0 for our recent cases), which images to build
(this includes all the images needed to run Knative and also all the
images required for running e2e tests) and which command to execute
for the CI jobs to run (more on this later).

Before we can create the ci-operator configs mentioned above, we need
to make sure there are Dockerfiles for all images that we need
(they’ll be referenced by the ci-operator config hence we need to
create them first). The [generate-dockerfiles.sh](https://github.com/openshift/tektoncd-triggers/blob/master/openshift/ci-operator/generate-dockerfiles.sh) script takes care of
creating all the Dockerfiles needed automatically. The files now need
to be committed to the branch that CI is being setup for.

The basic ci-operator configs mentioned above are generated via the
generate-release.sh file in the openshift/tektoncd-pipeline
repository. They are generated to alleviate the burden of having to
add all possible test images to the manifest manually, which is error
prone.

Once the file is generated, it must be committed to the
[openshift/release](https://github.com/openshift/release) repository, as the other manifests linked above. The
naming schema is `openshift-tektoncd-pipeline-BRANCH.yaml`, thus the
files existing already correspond to our existing releases and the
master branch itself.

After the file has been added to the folder as mentioned above, the
job manifests itself will need to be generated as is described in the
corresponding [ci-operator documentation](https://docs.google.com/document/d/1SQ_qlkcplqhe8h6ONXdgBr7YUVbs4oRSj4ISl3gpLW4/edit#heading=h.8w7nj9363nsd).

Once all of this is done (Dockerfiles committed, ci-operator config
created and job manifests generated) a PR must be opened against
[openshift/release](https://github.com/openshift/releaseopenshift/release)
to include all the ci-operator related files. Once
this PR is merged, the CI setup for that branch is active.

## Create a new release

### Deliverables:

- Tagged images on quay.io
- An OLM manifest referencing these images
- Install documentation

### High-level steps for a release

#### [Building upstream](RELEASE_PROCESS.md)

