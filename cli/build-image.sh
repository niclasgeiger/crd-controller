#!/bin/bash

IMAGE_NAME="crd-controller"
REPOSITORY="localhost:5000"
TAG="latest"

set -o errexit
set -o nounset
set -o pipefail

# adding minikube environment to docker
eval $(minikube docker-env)

echo "building new images..."

docker build -t ${IMAGE_NAME}:${TAG} .
docker tag ${IMAGE_NAME}:${TAG} ${REPOSITORY}/${IMAGE_NAME}:${TAG}
docker push ${REPOSITORY}/${IMAGE_NAME}:${TAG}

echo "...done"




