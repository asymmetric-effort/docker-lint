# docker-lint

Docker-lint is a minimal linter for Dockerfiles. It parses a Dockerfile using the BuildKit parser, normalizes stages into an intermediate representation, and evaluates a set of rules to report potential issues.

## Installation

```bash
go install github.com/asymmetric-effort/docker-lint/cmd/docker-lint@latest
```

## Usage

Provide one or more paths or glob patterns (supports `*` and `**`). Matching files are linted and findings are emitted as a JSON array to standard output.

```bash
docker-lint /path/to/Dockerfile
docker-lint './**/Dockerfile'
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

## Configuration

docker-lint reads an optional `.docker-lint.yaml` from the current working directory to configure file-based rule exclusions.

```yaml
exclude:
  Dockerfile.dev:
    - DL3007
```

The above disables rule `DL3007` for `Dockerfile.dev`.

## Linting Containers

docker-lint is also published as a container image. This allows you to lint Dockerfiles without installing the binary on your host system. Mount your project directory and provide the Dockerfile path inside the container:

```bash
docker run --rm -v "$(pwd):/src" ghcr.io/asymmetric-effort/docker-lint:latest /src/Dockerfile
```

## Rules

- [DL3050](docs/rules/DL3050.md) - Superfluous label(s) present when strict label validation is enabled.

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

