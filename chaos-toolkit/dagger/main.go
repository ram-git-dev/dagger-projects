package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	dagger "dagger.io/dagger"
)

type ChaosToolkit struct {
	client *dagger.Client
}

func NewChaosToolkit(ctx context.Context) (*ChaosToolkit, error) {
	client, err := dagger.Connect(ctx)
	if err != nil {
		return nil, err
	}
	return &ChaosToolkit{client: client}, nil
}

func (c *ChaosToolkit) ChaosTest(
	ctx context.Context,
	kubeconfigDir *dagger.Directory,
	minikubeDir *dagger.Directory,
	namespace string,
	deployment string,
	chaosType string,
	chaosDuration int,
	loadTestDuration string,
	loadTestVUs int,
	cleanup bool,
) error {
	fmt.Printf("Running chaos test for deployment=%s, namespace=%s, chaos=%s\n", deployment, namespace, chaosType)

	// Example: simulate chaos injection with pod-delete
	chaosCmd := []string{"echo", "Simulating chaos: " + chaosType}
	switch chaosType {
	case "pod-delete":
		chaosCmd = []string{"kubectl", "delete", "pod", "-n", namespace, "--selector=app=" + deployment}
	case "pod-network-latency":
		chaosCmd = []string{"echo", "Inject network latency"}
	case "pod-cpu-hog":
		chaosCmd = []string{"echo", "Inject CPU stress"}
	case "pod-memory-hog":
		chaosCmd = []string{"echo", "Inject Memory stress"}
	default:
		return fmt.Errorf("unknown chaos type: %s", chaosType)
	}

	// Run chaos in container
	_, err := c.client.Container().
		From("bitnami/kubectl:latest").
		WithMountedDirectory("/kube", kubeconfigDir).
		WithEnvVariable("KUBECONFIG", "/kube/config").
		WithExec(chaosCmd).
		ExitCode(ctx)
	if err != nil {
		return fmt.Errorf("failed chaos command: %w", err)
	}

	// k6 Load test
	k6Container := c.client.Container().From("loadimpact/k6:latest")
	_, err = k6Container.WithExec([]string{
		"k6", "run",
		"--vus", strconv.Itoa(loadTestVUs),
		"--duration", loadTestDuration,
		"/scripts/loadtest.js",
	}).ExitCode(ctx)
	if err != nil {
		return fmt.Errorf("failed k6 load test: %w", err)
	}

	fmt.Println("Chaos test + load test completed!")

	if cleanup {
		fmt.Println("Cleanup enabled - removing chaos artifacts")
		// example: just echo for now
		fmt.Println("Cleanup done")
	}

	return nil
}

func main() {
	ctx := context.Background()

	// Read environment/inputs
	namespace := os.Getenv("INPUT_NAMESPACE")
	if namespace == "" {
		namespace = "default"
	}
	deployment := os.Getenv("INPUT_DEPLOYMENT")
	if deployment == "" {
		deployment = "sample-app"
	}
	chaosType := os.Getenv("INPUT_CHAOS_TYPE")
	if chaosType == "" {
		chaosType = "pod-delete"
	}
	chaosDurationStr := os.Getenv("INPUT_CHAOS_DURATION")
	if chaosDurationStr == "" {
		chaosDurationStr = "60"
	}
	chaosDuration, _ := strconv.Atoi(chaosDurationStr)

	loadTestDuration := os.Getenv("INPUT_LOAD_TEST_DURATION")
	if loadTestDuration == "" {
		loadTestDuration = "5m"
	}

	loadTestVUsStr := os.Getenv("INPUT_LOAD_TEST_VUS")
	if loadTestVUsStr == "" {
		loadTestVUsStr = "10"
	}
	loadTestVUs, _ := strconv.Atoi(loadTestVUsStr)

	cleanup := true
	if os.Getenv("INPUT_CLEANUP_AFTER") == "false" {
		cleanup = false
	}

	// Connect to Dagger
	toolkit, err := NewChaosToolkit(ctx)
	if err != nil {
		log.Fatalf("Failed to create chaos toolkit: %v", err)
	}
	defer toolkit.client.Close()

	// Mount directories
	kubeconfigDir := toolkit.client.Host().Directory(os.Getenv("HOME") + "/.kube")
	minikubeDir := toolkit.client.Host().Directory(os.Getenv("HOME") + "/.minikube")

	// Run the chaos test
	if err := toolkit.ChaosTest(ctx, kubeconfigDir, minikubeDir, namespace, deployment, chaosType, chaosDuration, loadTestDuration, loadTestVUs, cleanup); err != nil {
		log.Fatalf("Chaos test failed: %v", err)
	}

	fmt.Println("All done! ✅")
}
