
# Image URL to use all building/pushing image targets
IMG ?= controller:latest
# Produce CRDs that work back to Kubernetes 1.11 (no version conversion)
CRD_OPTIONS ?= "crd:trivialVersions=true,preserveUnknownFields=false"

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

all: manager api-docs seals charts

# Run tests
ENVTEST_ASSETS_DIR=$(shell pwd)/testbin
test: generate fmt vet manifests
	mkdir -p ${ENVTEST_ASSETS_DIR}
	test -f ${ENVTEST_ASSETS_DIR}/setup-envtest.sh || curl -sSLo ${ENVTEST_ASSETS_DIR}/setup-envtest.sh https://raw.githubusercontent.com/kubernetes-sigs/controller-runtime/v0.9.5/hack/setup-envtest.sh
	source ${ENVTEST_ASSETS_DIR}/setup-envtest.sh; fetch_envtest_tools $(ENVTEST_ASSETS_DIR); setup_envtest_env $(ENVTEST_ASSETS_DIR); go test ./... -coverprofile cover.out

# Build manager binary
manager: generate fmt vet
	go build -o bin/manager main.go

seals: 
	go build -o bin/seals cmd/seals/main.go

cli: 
	GOOS=darwin GOARCH=amd64 go build -o build/darwin/amd64/seals cmd/seals/main.go
	cat build/darwin/amd64/seals | gzip > build/darwin/amd64/seals.gz
	GOOS=darwin GOARCH=arm64 go build -o build/darwin/arm64/seals cmd/seals/main.go
	cat build/darwin/arm64/seals | gzip > build/darwin/arm64/seals.gz
	GOOS=windows GOARCH=amd64 go build -o build/windows/amd64/seals.exe cmd/seals/main.go
	zip build/windows/amd64/seals.zip build/windows/amd64/seals.exe
	GOOS=linux GOARCH=amd64 go build -o build/linux/amd64/seals cmd/seals/main.go
	cat build/linux/amd64/seals | gzip > build/linux/amd64/seals.gz
	GOOS=linux GOARCH=arm64 go build -o build/linux/arm64/seals cmd/seals/main.go
	cat build/linux/arm64/seals | gzip > build/linux/arm64/seals.gz

# Run against the configured Kubernetes cluster in ~/.kube/config
run: generate fmt vet manifests
	go run ./main.go

# Install CRDs into a cluster
install: manifests kustomize
	$(KUSTOMIZE) build config/crd | kubectl apply -f -

# Uninstall CRDs from a cluster
uninstall: manifests kustomize
	$(KUSTOMIZE) build config/crd | kubectl delete -f -

# Deploy controller in the configured Kubernetes cluster in ~/.kube/config
deploy: manifests kustomize
	cd config/manager && $(KUSTOMIZE) edit set image controller=${IMG}
	$(KUSTOMIZE) build config/default | kubectl apply -f -

# UnDeploy controller from the configured Kubernetes cluster in ~/.kube/config
undeploy:
	$(KUSTOMIZE) build config/default | kubectl delete -f -

# Generate manifests e.g. CRD, RBAC etc.
manifests: controller-gen
	$(CONTROLLER_GEN) $(CRD_OPTIONS) rbac:roleName=manager-role webhook paths="./..." output:crd:artifacts:config=config/crd/bases

# Generate API reference documentation
api-docs: gen-crd-api-reference-docs
	$(API_REF_GEN) -api-dir=./apis/secrets/v1beta1 -config=./hack/api-docs/config.json -template-dir=./hack/api-docs/template -out-file=./docs/api/SealedSecret.md

# Run go fmt against code
fmt:
	go fmt ./...

# Run go vet against code
vet:
	go vet ./...

# Generate code
generate: controller-gen
	$(CONTROLLER_GEN) object:headerFile="hack/boilerplate.go.txt" paths="./..."

# Build the docker image
docker-build: test
	docker build -t ${IMG} .

# Push the docker image
docker-push:
	docker push ${IMG}

# Download controller-gen locally if necessary
controller-gen:
ifeq (, $(shell which controller-gen))
	@{ \
	set -e ;\
	CONTROLLER_GEN_TMP_DIR=$$(mktemp -d) ;\
	cd $$CONTROLLER_GEN_TMP_DIR ;\
	go mod init tmp ;\
	go get sigs.k8s.io/controller-tools/cmd/controller-gen@v0.6.1 ;\
	rm -rf $$CONTROLLER_GEN_TMP_DIR ;\
	}
CONTROLLER_GEN=$(GOBIN)/controller-gen
else
CONTROLLER_GEN=$(shell which controller-gen)
endif




# Download kustomize locally if necessary
KUSTOMIZE = $(shell pwd)/bin/kustomize
kustomize:
	$(call go-get-tool,$(KUSTOMIZE),sigs.k8s.io/kustomize/kustomize/v3@v3.8.7)

charts: manifests
	cp config/crd/bases/* charts/k8s-gitops-secrets-controller/crds

# go-get-tool will 'go get' any package $2 and install it to $1.
PROJECT_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))
define go-get-tool
@[ -f $(1) ] || { \
set -e ;\
TMP_DIR=$$(mktemp -d) ;\
cd $$TMP_DIR ;\
go mod init tmp ;\
echo "Downloading $(2)" ;\
GOBIN=$(PROJECT_DIR)/bin go get $(2) ;\
rm -rf $$TMP_DIR ;\
}
endef

# Find or download gen-crd-api-reference-docs
gen-crd-api-reference-docs:
ifeq (, $(shell which gen-crd-api-reference-docs))
	@{ \
	set -e ;\
	API_REF_GEN_TMP_DIR=$$(mktemp -d) ;\
	cd $$API_REF_GEN_TMP_DIR ;\
	go mod init tmp ;\
	go get github.com/ahmetb/gen-crd-api-reference-docs@v0.2.0 ;\
	rm -rf $$API_REF_GEN_TMP_DIR ;\
	}
API_REF_GEN=$(GOBIN)/gen-crd-api-reference-docs
else
API_REF_GEN=$(shell which gen-crd-api-reference-docs)
endif
