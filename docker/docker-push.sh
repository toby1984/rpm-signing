#!/bin/bash

set -e -x -o pipefail

source ../common.sh

if [ "$#" -gt "1" ] ; then
  echo "Usage: $0 <docker image prefix>"
  exit 1
fi
PREFIX="$1"

IMAGE_NAME="$PREFIX${IMAGE_NAME}"
FINAL_NAME="${IMAGE_NAME}:${IMAGE_VERSION}"

echo "PREFIX: $PREFIX"
echo "IMAGE_NAME: $IMAGE_NAME"
echo "FINAL_NAME: $FINAL_NAME"

docker tag ${FINAL_NAME} ${FINAL_NAME}
docker push ${FINAL_NAME} 

read -p "Push LATEST tag as well (Y/n) ?" answer
if [[ -z "$answer" || "$answer" == "Y" || "$answer" == "y" ]] ; then 
  docker tag ${FINAL_NAME} $IMAGE_NAME:latest
  docker push $IMAGE_NAME:latest
fi
