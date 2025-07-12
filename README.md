# Simple Temporal Workflow

A basic Go application demonstrating Temporal workflow patterns with activities.

## Overview

This project implements a simple Temporal workflow that executes two sequential activities: processing and validation.

## Files

- `workflow.go` - Contains the main workflow definition
- `activities.go` - Implements the workflow activities
- `worker.go` - Sets up and starts the Temporal worker
- `go.mod` - Go module dependencies

## Requirements

- Go 1.21+
- Temporal server running locally

## Usage

1. Start the Temporal server
2. Run the worker: `go run .`
3. Execute workflows through the Temporal CLI or web UI

## Dependencies

- `go.temporal.io/sdk` - Temporal Go SDK