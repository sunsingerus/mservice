#!/bin/bash

# Verify installation

watch -n1 "kubectl --namespace cert-manager get pods"

function cat_test() {
cat <<EOF
apiVersion: v1
kind: Namespace
metadata:
  name: cert-manager-test
---
apiVersion: cert-manager.io/v1alpha2
kind: Issuer
metadata:
  name: test-selfsigned
  namespace: cert-manager-test
spec:
  selfSigned: {}
---
apiVersion: cert-manager.io/v1alpha2
kind: Certificate
metadata:
  name: selfsigned-cert
  namespace: cert-manager-test
spec:
  dnsNames:
    - example.com
  secretName: selfsigned-cert-tls
  issuerRef:
    name: test-selfsigned
EOF
}

cat_test | kubectl apply -f -

watch -n1 "kubectl --namespace cert-manager-test get issuer,certificate"

kubectl --namespace cert-manager-test describe issuer,certificate | less -S

cat_test | kubectl delete -f -

clear
echo "Now read on hot to setup certificates"
echo "https://cert-manager.io/docs/configuration/"
