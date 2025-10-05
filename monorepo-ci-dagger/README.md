# Multi-Service Build Pipeline with Dagger

## What It Does
This project automatically builds only the Go services that have changed compared to `origin/main`. It uses Dagger to run builds inside containers, ensuring clean, consistent environments.

## Setup
1. Install Go (v1.20+)
2. Clone this repo and navigate to it
3. Ensure Dagger is installed and running (see [Dagger docs](https://dagger.io)

## How to Run Locally
From the repo root, run:

go run main.go

This will:

    Detect changed services using git diff against origin/main

    Build those services inside containerized Go environments

    Export built binaries to ./{service}-built

GitHub Actions CI

    Runs on pushes to feature branches (excluding main)

    Detects changed services and builds them

    Ensures only updated services are built and tested before merging to main

How Changed Services Are Detected

    Uses git diff --name-only origin/main...

    Looks for changes under services/{serviceName}/

    Builds only services with detected changes

Adding New Services

    Add your service under services/{new-service}

    Make sure it contains a valid Go project that builds with go build

Troubleshooting

    Ensure your git remote is up to date (run git fetch origin)

    Confirm Dagger daemon is running (dagger daemon or similar)

    Check logs for build failures and missing dependencies