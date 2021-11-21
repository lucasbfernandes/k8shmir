#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

readonly KIND_VERSION=${KIND_VERSION:-v0.10.0}
readonly HELM_VERSION=3.5.4
readonly CLUSTER_NAME=e2e-test
readonly ROOT=$( cd -- "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )
readonly KUBECONFIG_PATH=$ROOT/$CLUSTER_NAME.kubeconfig
readonly HELM_CONTAINER_NAME=helm_builder

install_kind() {
    if ! which kind >/dev/null; then
        echo 'Installing Kind'
        curl -sSLo /tmp/kind "https://github.com/kubernetes-sigs/kind/releases/download/${KIND_VERSION}/kind-linux-amd64"
        chmod +x /tmp/kind
        mv /tmp/kind /usr/local/bin/kind
    else
        echo 'Skipping Kind installation'
    fi
    echo $(kind --version)
}

create_kind_cluster() {
    echo 'Create Kind cluster'
    kind delete cluster --name "$CLUSTER_NAME" > /dev/null 2>&1 || true
    kind create cluster \
      --name "$CLUSTER_NAME" \
      --kubeconfig $KUBECONFIG_PATH \
      --config $ROOT/kind-config.yaml
}

deploy_infrastructure() {
    echo 'Setting up deployment container'
    setup_deployment_container

    echo 'Deploying helm applications'
    deploy_ingress_controller
    deploy_atomix_charts
    deploy_replicated_apps

    echo 'Killing deployment container'
    docker kill $HELM_CONTAINER_NAME > /dev/null 2>&1 || true
    docker rm $HELM_CONTAINER_NAME > /dev/null 2>&1 || true
}

# TODO improve path management
setup_deployment_container() {
    echo 'Creating helm/kubectl container'
    docker kill $HELM_CONTAINER_NAME > /dev/null 2>&1 || true
    docker run -it -d \
      --entrypoint '/bin/sh' \
      --network host \
      --name $HELM_CONTAINER_NAME \
      --volume $ROOT/charts:/e2e/charts \
      --volume $ROOT/../../../install:/e2e/install \
      --workdir /e2e \
      dtzar/helm-kubectl:$HELM_VERSION

    echo 'Copying kubeconfig to container'
    docker exec -i $HELM_CONTAINER_NAME mkdir -p /root/.kube
    docker cp $KUBECONFIG_PATH $HELM_CONTAINER_NAME:/root/.kube/config
}

deploy_ingress_controller() {
    echo 'Creating nginx ingress controller'
    docker exec -i $HELM_CONTAINER_NAME \
      kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/master/deploy/static/provider/kind/deploy.yaml

    # TODO improve this
    sleep 20

    docker exec -i $HELM_CONTAINER_NAME \
      kubectl wait --namespace ingress-nginx \
        --for=condition=ready pod \
        --selector=app.kubernetes.io/component=controller \
        --timeout=120s
}

deploy_atomix_charts() {
    echo 'Creating atomix infrastructure'
    docker exec -i $HELM_CONTAINER_NAME helm install atomix install/helm-chart --atomic --debug
}

# Change this method to deploy other apps
deploy_replicated_apps() {
    echo 'Deploying replicated apps'
    docker exec -i $HELM_CONTAINER_NAME helm install stress-app charts/stress --atomic --debug
}

main() {
    install_kind
    create_kind_cluster
    deploy_infrastructure
}

main