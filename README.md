# docker-lint

Docker-lint is a minimal linter for Dockerfiles. It parses a Dockerfile using the BuildKit parser, normalizes stages into an intermediate representation, and evaluates a set of rules to report potential issues.

## Installation

```bash
go install github.com/asymmetric-effort/docker-lint/cmd/docker-lint@latest
```

## Usage

Provide the path to a Dockerfile. Findings are emitted as a JSON array to standard output.

```bash
docker-lint /path/to/Dockerfile
```

Example output:

```json
[
  {
    "rule": "DL3007",
    "message": "Avoid implicit latest tag",
    "line": 1
  }
]
```

## Development

Common tasks can be run using [`make`](Makefile):

```bash
make clean   # Remove build artifacts
make lint    # Run static analysis
make test    # Run unit and integration tests
make build   # Build the docker-lint binary
```

## License

This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.

