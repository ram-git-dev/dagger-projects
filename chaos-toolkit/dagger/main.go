package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"dagger.io/dagger"
)

// ChaosToolkit holds your Dagger client and directories
type ChaosToolkit struct {
	Client       *dagger.Client
	Kubeconfig   *dagger.Directory
	MinikubeDir  *dagger.Directory
}

// Chaos modes
func (ct *ChaosToolkit) KillPod(ctx context.Context, podName string) error {
	fmt.Printf("[Chaos] Killing pod: %s\n", podName)
	// Example: simulate pod deletion
	time.Sleep(1 * time.Second)
	fmt.Println("[Chaos] Pod killed")
	return nil
}

func (ct *ChaosToolkit) StressCPU(ctx context.Context, duration time.Duration) error {
	fmt.Printf("[Chaos] Stressing CPU for %s\n", duration)
	time.Sleep(duration)
	fmt.Println("[Chaos] CPU stress finished")
	return nil
}

// K6 load test
func (ct *ChaosToolkit) RunK6(ctx context.Context, scriptPath string) error {
	fmt.Printf("[K6] Running load test: %s\n", scriptPath)
	// simulate k6 run
	time.Sleep(2 * time.Second)
	fmt.Println("[K6] Load test finished")
	return nil
}

func main() {
	ctx := context.Background()

	// Init Dagger client
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(log.Writer()))
	if err != nil {
		log.Fatalf("Failed to connect to Dagger: %v", err)
	}
	defer client.Close()

	ct := &ChaosToolkit{
		Client:      client,
		Kubeconfig:  client.Host().Directory("/home/user/.kube"),   // adjust path
		MinikubeDir: client.Host().Directory("/home/user/.minikube"), // adjust path
	}

	// Run your chaos strategies + k6 sequentially
	if err := ct.KillPod(ctx, "my-app-pod"); err != nil {
		log.Fatalf("Chaos failed: %v", err)
	}

	if err := ct.StressCPU(ctx, 5*time.Second); err != nil {
		log.Fatalf("Chaos failed: %v", err)
	}

	if err := ct.RunK6(ctx, "./k6-scripts/test.js"); err != nil {
		log.Fatalf("K6 test failed: %v", err)
	}

	fmt.Println("All chaos + k6 steps completed successfully ✅")
}
