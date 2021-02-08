#!/bin/bash

set -e

kubectl apply -f config/test-data/aws-kms-secret.yaml

kubectl -n k8s-gitops-secrets-system wait sealedsecrets/aws-kms-secret --for=condition=Ready --timeout=4m

[[ "$(kubectl get -n k8s-gitops-secrets-system secret/aws-kms-secret -o jsonpath='{.data.hello}')" == "d29ybGQ=" ]]

[[ "$(kubectl get -n k8s-gitops-secrets-system secret/aws-kms-secret -o jsonpath='{.data.hello2}')" == "d29ybGQ=" ]]

[[ "$(kubectl get -n k8s-gitops-secrets-system secret/aws-kms-secret -o jsonpath='{.metadata.labels.my-label}')" == "label-value" ]]

[[ "$(kubectl get -n k8s-gitops-secrets-system secret/aws-kms-secret -o jsonpath='{.metadata.annotations.my-annotation}')" == "annotation-value" ]]

