#!/bin/bash

# Exit immediately when a command fails
set -o errexit
# Error on unset variables
set -o nounset
# Only exit with zero if all commands of the pipeline exit successfully
set -o pipefail

# Source configuration
CUR_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)"
source "${CUR_DIR}/go_build_config.sh"

echo "Generating code with the following options:"
echo "      SRC_ROOT=${SRC_ROOT}"
echo ""
echo ""

PROTOC="${PROTOC:-protoc}"

# Check protoc is available
if [[ "${PROTOC}" > /dev/null ]]; then
    :
else
    echo "${PROTOC} is not available. Abort"
    exit 1
fi

PROTO_ROOT="${PKG_ROOT}/api"

function generate_from_proto() {
    FOLDER="${1}"

    if [[ -z "${FOLDER}" ]]; then
        echo "need to specify folder where to look for .proto files to generate code from "
        exit 1
    fi

    echo "Generate code from .proto files in ${FOLDER}"

    echo "Clean previously generated files"
    rm -f "${FOLDER}"/*.pb.go

    echo "Compile .proto files"
    # --go_out requires list of plugins to be used
    "${PROTOC}" -I "${FOLDER}" --go_out=plugins=grpc:"${FOLDER}" "${FOLDER}"/*.proto

    #protoc -I "${SRC_ROOT}" --go_out="${SRC_ROOT}" ./mservice.proto
}

generate_from_proto "${PROTO_ROOT}"/mservice
generate_from_proto "${PROTO_ROOT}"/health
