package main

import (
    "context"
    "fmt"
    "log"

    "dagger/chaos-toolkit/internal/dagger"
)

// ChaosToolkit is the exported receiver used by the Dagger module.
type ChaosToolkit struct{}

// Hello returns a greeting message
func (m *ChaosToolkit) Hello(ctx context.Context) string {
    return "Hello from Chaos Toolkit!"
}

// ChaosTest runs a complete chaos engineering test
func (m *ChaosToolkit) ChaosTest(
    ctx context.Context,
    // Target namespace
    namespace string,
    // Target deployment name
    deployment string,
    // Kubeconfig directory (contains config file)
    kubeconfigDir *dagger.Directory,
    // Minikube directory (optional) used when running against local Minikube
    // +optional
    minikubeDir *dagger.Directory,
    // Type of chaos experiment
    // +optional
    // +default="pod-delete"
    chaosType string,
    // Duration of chaos in seconds
    // +optional
    // +default="60"
    chaosDuration string,
    // Number of virtual users for load test
    // +optional
    // +default="10"
    loadTestVus string,
    // Load test duration
    // +optional
    // +default="5m"
    loadTestDuration string,
) (string, error) {

    // Prefer minikubeDir when provided
    if minikubeDir != nil {
        kubeconfigDir = minikubeDir
    }

    fmt.Println("🚀 Starting Chaos Engineering Test")
    fmt.Printf("Target: %s/%s\n", namespace, deployment)
    fmt.Printf("Chaos Type: %s\n", chaosType)
    fmt.Printf("Chaos Duration: %ss\n", chaosDuration)
    fmt.Printf("Load Test: %s VUs for %s\n", loadTestVus, loadTestDuration)

    // Phase 1: Pre-flight checks
    fmt.Println("\n📋 Phase 1: Pre-flight Checks")
    kubectl, err := m.kubectlContainer(ctx, kubeconfigDir)
    if err != nil {
        return "", fmt.Errorf("failed to prepare kubectl container: %w", err)
    }

    if err := m.preflightChecks(ctx, kubectl, namespace, deployment); err != nil {
        return "", fmt.Errorf("pre-flight checks failed: %w", err)
    }
    fmt.Println("✅ Pre-flight checks passed!")

    // Phase 2: Install operators (TODO)
    fmt.Println("\n📦 Phase 2: Installing Operators")
    fmt.Println("⚠️  Operator installation - Coming soon!")

    // Phase 3: Baseline test (TODO)
    fmt.Println("\n📊 Phase 3: Baseline Test")
    fmt.Println("⚠️  Baseline testing - Coming soon!")

    // Phase 4: Chaos injection (TODO)
    fmt.Println("\n💥 Phase 4: Chaos Injection")
    fmt.Printf("⚠️  Would inject %s chaos for %ss - Coming soon!\n", chaosType, chaosDuration)

    // Phase 5: Load test during chaos (TODO)
    fmt.Println("\n🔥 Phase 5: Load Test During Chaos")
    fmt.Printf("⚠️  Would run k6 with %s VUs - Coming soon!\n", loadTestVus)

    // Phase 6: Recovery measurement (TODO)
    fmt.Println("\n🔄 Phase 6: Recovery Measurement")
    fmt.Println("⚠️  Recovery testing - Coming soon!")

    // Phase 7: Report generation (TODO)
    fmt.Println("\n📄 Phase 7: Report Generation")
    fmt.Println("⚠️  Report generation - Coming soon!")

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

// preflightChecks validates that the target namespace and deployment exist
func (m *ChaosToolkit) preflightChecks(
    ctx context.Context,
    kubectl *dagger.Container,
    namespace string,
    deployment string,
) error {

    // Check namespace exists
    fmt.Printf("  → Checking namespace '%s'...\n", namespace)
    _, err := kubectl.
        WithExec([]string{"kubectl", "get", "namespace", namespace}).
        Sync(ctx)
    if err != nil {
        return fmt.Errorf("namespace '%s' not found: %w", namespace, err)
    }
    fmt.Println("    ✓ Namespace exists")

    // Check deployment exists
    fmt.Printf("  → Checking deployment '%s'...\n", deployment)
    _, err = kubectl.
        WithExec([]string{"kubectl", "get", "deployment", deployment, "-n", namespace}).
        Sync(ctx)
    if err != nil {
        return fmt.Errorf("deployment '%s' not found in namespace '%s': %w", deployment, namespace, err)
    }
    fmt.Println("    ✓ Deployment exists")

    // Check deployment status
    fmt.Println("  → Checking deployment status...")
    statusOutput, err := kubectl.
        WithExec([]string{
            "kubectl", "get", "deployment", deployment, "-n", namespace,
            "-o", "jsonpath={.status.readyReplicas}/{.status.replicas}",
        }).
        Stdout(ctx)
    if err != nil {
        return fmt.Errorf("failed to get deployment status: %w", err)
    }
    fmt.Printf("    ✓ Ready replicas: %s\n", statusOutput)

    return nil
}

// kubectlContainer returns a container with kubectl installed and kubeconfig mounted
// it connects to the dagger engine, builds the container and returns it.
// caller is responsible for providing ctx and letting the returned container be used within that ctx lifecycle.
func (m *ChaosToolkit) kubectlContainer(ctx context.Context, kubeconfigDir *dagger.Directory) (*dagger.Container, error) {
    client, err := dagger.Connect(ctx)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to dagger: %w", err)
    }
    // do not close client here if caller expects to use returned container across client lifetime;
    // closing here would invalidate returned objects. If you prefer short-lived usage, call client.Close() in caller.
    // For a simple flow we will not close here.

    kubeconfigFile := kubeconfigDir.File("config")

    ctr := client.Container().
        From("alpine:latest").
        WithExec([]string{"apk", "add", "--no-cache", "curl"}).
        WithExec([]string{"sh", "-c", "curl -LO https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"}).
        WithExec([]string{"chmod", "+x", "./kubectl"}).
        WithExec([]string{"mv", "./kubectl", "/usr/local/bin/kubectl"}).
        WithFile("/root/.kube/config", kubeconfigFile).
        WithExec([]string{"chmod", "600", "/root/.kube/config"}).
        WithEnvVariable("KUBECONFIG", "/root/.kube/config")

    return ctr, nil
}

// main is a minimal entrypoint to keep local builds simple. Remove when using generated Dagger wiring.
func main() {
    ctx := context.Background()
    log.Println("ChaosToolkit build OK. Replace main when using generated Dagger module.")
    _ = ctx
}
