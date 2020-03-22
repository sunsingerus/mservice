#!/bin/bash

CUR_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)"
source "${CUR_DIR}/zookeeper-operator-common.sh"

# OPERATOR_NAMESPACE is declared in common/configuration file. Check in there for default value

echo "Delete operator from ${OPERATOR_NAMESPACE} namespace"

# Operator
kubectl delete --namespace="${OPERATOR_NAMESPACE}" -f https://raw.githubusercontent.com/pravega/zookeeper-operator/master/deploy/default_ns/operator.yaml


# RBAC
kubectl delete --namespace="${OPERATOR_NAMESPACE}" -f <( \
    get_file https://raw.githubusercontent.com/pravega/zookeeper-operator/master/deploy/default_ns/rbac.yaml | \
        OPERATOR_NAMESPACE=${OPERATOR_NAMESPACE} sed "s/namespace: default/namespace: ${OPERATOR_NAMESPACE}/" \
)

# CRD
kubectl delete -f https://raw.githubusercontent.com/pravega/zookeeper-operator/master/deploy/crds/zookeeper_v1beta1_zookeepercluster_crd.yaml
