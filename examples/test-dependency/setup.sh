#!/bin/bash

# SPDX-License-Identifier: Apache-2.0
# Copyright (c) 2020 Intel Corporation

set -o errexit
set -o nounset
set -o pipefail


source ../multi-cluster/_common.sh

test_folder=../../tests/
demo_folder=../../demo/
deployment_folder=../../../deployments/

HOST_IP=${HOST_IP:-"oops"}
KUBE_PATH1=${KUBE_PATH1:-"oops"}
KUBE_PATH2=${KUBE_PATH2:-"oops"}
OUTPUT_DIR=output

function create_common_values {
    local output_dir=$1
    local host_ip=$2

    create_apps $output_dir
    create_config_file $host_ip

        cat << NET > values.yaml
    PackagesPath: $output_dir
    ProjectName: proj-depend-1
    ClusterProvider: provider-depend
    ClusterLabel: edge-cluster
    LogicalCloud: default
    CompositeApp: test-composite-app
    CompositeProfile: test-composite-profile
    GenericPlacementIntent: test-placement-intent
    DeploymentIntentGroup: test-deployment-intent
    Intent: intent
    RsyncPort: 30431
    GacPort: 30433
    OvnPort: 30432
    DtcPort: 30448
    NpsPort: 30438
    HostIP: $host_ip
    APP1: collectd
    APP2: operator
    APP3: nginx
    Clusters:
      - Name: cluster1
        KubeConfig: $KUBE_PATH1
      - Name: cluster2
        KubeConfig: $KUBE_PATH2
    Applist:
      - Name: collectd
        Cluster:
          - cluster1
          - cluster2
      - Name: operator
        Cluster:
          - cluster1
      - Name: nginx
        Cluster:
          - cluster1
          - cluster2

NET
}


function usage {
    echo "Usage: $0 create|cleanup"
}

function cleanup {
    rm -f *.tar.gz
    rm -f values.yaml
    rm -f emco-cfg.yaml
    rm -rf $OUTPUT_DIR
}

case "$1" in
    "create" )
        if [ "${HOST_IP}" == "oops" ] ; then
            echo -e "ERROR - Environment varaible HOST_IP must be set"
            exit
        fi
        if [ "${KUBE_PATH1}" == "oops"  ] ; then
            echo -e "ERROR - KUBE_PATH1 must be defined"
            exit
        fi
        if [ "${KUBE_PATH2}" == "oops" ] ; then
            echo -e "KUBE_PATH2 must be defined"
        else
            create_common_values $OUTPUT_DIR $HOST_IP
            echo "Done create!!!"
        fi
        ;;
    "cleanup" )
        cleanup
    ;;
    *)
        usage ;;
esac
