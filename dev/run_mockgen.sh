#!/bin/bash

CUR_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)"
source "${CUR_DIR}/go_build_config.sh"

GO111MODULE=on ~/gocode/bin/mockgen \
    -package controller_client \
    github.com/sunsingerus/mservice/pkg/api/mservice MServiceControlPlaneClient,MServiceControlPlane_DataClient > ${SRC_ROOT}/pkg/controller/client/client_test_mock_mservice.go
