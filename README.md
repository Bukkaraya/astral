# Astral

A Rails-like framework for Temporal workflows in Go. Build robust distributed applications with automatic code generation, type-safe clients, and convention-over-configuration patterns.

## Overview

Astral provides Rails-like conventions for building Temporal applications, including automatic client generation, structured project layouts, and developer-friendly tooling.

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