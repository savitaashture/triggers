#This makefile is used by ci-operator

CGO_ENABLED=0
GOOS=linux
CORE_IMAGES=./cmd/controller ./cmd/eventlistenersink ./cmd/webhook ./cmd/interceptors

##
# You need to provide a RELEASE_VERSION when using targets like `push-image`, you can do it directly
# on the command like this: `make push-image RELEASE_VERSION=0.4.0`
RELEASE_VERSION=
REGISTRY_CI_URL=registry.ci.openshift.org/openshift/tektoncd-v$(RELEASE_VERSION):tektoncd-triggers
REGISTRY_RELEASE_URL=quay.io/openshift-pipeline/tektoncd-triggers

# Install core images
install: installuidwrapper
	go install $(CORE_IMAGES)
.PHONY: install

# Run E2E tests on OpenShift
test-e2e: check-images
	./openshift/e2e-tests-openshift.sh
.PHONY: test-e2e

# Make sure we have all images in the makefile variable or that would be a new
# binary that needs to be added
check-images:
	@for cmd in ./cmd/*;do \
		found="" ;\
		for image in $(CORE_IMAGES);do \
			if [[ $$image == $$cmd ]];then \
				found=1 ;\
			fi \
		done ;\
		test -n $$found || { \
			echo "*ERROR*: Could not find $$cmd in the Makefile variables CORE_IMAGES" ;\
			echo "" ;\
			echo "If it it's a new binary that was added upstream, then do the following :" ;\
			echo "- Add the binary to openshift/release like this: https://git.io/fj18c" ;\
			echo "- Add to the CORE_IMAGES variables in the Makefile" ;\
			echo "- Generate the dockerfiles by running 'make generate-dockerfiles'" ;\
			echo "- Commit and PR these to 'openshift/release-next' remote/branch and 'openshift/master'" ;\
			exit 1 ;\
		} && exit 0 ;\
	done
.PHONY: check-images

# Generate Dockerfiles used by ci-operator. The files need to be committed manually.
generate-dockerfiles:
	./openshift/ci-operator/generate-dockerfiles.sh openshift/ci-operator/Dockerfile.in openshift/ci-operator/tekton-images $(CORE_IMAGES)
.PHONY: generate-dockerfiles

# NOTE(chmou): Install uidwraper for launching some binaries with fixed uid
UIDWRAPPER_PATH=./openshift/ci-operator/uidwrapper
installuidwrapper: $(UIDWRAPPER_PATH)
	install -m755 $(UIDWRAPPER_PATH) $(GOPATH)/bin/

# Generates a ci-operator configuration for a specific branch.
generate-ci-config:
	./openshift/ci-operator/generate-ci-config.sh $(BRANCH) > ci-operator-config.yaml
.PHONY: generate-ci-config

# Generate an aggregated knative yaml file with replaced image references
generate-release:
	@test $(RELEASE_VERSION) || { echo "You need to set the RELEASE_VERSION on the command line i.e: make RELEASE_VERSION=0.4.0"; exit ;1;}
	@./openshift/release/generate-release.sh v$(RELEASE_VERSION)
.PHONY: generate-release

.PHONY: push-image
