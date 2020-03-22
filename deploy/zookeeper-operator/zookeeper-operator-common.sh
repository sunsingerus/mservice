#!/bin/bash

# Namespace to install operator
OPERATOR_NAMESPACE="${OPERATOR_NAMESPACE:-default}"

#
# Check whether kubectl is available
#
function is_kubectl_available() {
    if ! kubectl version > /dev/null; then
        echo "kubectl is unavailable, can not continue"
        exit 1
    fi
}

#
# Check whether curl is available
#
function is_curl_available() {
    if ! curl --version > /dev/null; then
        echo "curl is unavailable, can not continue"
        exit 1
    fi
}

#
# Check whether wget is available
#
function is_wget_available() {
    if ! wget --version > /dev/null; then
        echo "wget is unavailable, can not continue"
        exit 1
    fi
}

#
# Check whether any download tool (curl, wget) is available
#
function check_file_getter_available() {
    if curl --version > /dev/null; then
        # curl is available - use it
        :
    elif wget --version > /dev/null; then
        # wget is available - use it
        :
    else
        echo "neither curl nor wget is available, can not continue"
        exit 1
    fi
}


#
# Check whether envsubst is available
#
function check_envsubst_available() {
    if ! envsubst --version > /dev/null; then
        echo "envsubst is unavailable, can not continue"
        exit 1
    fi
}

#
# Get file
#
function get_file() {
    local URL="$1"

    if curl --version > /dev/null; then
        # curl is available - use it
        curl -s "${URL}"
    elif wget --version > /dev/null; then
        # wget is available - use it
        wget -qO- "${URL}"
    else
        echo "neither curl nor wget is available, can not continue"
        exit 1
    fi
}
