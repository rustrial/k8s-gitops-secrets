FROM nexxiot-registry-oci.jfrog.io/golang:1.24.2-alpine3.21

RUN apk --no-cache add git bash make zip

RUN go install golang.org/x/tools/gopls@v0.19.1
