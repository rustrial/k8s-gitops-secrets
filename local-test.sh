#!/bin/bash

export KIND=3e82e061-dae1-4b88-a760-549d3f9c161d
export KUBECONFIG=".kube-config"

kind create cluster --name $KIND --kubeconfig $KUBECONFIG --image 'kindest/node:v1.26.6'

kubectl delete deployment -n k8s-gitops-secrets-system k8s-gitops-secrets-controller-manager

./.github/install.sh "$@"

./.github/e2e-tests.sh