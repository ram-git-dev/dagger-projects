# Chaos Engineering Toolkit

Automated chaos testing pipeline for Kubernetes using Dagger, LitmusChaos, and k6. Connect any cluster, select your target, run experiments, get results.

## What It Does

- Injects faults into your K8s deployments (pod kills, latency, resource stress)
- Runs concurrent load tests to measure impact
- Collects metrics and generates HTML reports
- All orchestrated via GitHub Actions + Dagger

## Architecture

```
GitHub Actions → Dagger Pipeline → Your K8s Cluster
                      ↓
              ┌───────┴────────┐
              ↓                ↓
         LitmusChaos         k6
         (fault injection)   (load generation)
              ↓                ↓
           Metrics Collection
              ↓
         HTML Report
```

## Prerequisites

- Kubernetes cluster with `kubectl` access
- GitHub account
- Dagger CLI (optional for local runs)

## Setup

### 1. Add Kubeconfig Secret

```bash
cat ~/.kube/config | base64 > kubeconfig.b64
```

Add to GitHub repo:
- Settings → Secrets → Actions
- Name: `KUBECONFIG_BASE64`
- Value: contents of `kubeconfig.b64`

### 2. Run Test

Actions → "Chaos Engineering Test" → Run workflow

**Required inputs:**
- `namespace`: Target namespace
- `deployment`: Target deployment name  
- `chaos_type`: Fault type (see below)
- `chaos_duration`: Duration in seconds
- `load_test_vus`: Concurrent users
- `load_test_duration`: Load test duration

### 3. Get Results

Download artifact `chaos-report-{run_number}.zip` from workflow run. Contains:
- `report.html` - Visual report with charts
- `summary.json` - Raw metrics

## Supported Chaos Types

| Type | Effect | Tunable Parameters |
|------|--------|-------------------|
| `pod-delete` | Kills pods randomly | `CHAOS_INTERVAL`, `PODS_AFFECTED_PERC` |
| `pod-network-latency` | Injects network delay | `NETWORK_LATENCY` (ms) |
| `pod-cpu-hog` | Consumes CPU | `CPU_CORES`, `CPU_LOAD` |
| `pod-memory-hog` | Consumes memory | `MEMORY_CONSUMPTION` (MB) |

## Metrics Collected

- **Error rate**: Failed requests / total requests
- **Latency percentiles**: p50, p95, p99
- **Throughput**: Requests per second
- **Recovery time**: Time to restore healthy state
- **Blast radius**: Percentage of pods affected

## Test Phases

```
1. Pre-flight checks (namespace, deployment validation)
2. Operator installation (Litmus, k6)
3. Baseline test (2min, no chaos)
4. Chaos injection + load test (configurable duration)
5. Recovery measurement
6. Report generation
7. Cleanup (optional)
```

## Pass/Fail Criteria

Test passes if:
- Error rate < 5%
- P99 latency < 500ms
- Recovery time < 60s

Thresholds configurable in `dagger/main.go`.

## Local Development

```bash
cd chaos-toolkit/dagger

# Test Dagger pipeline locally
dagger call chaos-test \
  --namespace=default \
  --deployment=nginx \
  --chaos-type=pod-delete \
  --chaos-duration=60 \
  --load-test-duration=5m \
  --load-test-vus=10

# Output saved to output/
```

## Customization

### Custom Load Test

Edit `manifests/k6/test.js`:

```javascript
export default function () {
  http.post('http://my-api/endpoint', JSON.stringify({...}));
}
```

### Custom Chaos Parameters

Edit manifest templates in `manifests/litmus/*.yaml`. Variables replaced at runtime:
- `{{NAMESPACE}}`
- `{{DEPLOYMENT}}`
- `{{CHAOS_DURATION}}`

### Custom Thresholds

Edit `dagger/main.go`:

```go
const (
    MaxErrorRate = 0.05    // 5%
    MaxP99Latency = 500.0  // ms
    MaxRecoveryTime = 60   // seconds
)
```

## File Structure

```
chaos-toolkit/
├── .github/workflows/
│   └── chaos-test.yml          # GitHub Actions workflow
├── dagger/
│   ├── main.go                 # Dagger pipeline logic
│   ├── go.mod
│   └── go.sum
├── manifests/
│   ├── litmus/
│   │   ├── pod-delete.yaml
│   │   ├── network-latency.yaml
│   │   ├── cpu-hog.yaml
│   │   └── memory-hog.yaml
│   └── k6/
│       ├── test.js             # Load test script
│       └── testrun.yaml        # K6 TestRun CRD
├── templates/
│   └── report.html             # HTML report template
└── scripts/
    ├── install-operators.sh
    └── cleanup.sh
```

## CI/CD Integration

Trigger from deployment pipeline:

```yaml
- name: Chaos Test
  uses: ram-git-dev/dagger-projects/.github/workflows/chaos-test.yml@main
  with:
    namespace: staging
    deployment: api-server
    chaos_type: pod-delete
    chaos_duration: 60
```

## Dependencies

**Installed by pipeline:**
- LitmusChaos 3.x
- k6-operator 0.x

**Runtime:**
- Dagger 0.11+
- Go 1.21+

## Troubleshooting

**"namespace not found"**
- Verify namespace exists: `kubectl get ns`

**"deployment not found"**  
- Check deployment name: `kubectl get deploy -n <namespace>`

**"chaos experiments failing"**
- Check ServiceAccount permissions
- Verify Litmus installed: `kubectl get pods -n litmus`

**"no metrics collected"**
- Check k6 pods: `kubectl get pods -n <namespace> -l runner=chaos-load-test`
- Check logs: `kubectl logs -n <namespace> -l runner=chaos-load-test`

## Contributing

PRs welcome. Please:
1. Test locally with Dagger first
2. Ensure Go code passes `go vet` and `go fmt`
3. Update README if adding features

## License

MIT

## References

- [Principles of Chaos Engineering](https://principlesofchaos.org/)
- [LitmusChaos Docs](https://docs.litmuschaos.io/)
- [k6 Load Testing](https://k6.io/docs/)
- [Dagger Documentation](https://docs.dagger.io/)