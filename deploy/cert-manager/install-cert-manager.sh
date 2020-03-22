#!/bin/bash

# Is Kubernetes 1.15+
IS_MODERN_K8S_VERSION="no"



if [[ "${IS_MODERN_K8S_VERSION}" == "yes" ]]; then
    kubectl apply --validate=false -f https://github.com/jetstack/cert-manager/releases/download/v0.14.1/cert-manager.yaml
else
    kubectl apply --validate=false -f https://github.com/jetstack/cert-manager/releases/download/v0.14.1/cert-manager-legacy.yaml
fi


# Note: If you are running Kubernetes v1.15.4 or below, you will need to add the --validate=false flag
# to your kubectl apply command above else you will receive a validation error relating to
# the x-kubernetes-preserve-unknown-fields field in cert-manager’s CustomResourceDefinition resources.
# This is a benign error and occurs due to the way kubectl performs resource validation.

# Verify installation

kubectl --namespace cert-manager get pods

clear
echo "Now you can VERIFY cert-manager installation"
