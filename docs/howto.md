# How To Use docker-lint

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

To display the current version:

```bash
docker-lint --version
docker-lint -version
docker-lint version
```
