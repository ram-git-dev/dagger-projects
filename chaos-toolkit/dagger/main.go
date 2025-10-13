package main

     (
    "context"
    "fmt"
    "strings"

    chaosdagger "dagger/chaos-toolkit/internal/chaosdagger"

    dg "dagger.io/dagger"
)

type ChaosToolkit struct{}

func (m *ChaosToolkit) Hello(ctx context.Context) string {
    return "Hello from Chaos Toolkit!"
}

func (m *ChaosToolkit) ChaosTest(
    ctx context.Context,
    namespace string,
    deployment string,
    kubeconfigDir *dg.Directory,
    minikubeDir *dg.Directory,
    chaosType string,
    chaosDuration string,
    loadTestDuration string,
    loadTestVus string,
    cleanup bool,
) (string, error) {

    if minikubeDir != nil {
        kubeconfigDir = minikubeDir
    }

    fmt.Println("🚀 Starting Chaos Engineering Test")
    fmt.Printf("Target: %s/%s\n", namespace, deployment)
    fmt.Printf("Chaos Type: %s\n", chaosType)
    fmt.Printf("Chaos Duration: %ss\n", chaosDuration)
    fmt.Printf("Load Test: %s VUs for %s\n", loadTestVus, loadTestDuration)
    fmt.Printf("Cleanup after test: %v\n", cleanup)

    fmt.Println("\n📋 Phase 1: Pre-flight Checks")
    kubectl, err := m.kubectlContainer(ctx)
    if err != nil {
        return "", fmt.Errorf("failed to prepare kubectl container: %w", err)
    }

    if err := m.preflightChecks(ctx, kubectl, namespace, deployment); err != nil {
        return "", fmt.Errorf("pre-flight checks failed: %w", err)
    }
    fmt.Println("✅ Pre-flight checks passed!")

    fmt.Println("\n📦 Phase 2: Installing Operators - TODO")
    fmt.Println("\n📊 Phase 3: Baseline Test - TODO")
    fmt.Println("\n💥 Phase 4: Chaos Injection - TODO")
    fmt.Println("\n🔥 Phase 5: Load Test During Chaos - TODO")
    fmt.Println("\n🔄 Phase 6: Recovery Measurement - TODO")
    fmt.Println("\n📄 Phase 7: Report Generation - TODO")

    result := fmt.Sprintf(`
✅ Chaos Test Complete!

Target: %s/%s
Chaos Type: %s
Duration: %ss
Load Test: %s VUs for %s

Status: Pre-flight checks passed ✅
Next: Implement chaos injection, load testing, and reporting
`, namespace, deployment, chaosType, chaosDuration, loadTestVus, loadTestDuration)

    return result, nil
}

func (m *ChaosToolkit) preflightChecks(
    ctx context.Context,
    kubectl *dg.Container,
    namespace string,
    deployment string,
) error {

    fmt.Printf("  → Checking namespace '%s'...\n", namespace)
    _, err := kubectl.WithExec([]string{"kubectl", "get", "namespace", namespace}).Sync(ctx)
    if err != nil {
        return fmt.Errorf("namespace '%s' not found: %w", namespace, err)
    }
    fmt.Println("    ✓ Namespace exists")

    fmt.Printf("  → Checking deployment '%s'...\n", deployment)
    _, err = kubectl.WithExec([]string{"kubectl", "get", "deployment", deployment, "-n", namespace}).Sync(ctx)
    if err != nil {
        return fmt.Errorf("deployment '%s' not found in namespace '%s': %w", deployment, namespace, err)
    }
    fmt.Println("    ✓ Deployment exists")

    fmt.Println("  → Checking deployment status...")
    statusOutput, err := kubectl.WithExec([]string{
        "kubectl", "get", "deployment", deployment, "-n", namespace,
        "-o", "jsonpath={.status.readyReplicas}/{.status.replicas}",
    }).Stdout(ctx)
    if err != nil {
        return fmt.Errorf("failed to get deployment status: %w", err)
    }
    fmt.Printf("    ✓ Ready replicas: %s\n", statusOutput)

    return nil
}

func (m *ChaosToolkit) kubectlContainer(ctx context.Context) (*dg.Container, error) {
    client, err := dg.Connect(ctx)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to Dagger: %w", err)
    }

    saDir := client.Host().Directory("/var/run/secrets/kubernetes.io/serviceaccount")

    token, err := saDir.File("token").Contents(ctx)
    if err != nil {
        return nil, fmt.Errorf("failed to read token: %w", err)
    }
    namespace, err := saDir.File("namespace").Contents(ctx)
    if err != nil {
        return nil, fmt.Errorf("failed to read namespace: %w", err)
    }

    kubeconfigYAML := fmt.Sprintf(`
apiVersion: v1
kind: Config
clusters:
- name: in-cluster
  cluster:
    server: https://kubernetes.default.svc
    certificate-authority: /var/run/secrets/kubernetes.io/serviceaccount/ca.crt
contexts:
- name: in-cluster
  context:
    cluster: in-cluster
    user: in-cluster
    namespace: %s
current-context: in-cluster
users:
- name: in-cluster
  user:
    token: %s
`, strings.TrimSpace(namespace), strings.TrimSpace(token))

    kubeconfigDir := client.Directory().WithNewFile("config", kubeconfigYAML)

    ctr := client.Container().
        From("alpine:latest").
        WithExec([]string{"apk", "add", "--no-cache", "curl", "bash"}).
        WithExec([]string{"sh", "-c", `
curl -LO https://dl.k8s.io/release/$(curl -sL https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl &&
chmod +x kubectl &&
mv kubectl /usr/local/bin/kubectl
`}).
        WithDirectory("/root/.kube", kubeconfigDir).
        WithFile("/var/run/secrets/kubernetes.io/serviceaccount/ca.crt", saDir.File("ca.crt")).
        WithEnvVariable("KUBECONFIG", "/root/.kube/config")

    return ctr, nil
}
