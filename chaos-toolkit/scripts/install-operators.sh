#!/bin/bash
set -e

echo "📦 Installing Chaos Engineering Operators..."

# Install Litmus
echo "Installing Litmus Chaos..."
kubectl create ns litmus --dry-run=client -o yaml | kubectl apply -f -

helm repo add litmuschaos https://litmuschaos.github.io/litmus-helm/ 2>/dev/null || true
helm repo update

helm upgrade --install litmus litmuschaos/litmus \
  --namespace litmus \
  --set portal.frontend.service.type=ClusterIP \
  --wait --timeout=5m

# Create service account for chaos
kubectl apply -f - <<EOF
apiVersion: v1
kind: ServiceAccount
metadata:
  name: litmus-admin
  namespace: litmus
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: litmus-admin
rules:
  - apiGroups: [""]
    resources: ["pods", "events", "services"]
    verbs: ["get", "list", "watch", "delete", "create"]
  - apiGroups: ["apps"]
    resources: ["deployments", "replicasets"]
    verbs: ["get", "list"]
  - apiGroups: ["litmuschaos.io"]
    resources: ["*"]
    verbs: ["*"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: litmus-admin
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: litmus-admin
subjects:
  - kind: ServiceAccount
    name: litmus-admin
    namespace: litmus
EOF

# Install k6 Operator
echo "Installing k6 Operator..."
kubectl create ns k6-operator --dry-run=client -o yaml | kubectl apply -f -

helm repo add grafana https://grafana.github.io/helm-charts 2>/dev/null || true
helm repo update

helm upgrade --install k6-operator grafana/k6-operator \
  --namespace k6-operator \
  --wait --timeout=5m

echo "✅ All operators installed successfully!"
echo ""
echo "Verify installation:"
echo "  kubectl get pods -n litmus"
echo "  kubectl get pods -n k6-operator"