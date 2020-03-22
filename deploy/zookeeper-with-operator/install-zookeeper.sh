#!/bin/bash

# Namespace to install Zookeeper
NAMESPACE="${NAMESPACE:-default}"

CUR_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)"

echo "Install Zookeeper into ${NAMESPACE} namespace"

kubectl apply --namespace="${NAMESPACE}" -f "${CUR_DIR}/zookeeper.yaml"