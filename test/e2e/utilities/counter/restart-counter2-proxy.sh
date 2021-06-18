#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

readonly HELM_VERSION=3.5.4
readonly CLUSTER_NAME=e2e-test
readonly ROOT=$( cd -- "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )
readonly KUBECONFIG_PATH=$ROOT/../../bootstrap/$CLUSTER_NAME.kubeconfig
readonly KUBECTL_CONTAINER_NAME=kube_builder

# TODO improve path management
setup_deployment_container() {
    echo 'Creating helm/kubectl container'
    docker kill $KUBECTL_CONTAINER_NAME > /dev/null 2>&1 || true
    docker rm $KUBECTL_CONTAINER_NAME > /dev/null 2>&1 || true
    docker run -it -d \
      --entrypoint '/bin/sh' \
      --network host \
      --name $KUBECTL_CONTAINER_NAME \
      --workdir /e2e \
      dtzar/helm-kubectl:$HELM_VERSION

    echo 'Copying kubeconfig to container'
    docker exec -i $KUBECTL_CONTAINER_NAME mkdir -p /root/.kube
    docker cp $KUBECONFIG_PATH $KUBECTL_CONTAINER_NAME:/root/.kube/config
}

restart_app_container() {
    echo 'Restarting proxy container'
    docker exec -i $KUBECTL_CONTAINER_NAME \
      kubectl exec counter2 -c k8s-smr -- /bin/sh -c "kill 1"

    # TODO improve this
    sleep 10

    echo 'Waiting for container restart and pod ready status'
    docker exec -i $KUBECTL_CONTAINER_NAME \
      kubectl wait \
        --for=condition=ready \
        pod/counter2\
        --timeout=360s

    echo 'Killing kubectl container'
    docker kill $KUBECTL_CONTAINER_NAME > /dev/null 2>&1 || true
    docker rm $KUBECTL_CONTAINER_NAME > /dev/null 2>&1 || true
}

main() {
    setup_deployment_container
    restart_app_container
    sleep 5
}

main
