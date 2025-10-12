package main

import (
	"context"
	"log"
	"os"
	"os/exec"

	"dagger.io/dagger"
)

func main() {
	ctx := context.Background()
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		log.Fatalf("Failed to connect Dagger: %v", err)
	}
	defer client.Close()

	// 1️⃣ Apply chaos experiments
	litmusFolder := "./manifest/litmus"
	experiments := []string{"cpu-hog.yaml", "memory-hog.yaml", "network-latency.yaml", "pod-delete.yaml"}

	namespace := os.Getenv("NAMESPACE")
	if namespace == "" {
		namespace = "default"
	}

	for _, exp := range experiments {
		expPath := litmusFolder + "/" + exp
		cmd := exec.Command("kubectl", "apply", "-f", expPath, "-n", namespace)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Fatalf("Failed to apply chaos experiment %s: %v", exp, err)
		}
		log.Printf("Applied chaos experiment: %s\n", exp)
	}

	// 2️⃣ Run k6 tests
	k6Folder := "./k6"
	k6Container := client.Container().
		From("loadimpact/k6:latest").
		WithMountedDirectory("/tests", client.Host().Directory(k6Folder, dagger.HostDirectoryOpts{})).
		WithWorkdir("/tests")

	serviceURL := os.Getenv("SERVICE_URL")
	if serviceURL == "" {
		serviceURL = "http://sample-app.default.svc.cluster.local"
	}

	vus := os.Getenv("VUS")
	if vus == "" {
		vus = "10"
	}
	duration := os.Getenv("DURATION")
	if duration == "" {
		duration = "5m"
	}

	k6Container = k6Container.WithExec([]string{"run", "--vus", vus, "--duration", duration, "test.js"})
	output, err := k6Container.Stdout(ctx)
	if err != nil {
		log.Fatalf("Failed to run k6: %v", err)
	}

	log.Println("=== K6 Test Output ===")
	log.Println(output)

	// 3️⃣ Optional cleanup
	for _, exp := range experiments {
		expPath := litmusFolder + "/" + exp
		cmd := exec.Command("kubectl", "delete", "-f", expPath, "-n", namespace)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Printf("Warning: Failed to delete chaos experiment %s: %v", exp, err)
		} else {
			log.Printf("Deleted chaos experiment: %s\n", exp)
		}
	}
}
