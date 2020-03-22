#!/bin/bash

# Universal docker image builder

# Externally configurable build-dependent options
TAG="${TAG}"
DOCKERHUB_LOGIN="${DOCKERHUB_LOGIN}"
DOCKERHUB_PUBLISH="${DOCKERHUB_PUBLISH:-yes}"
MINIKUBE="${MINIKUBE:-no}"

DOCKERFILE_DIR="${DOCKERFILE_DIR}"
DOCKERFILE="${DOCKERFILE}"

if [[ -z "${TAG}" ]]; then
    echo "TAG has to be specified. Abort"
    exit 1
fi

if [[ ! -z "${DOCKERFILE}" ]]; then
    # Dockerfile is explicitly specified, ok
    :
elif [[ ! -z "${DOCKERFILE_DIR}" ]]; then
    #
    DOCKERFILE="${DOCKERFILE_DIR}/Dockerfile"
else
    echo "DOCKERFILE or DOCKERFILE_DIR has to be specified. Abort"
    exit 1
fi

# Source-dependent options
CUR_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)"
SRC_ROOT="$(realpath "${CUR_DIR}/..")"
source "${CUR_DIR}/go_build_config.sh"

# Build clickhouse-operator install .yaml manifest
#"${MANIFESTS_ROOT}/operator/build-clickhouse-operator-install-yaml.sh"

# Build image with Docker
if [[ "${MINIKUBE}" == "yes" ]]; then
    # We'd like to build for minikube
    eval $(minikube docker-env)
fi
docker build -t "${TAG}" -f "${DOCKERFILE}" "${SRC_ROOT}"

# Publish image
if [[ "${DOCKERHUB_PUBLISH}" == "yes" ]]; then
    if [[ ! -z "${DOCKERHUB_LOGIN}" ]]; then
        echo "Dockerhub login specified: '${DOCKERHUB_LOGIN}', perform login"
        docker login -u "${DOCKERHUB_LOGIN}"
    fi
    docker push "${TAG}"
fi
