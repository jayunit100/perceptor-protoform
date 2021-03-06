#!/bin/bash
cat << EOF > rbac.yaml
kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1beta1
metadata:
  name: protoform
subjects:
- kind: ServiceAccount
  name: protoform
  namespace: ${_arg_pcp_namespace}
roleRef:
  kind: ClusterRole
  name: cluster-admin
  apiGroup: ""
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: protoform
  namespace: ${_arg_pcp_namespace}
EOF
