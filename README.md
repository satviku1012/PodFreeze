# PodFreeze

PodFreeze is an initial scaffolding project for a CRIU-based hibernation controller for idle Kubernetes workloads. The current milestone focuses on the orchestration contract, a basic gRPC daemon, and a buffering reverse proxy that pauses incoming requests while the container is restored.

## Quick start

```bash
make deps
make proto
make build-proxy
make build-daemon
```

## Running the components

Start the daemon:

```bash
./bin/daemon
```

In a second terminal, start the proxy:

```bash
./bin/proxy
```

Then send a request with a container identifier and target address:

```bash
curl -H 'X-Container-ID: demo' -H 'X-Target-IP: 127.0.0.1:8080' http://127.0.0.1:8081/hello
```
