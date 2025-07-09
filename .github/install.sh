#!/bin/bash

set -e
set -x

make docker-build IMG=test/secrets-controller:latest

if ! kind load docker-image test/secrets-controller:latest --name "${KIND:-kind}" -v 10; then
    docker save -o /tmp/image-archive.tar test/secrets-controller:latest
    kind load image-archive /tmp/image-archive.tar --name "${KIND:-kind}";
fi


helm upgrade k8s-gitops-secrets-controller charts/k8s-gitops-secrets-controller --install -n k8s-gitops-secrets-system --create-namespace --set fullnameOverride=k8s-gitops-secrets-controller-manager --set image.repository=test/secrets-controller --set image.tag=latest
