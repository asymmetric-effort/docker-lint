# Configuration

Docker-lint accepts configuration files in the same format as [hadolint](https://github.com/hadolint/hadolint). Place a
`.docker-lint.yaml` file in the project root to adjust rule behavior.

## Example

```yaml
ignored:
  - DL3006
  - DL3008
  - DL3026
  - DL3050
  - SC3050
  - SC3020
override:
  warning:
    - SC1099
failure-threshold: warning
trustedRegistries:
  - ghcr.io
strict-labels: true
label-schema:
  author: text
  base-image: text
  contact: text
  created: text
  documentation: text
  git-commit: text
  license: text
  version: text
```

See the [hadolint documentation](https://github.com/hadolint/hadolint#configure) for the meaning of these fields. Docker-lint
currently honours the `ignored` list and parses the remaining fields for forward compatibility.
