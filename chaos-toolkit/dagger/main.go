package main

import (
  "context"
  "fmt"
  "log"
  "os"
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

    // If you still want this process to run the orchestration as a client,
    // move the client-orchestration code into a separate command or binary.
}

// getenv is a helper to read env with default
func getenv(key, def string) string {
    if val := os.Getenv(key); val != "" {
        return val
    }
    return def
}
