set -x

source ../common/parse-or-gather-user-input.sh "${@}"

kubectl create ns $_arg_pcp_namespace

source ../common/rbac.yaml.sh $_arg_pcp_namespace
kubectl create -f protoform.yaml -n $_arg_pcp_namespace
kubectl create -f rbac.yaml -n $_arg_pcp_namespace

# Optional Prometheus
if [[ $_arg_proto_prometheus_metrics == "on" ]]; then
  kubectl create -f ../common/prometheus-deployment.yaml
fi
