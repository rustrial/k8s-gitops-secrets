FROM nexxiot-registry-oci.jfrog.io/golang:1.24.5-alpine3.22

RUN apk --no-cache add git bash make zip kubectl docker fuse-overlayfs helm curl npm # kind

# Kind from package manager is too old, so download it from upstream.
RUN ([ $(uname -m) = x86_64 ] && curl -Lo /usr/bin/kind https://kind.sigs.k8s.io/dl/v0.29.0/kind-linux-amd64) || true

RUN ([ $(uname -m) = aarch64 ] && curl -Lo /usr/bin/kind https://kind.sigs.k8s.io/dl/v0.29.0/kind-linux-arm64) || true

RUN chmod +x /usr/bin/kind

RUN go install golang.org/x/tools/gopls@v0.19.1
