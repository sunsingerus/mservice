#!/bin/bash

# Run binary
# Do not forget to update version

# Source configuration
CUR_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)"
source "${CUR_DIR}/go_build_config.sh"
LOG_DIR="${CUR_DIR}/log"

EXECUTABLE_BINARY="${EXECUTABLE_BINARY:-$CLIENT_BIN}"

if [[ $1 == "nobuild" ]]; then
    echo "Build step skipped, starting old binary"
else
    echo -n "Building ${EXECUTABLE_BINARY}, please wait..."
    if "${CUR_DIR}/go_build_client.sh"; then
        echo "Successfully built ${EXECUTABLE_BINARY}."
    else
        echo "Unable to build ${EXECUTABLE_BINARY}. Abort."
        exit 1
    fi
fi

if [[ ! -x "${EXECUTABLE_BINARY}" ]]; then
    echo "Unable to start ${EXECUTABLE_BINARY} Is not executable or found. Abort"
    exit 2
fi

    echo "Starting ${EXECUTABLE_BINARY}..."

    mkdir -p "${LOG_DIR}"
    "${EXECUTABLE_BINARY}" \
        send \
    	--service-address "localhost:10000" \
    	--file "${CUR_DIR}/example.txt"

if [[ $2 == "noclean" ]]; then
    echo "Clean step skipped"
else
    # And clean binary after run. It'll be rebuilt next time
    "${CUR_DIR}/go_build_client_clean.sh"
fi

#    echo "======================"
#    echo "=== Logs available ==="
#    echo "======================"
#    ls "${LOG_DIR}"/*
