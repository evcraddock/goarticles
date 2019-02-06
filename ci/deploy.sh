#!/bin/sh

curl -LO https://storage.googleapis.com/kubernetes-release/release/v1.13.0/bin/linux/amd64/kubectl  && \
    chmod +x ./kubectl && \
    mv ./kubectl /usr/local/bin/kubectl

echo "${GKE_KEY}" > key.json

gcloud auth activate-service-account --key-file key.json

gcloud container clusters get-credentials "${GKE_CLUSTER}" --project "${GKE_PROJECT}" --region "${GKE_REGION}"

repository=$(cat docker-file/repository)
tag=$(cat version/version)

IMAGE_NAME=${repository}:${tag}
echo "Deploying ${IMAGE_NAME}"
sed -i -e "s,((IMAGE_NAME)),$IMAGE_NAME,g" source/deployments/api/deployment.yml

kubectl apply -f source/deployments/api/deployment.yml

