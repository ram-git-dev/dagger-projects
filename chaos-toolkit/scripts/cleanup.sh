#!/bin/bash
set -e

echo "🧹 Cleaning up chaos testing resources..."

# Delete chaos engines
echo "Removing chaos engines..."
kubectl delete chaosengine --all --all-namespaces --ignore-not-found=true

# Delete k6 test runs
echo "Removing k6 test runs..."
kubectl delete testrun --all --all-namespaces --ignore-not-found=true

# Optional: Remove operators (comment out if you want to keep them)
# echo "Removing operators..."
# helm uninstall litmus -n litmus || true
# helm uninstall k6-operator -n k6-operator || true
# kubectl delete ns litmus k6-operator --ignore-not-found=true

echo "✅ Cleanup complete!"