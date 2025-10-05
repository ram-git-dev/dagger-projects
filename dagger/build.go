package main

import (
    "context"
    "fmt"
    "os"

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

    // Build service-a
    serviceABin := buildService(ctx, client, "service-a")
    fmt.Println("✅ service-a built and exported to ./service-a-built")

    // Build service-b
    serviceBBin := buildService(ctx, client, "service-b")
    fmt.Println("✅ service-b built and exported to ./service-b-built")

    // Export built binaries locally
    _, err = serviceABin.Export(ctx, "./service-a-built")
    if err != nil {
        panic(err)
    }
    _, err = serviceBBin.Export(ctx, "./service-b-built")
    if err != nil {
        panic(err)
    }
}

// buildService builds the given service folder and returns the output File
func buildService(ctx context.Context, client *dagger.Client, serviceName string) *dagger.File {
    output := client.Container().
        From("golang:1.20").
        WithMountedDirectory("/src", client.Host().Directory("./services/" + serviceName)).
        WithWorkdir("/src").
        WithExec([]string{"go", "build", "-o", "app"}).
        File("app")

    return output
}
