#!/bin/bash

set -e

dump_status() {
    set +e
    echo "#========================================================================================"
    echo "# Deployment:"
    echo "#========================================================================================"
    kubectl -n k8s-gitops-secrets-system get deploy/k8s-gitops-secrets-controller-manager -o yaml
    echo ""
    echo "#========================================================================================"
    echo "# SealedSecret:"
    echo "#========================================================================================"
    kubectl -n k8s-gitops-secrets-system get sealedsecrets/aws-kms-secret -o yaml
    echo ""
    echo "#========================================================================================"
    echo "# Secret:"
    echo "#========================================================================================"
    kubectl -n k8s-gitops-secrets-system get secrets/aws-kms-secret -o yaml
    echo ""
    echo "#========================================================================================"
    echo "# Logs:"
    echo "#========================================================================================"
    kubectl -n k8s-gitops-secrets-system logs deployment/k8s-gitops-secrets-controller-manager
    set -e
    false
}

kubectl apply -f config/test-data/aws-kms-secret.yaml

kubectl -n k8s-gitops-secrets-system wait deploy/k8s-gitops-secrets-controller-manager --for=condition=Available --timeout=8m || dump_status

kubectl -n k8s-gitops-secrets-system wait sealedsecrets/aws-kms-secret --for=condition=Ready --timeout=4m || dump_status

[[ "$(kubectl get -n k8s-gitops-secrets-system secret/aws-kms-secret -o jsonpath='{.data.hello}')" == "d29ybGQ=" ]]

[[ "$(kubectl get -n k8s-gitops-secrets-system secret/aws-kms-secret -o jsonpath='{.data.hello2}')" == "d29ybGQ=" ]]

[[ "$(kubectl get -n k8s-gitops-secrets-system secret/aws-kms-secret -o jsonpath='{.metadata.labels.my-label}')" == "label-value" ]]

[[ "$(kubectl get -n k8s-gitops-secrets-system secret/aws-kms-secret -o jsonpath='{.metadata.annotations.my-annotation}')" == "annotation-value" ]]
