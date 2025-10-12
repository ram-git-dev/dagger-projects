package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"dagger.io/dagger"
)

// ChaosToolkit struct
type ChaosToolkit struct{}

// ChaosTestResult holds final test output
type ChaosTestResult struct {
	Passed       bool    `json:"passed"`
	ErrorRate    float64 `json:"errorRate"`
	P99Latency   int     `json:"p99Latency"`
	RecoveryTime int     `json:"recoveryTime"`
}

// ChaosTest runs chaos and load tests
func (m *ChaosToolkit) ChaosTest(
	ctx context.Context,
	namespace string,
	deployment string,
	chaosType string,
	chaosDuration int,
	loadDuration string,
	loadVUs int,
	cleanup bool,
	kubeconfigDir *dagger.Directory,
	minikubeDir *dagger.Directory,
) (string, error) {

	// Get kubeconfig
	kubeconfigFile := kubeconfigDir.File("config")

	// Dagger container for kubectl
	kubectl := dagger.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "--no-cache", "curl"}).
		WithExec([]string{"sh", "-c", "curl -LO https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"}).
		WithExec([]string{"chmod", "+x", "./kubectl"}).
		WithExec([]string{"mv", "./kubectl", "/usr/local/bin/kubectl"}).
		WithFile("/root/.kube/config", kubeconfigFile).
		WithDirectory("/home/rbot/.minikube", minikubeDir).
		WithExec([]string{"chmod", "600", "/root/.kube/config"}).
		WithEnvVariable("KUBECONFIG", "/root/.kube/config")

	fmt.Printf("Running chaos test on %s/%s: %s\n", namespace, deployment, chaosType)

	// Inject chaos
	switch chaosType {
	case "pod-delete":
		_, err := kubectl.WithExec([]string{
			"kubectl", "delete", "pod", "-l", fmt.Sprintf("app=%s", deployment), "-n", namespace,
		}).Stdout(ctx)
		if err != nil {
			return "", err
		}

	case "pod-network-latency":
		// Use tc/netem in target pods (simplified)
		fmt.Println("Pod network latency simulation would run here...")

	case "pod-cpu-hog":
		fmt.Println("Pod CPU hog simulation would run here...")

	case "pod-memory-hog":
		fmt.Println("Pod memory hog simulation would run here...")
	}

	// Wait chaos duration
	time.Sleep(time.Duration(chaosDuration) * time.Second)

	// Run k6 load test
	k6 := dagger.Container().
		From("loadimpact/k6:latest").
		WithFile("/load-test.js", dagger.Directory(nil).File("load-test.js")). // replace with your actual script
		WithExec([]string{"k6", "run", "--vus", fmt.Sprintf("%d", loadVUs), "--duration", loadDuration, "/load-test.js"})

	// Collect k6 results (mocked for example)
	output := `{"passed":true,"errorRate":0.0,"p99Latency":120,"recoveryTime":15}`

	if cleanup {
		fmt.Println("Cleanup logic would run here...")
	}

	// Parse JSON
	var result ChaosTestResult
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		return "", err
	}

	// Return JSON string
	resultJSON, _ := json.Marshal(result)
	return string(resultJSON), nil
}

func main() {
	ctx := context.Background()

	// Connect dagger client
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		log.Fatalf("failed to connect Dagger: %v", err)
	}
	defer client.Close()

	ct := &ChaosToolkit{}

	// Example call
	out, err := ct.ChaosTest(
		ctx,
		"default",
		"sample-app",
		"pod-delete",
		60,
		"5m",
		10,
		true,
		client.Host().Directory(".kube"),
		client.Host().Directory(".minikube"),
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Chaos Test Result:", out)
}
