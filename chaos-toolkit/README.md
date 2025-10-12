# 🧪 Chaos Engineering Toolkit

A universal, plug-and-play chaos engineering platform that works with **any Kubernetes cluster**. Run chaos experiments and load tests against your applications with just a few clicks.

## 🚀 Features

- ✅ **Universal**: Works with any K8s cluster (minikube, EKS, GKE, AKS, etc.)
- ✅ **No Installation Required**: Everything runs via GitHub Actions
- ✅ **Multiple Chaos Types**: Pod delete, network latency, CPU/memory stress
- ✅ **Automated Load Testing**: Built-in k6 integration
- ✅ **Beautiful Reports**: HTML reports with charts and metrics
- ✅ **Safe**: Automated cleanup and recovery checks

## 📋 Prerequisites

- A Kubernetes cluster (any provider)
- `kubectl` configured locally
- A GitHub account

## 🎯 Quick Start

### 1. Fork This Repository

Click the "Fork" button at the top right of this page.

### 2. Add Your Kubeconfig

```bash
# Encode your kubeconfig
cat ~/.kube/config | base64 | pbcopy  # macOS
cat ~/.kube/config | base64 | xclip   # Linux
```

Go to your forked repo:
- Settings → Secrets and variables → Actions
- Click "New repository secret"
- Name: `KUBECONFIG_BASE64`
- Value: Paste the base64 string
- Click "Add secret"

### 3. Run Your First Chaos Test

1. Go to **Actions** tab
2. Click **"Chaos Engineering Test"**
3. Click **"Run workflow"**
4. Fill in the form:
   - **Namespace**: `default` (or your namespace)
   - **Deployment**: `my-app` (your deployment name)
   - **Chaos Type**: `pod-delete`
   - **Duration**: `60` (seconds)
   - **VUs**: `10` (virtual users)
5. Click **"Run workflow"**

### 4. View Results

After the test completes (~5-10 minutes):

1. Go to the workflow run
2. Scroll to **Artifacts** section
3. Download `chaos-report-XXX`
4. Unzip and open `report.html` in your browser

## 🧪 Supported Chaos Experiments

| Chaos Type | Description | Duration |
|------------|-------------|----------|
| `pod-delete` | Randomly deletes pods | 30-300s |
| `pod-network-latency` | Injects network delay | 30-300s |
| `pod-cpu-hog` | Consumes CPU resources | 30-300s |
| `pod-memory-hog` | Consumes memory | 30-300s |

## 📊 What Gets Tested?

The pipeline measures:

- **Error Rate**: % of failed requests during chaos
- **Latency**: p50, p95, p99 response times
- **Throughput**: Requests per second
- **Recovery Time**: How long to return to normal
- **Blast Radius**: % of pods affected

## 🎯 Success Criteria

A test **passes** if:
- ✅ Error rate < 5% during chaos
- ✅ P99 latency < 500ms
- ✅ Recovery time < 60s
- ✅ All pods recover successfully

## 🔧 Advanced Usage

### Custom Load Test Script

Create `k6-script.js` in your repo:

```javascript
import http from 'k6/http';
import { check } from 'k6';

export default function () {
  const res = http.get('http://your-service/api');
  check(res, {
    'status is 200': (r) => r.status === 200,
  });
}
```

### Run Multiple Chaos Types

You can trigger multiple tests with different chaos types to compare resilience.

### Integration with CI/CD

Add to your deployment pipeline:

```yaml
- name: Chaos Test
  uses: ./.github/workflows/chaos-test.yml
  with:
    namespace: production
    deployment: api-server
    chaos_type: pod-delete
```

## 📈 Example Report

The HTML report includes:

- **Summary Dashboard**: Pass/fail status, key metrics
- **Timeline Chart**: Error rate and latency over time
- **Chaos Events**: When chaos was injected/stopped
- **Recovery Analysis**: Pod status during recovery
- **Recommendations**: Suggestions for improvement

## 🛡️ Safety Features

- **Dry-run mode**: Test without actual chaos
- **Automatic rollback**: If error rate exceeds threshold
- **Cleanup**: Removes all test resources after completion
- **Isolation**: Operators installed in temporary namespaces

## 🤝 Contributing

Contributions welcome! Please:

1. Fork the repo
2. Create a feature branch
3. Add tests
4. Submit a PR

## 📝 License

MIT License - see LICENSE file for details

## 🙋 FAQ

### Can I use this in production?

Yes! But start with staging environments first. Use shorter chaos durations and monitor closely.

### Does this require installing anything in my cluster?

Temporarily, yes. The pipeline installs k6-operator and Litmus, but removes them after testing (if cleanup enabled).

### Can I test multiple services at once?

Not yet, but it's on the roadmap! For now, run separate workflows for each service.

### How do I customize the load test?

Edit the k6 script in `manifests/k6/test.js` to match your API endpoints and patterns.

### What if my cluster doesn't have Prometheus?

No problem! The pipeline collects metrics from k6 directly. Prometheus is optional.

## 🎓 Learn More

- [Chaos Engineering Principles](https://principlesofchaos.org/)
- [k6 Documentation](https://k6.io/docs/)
- [LitmusChaos Docs](https://docs.litmuschaos.io/)
- [Dagger Documentation](https://docs.dagger.io/)

---

**Made with ❤️ for chaos engineers everywhere**