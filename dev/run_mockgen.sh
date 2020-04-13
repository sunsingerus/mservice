#!/bin/bash

CUR_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)"
source "${CUR_DIR}/go_build_config.sh"

GO111MODULE=on ~/gocode/bin/mockgen \
    github.com/sunsingerus/mservice/pkg/api/mservice MServiceControlPlaneClient,MServiceControlPlane_DataClient > ${SRC_ROOT}/pkg/mock/mock_mservice.go
