# Golang Service Start to Production

### Overview

Create a robust Golang service, deploy it using Docker and Kubernetes, and enhance it with logging, configuration metrics, (continue ............)

## Table of Contents

- [Introduction](#introduction)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
- [Usage](#usage)
  - [Running the Application](#running-the-application)
  - [Kubernetes Cluster](#kubernetes-cluster)
    - [Setting up a Local Kubernetes Cluster](#setting-up-a-local-kubernetes-cluster)
    - [Deploying the Service to Kubernetes](#deploying-the-service-to-kubernetes)
- [Commands](#commands)
- [Module Vendor Support](#module-vendor-support)

## Things Covered

- [x] Documentation using open api docs
- [x] Packaged using Kubernetes
- [x] Dockerized
- [x] Middleware with JWT authentication
- [x] Logging with Zap logger
- [x] Microservice architecture
- [x] Onion Layering
- [x] sqlx used as Database ORM with PostgreSQL
- [x] Unit testing
- [x] Integration testing
- [ ] Opentelemetry with Prometheus, grafana, Loki, etc
- [ ] github workflow

## Introduction

<!-- Provide a brief overview of the project, its purpose, and its key features. -->

## Getting Started

<!-- Explain how to set up the project locally and any prerequisites that need to be installed. -->

### Prerequisites

<!--
List any software or tools that must be installed before running the application or deploying to Kubernetes. -->

- Docker
- Go (Golang)
- Kubernetes (optional for deploying to a local cluster)

### Installation

Instructions on how to clone the repository and install the necessary dependencies.

---

```bash
$ git clone https://github.com/avyukth/service-s2p.git
$ cd project
$ go mod tidy
$ go mod vendor
```

---

## Usage

<!-- Provide instructions on how to use the project and how to interact with it. -->

### Running the Application

To run the application locally, execute the following command:

---

```bash
$ go run main.go
```

---

### Kubernetes Cluster

If you want to deploy the service to a Kubernetes cluster, follow the steps below.

#### Setting up a Local Kubernetes Cluster

Create a local Kubernetes cluster using Kind:

---

```bash
$ make kind-up
```

---

#### Deploying the Service to Kubernetes

Build and deploy the Docker image to the Kubernetes cluster:

---

```bash
$ make kind-update-apply
```

---

## Commands

List and explain the different make commands available for managing the project.

- `make run`: Run the application locally using Go.
- `make kind-up`: Create a local Kubernetes cluster using Kind.
- `make kind-down`: Delete the local Kubernetes cluster.
- `make kind-load`: Load the Docker image into the Kind cluster.
- `make kind-status`: Get status information about the Kind cluster and its resources.
- `make kind-apply`: Apply Kubernetes manifests to the Kind cluster.
- `make kind-status-service`: Get status information about the deployed service pods.
- `make kind-logs`: Tail the logs for the service pods.
- `make kind-restart`: Restart the service pods in the Kind cluster.
- `make kind-update`: Build the Docker image, load it into the Kind cluster, and restart the service pods.
- `make kind-describe`: Get detailed information about the deployed service pods.
- `make kind-update-apply`: Build the Docker image, load it into the Kind cluster, and apply Kubernetes manifests.

## Module Vendor Support

This section is dedicated to the support for Go modules vendor directory.

- `make tidy`: Perform `go mod tidy` and `go mod vendor` to clean up the Go module and vendor dependencies.
