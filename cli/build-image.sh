#!/bin/bash

trap cleanup SIGTERM SIGINT
cleanup(){
	echo "stopped building"
	rm ./main
	exit
}

IMAGE_NAME="crd-controller"
REPOSITORY="localhost:5000"
TAG="latest"

set -o errexit
set -o nounset
set -o pipefail

# build binary
echo "compiling binary..."
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .
echo "...done"

# adding minikube environment to docker
eval $(minikube docker-env)

echo "building new images..."

docker build -t ${IMAGE_NAME}:${TAG} .
docker tag ${IMAGE_NAME}:${TAG} ${REPOSITORY}/${IMAGE_NAME}:${TAG}
docker push ${REPOSITORY}/${IMAGE_NAME}:${TAG}

echo "...done"

rm ./main


