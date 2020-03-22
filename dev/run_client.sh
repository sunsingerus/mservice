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
    	-alsologtostderr=true \
    	-log_dir=log \
	-tls \
	-service-address="localhost:10000" \
	-read-filename="${CUR_DIR}/example.txt" \
	-oauth \
	-client-id="98ea8297-de36-4296-a923-07caf5cb450a" \
	-client-secret="6b4a1413-6563-4e5f-9722-d4acfc4f43dd" \
	-token-url="http://localhost:8080/auth/realms/realm1/protocol/openid-connect/token" \
    	-v=1
#	-logtostderr=true \
#	-stderrthreshold=FATAL \
#	-client-id="client1" \
#	-client-secret="6d71e2dc-7f2d-4681-8c2b-571e7ede18f8" \

# -log_dir=log Log files will be written to this directory instead of the default temporary directory
# -alsologtostderr=true Logs are written to standard error as well as to files
# -logtostderr=true  Logs are written to standard error instead of to files
# -stderrthreshold=FATAL Log events at or above this severity are logged to standard	error as well as to files

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
