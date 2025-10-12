// main.go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"dagger.io/dagger"
)

func main() {
	ctx := context.Background()

	// ---- Read workflow inputs from env ----
	namespace := getenv("NAMESPACE", "default")
	deployment := getenv("DEPLOYMENT", "sample-app")
	chaosType := getenv("CHAOS_TYPE", "pod-delete")
	chaosDuration := getenv("CHAOS_DURATION", "60") // seconds
	loadTestDuration := getenv("LOAD_TEST_DURATION", "5m")
	loadTestVUs := getenv("LOAD_TEST_VUS", "10")
	cleanup := getenv("CLEANUP_AFTER", "true")

	fmt.Printf("Starting Chaos Test:\nNamespace: %s\nDeployment: %s\nChaos Type: %s\nDuration: %ss\nLoad Test: %s, VUs: %s\nCleanup: %s\n",
		namespace, deployment, chaosType, chaosDuration, loadTestDuration, loadTestVUs, cleanup)

	// ---- Connect to Dagger ----
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		log.Fatal("Failed to connect to Dagger:", err)
	}
	defer client.Close()

	// ---- Apply Chaos Manifest ----
	manifestPath := filepath.Join("manifest/litmus", chaosType+".yaml")
	fmt.Println("Applying chaos manifest:", manifestPath)
	applyCmd := exec.Command("kubectl", "apply", "-f", manifestPath, "-n", namespace)
	applyCmd.Stdout = os.Stdout
	applyCmd.Stderr = os.Stderr
	if err := applyCmd.Run(); err != nil {
		log.Fatal("Failed to apply chaos manifest:", err)
	}

	// ---- Wait for chaos duration ----
	durationSec, _ := time.ParseDuration(chaosDuration + "s")
	fmt.Println("Chaos running for", durationSec)
	time.Sleep(durationSec)

	// ---- Run k6 Load Test ----
	k6Script := filepath.Join("k6", "test.js")
	fmt.Println("Running k6 test:", k6Script)
	k6Cmd := exec.Command("k6", "run",
		"--vus", loadTestVUs,
		"--duration", loadTestDuration,
		k6Script,
	)
	k6Cmd.Stdout = os.Stdout
	k6Cmd.Stderr = os.Stderr
	if err := k6Cmd.Run(); err != nil {
		log.Fatal("k6 test failed:", err)
	}

	// ---- Cleanup Chaos if requested ----
	if cleanup == "true" {
		fmt.Println("Cleaning up chaos manifest...")
		delCmd := exec.Command("kubectl", "delete", "-f", manifestPath, "-n", namespace)
		delCmd.Stdout = os.Stdout
		delCmd.Stderr = os.Stderr
		if err := delCmd.Run(); err != nil {
			log.Println("Warning: failed to cleanup manifest:", err)
		}
	}

	fmt.Println("Chaos + k6 orchestration complete ✅")
}

// getenv is a helper to read env with default
func getenv(key, def string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return def
}
