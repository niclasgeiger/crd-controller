#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT=$(dirname ${BASH_SOURCE})/..
CODEGEN_PKG=${GOPATH}/src/k8s.io/code-generator

vendor/k8s.io/code-generator/generate-groups.sh all \
  github.com/niclasgeiger/crd-controller/client github.com/niclasgeiger/crd-controller/apis \
  niclasgeiger.com:v1
