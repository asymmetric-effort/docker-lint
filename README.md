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

## Development

Common tasks are defined in the [Taskfile](Taskfile.yml). To run them manually:

```bash
go mod tidy
go test ./... -short -cover
go test ./... -run Integration -covermode=atomic -coverpkg=./... -coverprofile=coverage.out
go tool cover -func=coverage.out | awk '/total:/ { print; if ($3+0 < 80) exit 1 }'
go build -trimpath -ldflags "-s -w" ./cmd/docker-lint
```

## License

This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.

