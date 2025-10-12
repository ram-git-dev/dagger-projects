package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"dagger.io/dagger"
)

type ChaosToolkit struct {
	Client *dagger.Client
}

func NewChaosToolkit(client *dagger.Client) *ChaosToolkit {
	return &ChaosToolkit{Client: client}
}

// Run actual chaos strategies
func (c *ChaosToolkit) ChaosTest(ctx context.Context, kubeconfigDir, minikubeDir *dagger.Directory) error {
	fmt.Println("Running chaos strategies...")

	// Example: Pod delete simulation
	podDelete := c.Client.Container().From("bitnami/kubectl:latest").
		WithMountedDirectory("/kube", kubeconfigDir).
		WithExec([]string{"kubectl", "delete", "pod", "-n", "ingress-nginx", "--all"})
	if _, err := podDelete.ExitCode(ctx); err != nil {
		return fmt.Errorf("pod delete failed: %w", err)
	}
	fmt.Println("Pods deleted successfully")

	// Example: CPU spike simulation (dummy sleep)
	fmt.Println("Simulating CPU spike...")
	time.Sleep(2 * time.Second)

	// Example: Network delay simulation (dummy sleep)
	fmt.Println("Simulating network delay...")
	time.Sleep(2 * time.Second)

	// Run k6 load test
	fmt.Println("Running k6 load test...")
	k6 := c.Client.Container().From("loadimpact/k6:latest").
		WithMountedDirectory("/tests", minikubeDir). // place your k6 script here
		WithExec([]string{"run", "/tests/test.js"})
	if _, err := k6.ExitCode(ctx); err != nil {
		return fmt.Errorf("k6 test failed: %w", err)
	}

	return nil
}

func main() {
	ctx := context.Background()

	client, err := dagger.Connect(ctx, dagger.WithLogOutput(nil))
	if err != nil {
		log.Fatalf("Failed to connect to Dagger: %v", err)
	}
	defer client.Close()

	fmt.Println("Connected to Dagger engine")

	kubeconfigDir := client.Host().Directory("/home/r/.kube")
	minikubeDir := client.Host().Directory("/home/r/.minikube")

	ct := NewChaosToolkit(client)

	if err := ct.ChaosTest(ctx, kubeconfigDir, minikubeDir); err != nil {
		log.Fatalf("Chaos test failed: %v", err)
	}

	fmt.Println("Chaos test completed successfully")
}
