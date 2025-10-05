package main

import (
    "context"
    "fmt"
    "os"
    "os/exec"
    "strings"

    "dagger.io/dagger"
)

func main() {
    ctx := context.Background()

    client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
    if err != nil {
        panic(err)
    }
    defer client.Close()

    fmt.Println("✅ Dagger pipeline running...")

    changedServices, err := getChangedServices()
    if err != nil {
        panic(err)
    }

    if len(changedServices) == 0 {
        fmt.Println("No changed services detected — nothing to build.")
        return
    }

    for _, service := range changedServices {
        fmt.Println("🔨 Building service:", service)

        serviceBin := buildService(ctx, client, service)

        outputPath := fmt.Sprintf("./%s-built", service)

        _, err := serviceBin.Export(ctx, outputPath)
        if err != nil {
            panic(fmt.Sprintf("Failed to export %s binary: %v", service, err))
        }

        fmt.Printf("✅ %s built and exported to %s\n", service, outputPath)
    }
}

func buildService(ctx context.Context, client *dagger.Client, serviceName string) *dagger.File {
    return client.Container().
        From("golang:1.20").
        WithMountedDirectory("/src", client.Host().Directory("../services/" + serviceName)).
        WithWorkdir("/src").
        WithExec([]string{"go", "build", "-o", "app"}).
        File("app")
}

// getChangedServices uses `git diff` to find which services have changed
func getChangedServices() ([]string, error) {
    // Run git diff to compare current branch with origin/main, list only filenames changed
    cmd := exec.Command("git", "diff", "--name-only", "origin/main...")

    // Capture the output of the git command
    out, err := cmd.Output()
    if err != nil {
        return nil, err // Return error if git command fails
    }

    // Split output string by newline to get a list of changed files
    files := strings.Split(string(out), "\n")

    // Create a map to store unique service names found from changed files
    serviceSet := map[string]bool{}

    // Iterate over each changed file path
    for _, file := range files {
        // Check if file path starts with "services/", meaning it belongs to a service
        if strings.HasPrefix(file, "services/") {
            // Split the file path by '/' to isolate parts
            parts := strings.Split(file, "/")
            // If valid path with service name (second part), add to map
            if len(parts) >= 2 {
                serviceSet[parts[1]] = true
            }
        }
    }

    // Convert the keys of the map into a slice of service names to return
    services := []string{}
    for svc := range serviceSet {
        services = append(services, svc)
    }

    return services, nil // Return list of changed services, no error
}
