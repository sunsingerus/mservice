#!/bin/bash

# Namespace to install Zookeeper
NAMESPACE="${NAMESPACE:-default}"

CUR_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)"

echo "Delete Zookeeper from ${NAMESPACE} namespace"

kubectl delete --namespace="${NAMESPACE}" -f "${CUR_DIR}/zookeeper.yaml"
