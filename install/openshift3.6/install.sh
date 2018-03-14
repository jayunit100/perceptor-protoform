#!/bin/bash
source ../parse-or-gather-user-input.sh "${@}"

oc new-project $_arg_pcp_namespace

source ../oadm-policy-init.sh $arg_pcp_namespace

oc project $_arg_pcp_namespace
../common/oadm-init.sh
oc create -f config.yml
oc create -f proto.yml

# Optional Prometheus
if [[ $_arg_proto_prometheus_metrics == "on" ]]; then
  oc create -f ../common/prometheus-deployment.yaml
fi
