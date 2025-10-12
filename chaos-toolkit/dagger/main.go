package main

import (
    "context"
    "fmt"

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
    namespace string,
    deployment string,
    kubeconfigDir *dagger.Directory,
    minikubeDir *dagger.Directory,
    chaosType string,
    chaosDuration string,
    loadTestVus string,
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

    // Remaining phases - TODO
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
    kubectl *dagger.Container,
    namespace string,
    deployment string,
) error {

    fmt.Printf("  → Checking namespace '%s'...\n", namespace)
    _, err := kubectl.
        WithExec([]string{"kubectl", "get", "namespace", namespace}).
        Sync(ctx)
    if err != nil {
        return fmt.Errorf("namespace '%s' not found: %w", namespace, err)
    }
    fmt.Println("    ✓ Namespace exists")

    fmt.Printf("  → Checking deployment '%s'...\n", deployment)
    _, err = kubectl.
        WithExec([]string{"kubectl", "get", "deployment", deployment, "-n", namespace}).
        Sync(ctx)
    if err != nil {
        return fmt.Errorf("deployment '%s' not found in namespace '%s': %w", deployment, namespace, err)
    }
    fmt.Println("    ✓ Deployment exists")

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

// kubectlContainer returns a container built via the Dagger client.
// This implementation matches SDKs where dagger.Connect() returns a client value.
// It does NOT attempt to close the client because some SDK versions do not provide Close.
func (m *ChaosToolkit) kubectlContainer(ctx context.Context, kubeconfigDir *dagger.Directory) (*dagger.Container, error) {
    // Connect to the Dagger engine; adapt if your SDK has a different signature.
    client := dagger.Connect()

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

// minimal local main to keep `go build` happy; remove or replace when using generated Dagger wiring.
//func main() {
//    ctx := context.Background()
//    log.Println("ChaosToolkit build OK. Using local main.")
//    _ = ctx
//}
