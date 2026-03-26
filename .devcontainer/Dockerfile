FROM nexxiot-registry-oci.jfrog.io/golang:1.26.1-alpine3.23

RUN apk --no-cache add git bash make zip kubectl docker fuse-overlayfs helm curl npm # kind

# Kind from package manager is too old, so download it from upstream.
# https://github.com/kubernetes-sigs/kind/releases
#
RUN ([ $(uname -m) = x86_64 ] && curl -Lo /usr/bin/kind https://kind.sigs.k8s.io/dl/v0.31.0/kind-linux-amd64) || true

RUN ([ $(uname -m) = aarch64 ] && curl -Lo /usr/bin/kind https://kind.sigs.k8s.io/dl/v0.31.0/kind-linux-arm64) || true

RUN chmod +x /usr/bin/kind

# https://github.com/golang/tools/releases
RUN go install golang.org/x/tools/gopls@v0.21.1
