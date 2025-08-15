# Hadolint Rule DL3049 — Required label is missing

## 1. Purpose
This rule ensures Docker images include metadata labels that our policy requires. It helps track provenance, link SBOMs, and provide essential information like source repository or version.

## 2. Scope
DL3049 checks for the presence of required label keys on every stage in a Dockerfile when docker-lint is run with `--require-label` or the equivalent configuration.

## 3. Rule Statement (normative)
Emit DL3049 when a required label is not defined in the relevant build context.

**Message:** `Label "<key>" is missing.`

The rule only evaluates label keys explicitly configured via `--require-label` or in `.hadolint.yaml`.

## 4. Rationale
Required labels enforce minimum metadata for policy compliance, provenance, and discoverability in registries.

## 5. Configuration
Declare required labels either on the command line or in configuration:

```sh
hadolint --require-label org.opencontainers.image.source \
         --require-label org.opencontainers.image.version Dockerfile
```

Or in `.hadolint.yaml`:

```yaml
require-label:
  - org.opencontainers.image.source
  - org.opencontainers.image.version
  - org.opencontainers.image.revision:"[0-9a-f]{7,40}"
```

## 6. Multi-stage semantics (normative)
Hadolint tracks labels per stage:
* Final stage must define all required labels.
* Builder stages with aliases referenced by `COPY --from=<alias>` are treated as silent.
* Defining a required label marks a stage as compliant.

## 7. Detection Logic (high-level)
1. Track the current stage for each `FROM` instruction.
2. Mark a stage as compliant when a `LABEL` defines the required key.
3. Mark a stage as silent if it is aliased and only used as a copy source.
4. Emit DL3049 for any non‑silent stage that never defined the required label.

## 8. Message and Severity
* **ID:** `DL3049`
* **Severity:** `info` (configurable)
* **Title:** `Label "<key>" is missing.`
* **Auto-fix:** Not applicable; add a `LABEL` instruction.

## 9. Examples
### Non‑compliant (single stage)
```Dockerfile
FROM alpine:3.19
# Missing required label
```

### Compliant (single stage)
```Dockerfile
FROM alpine:3.19
LABEL org.opencontainers.image.source="https://github.com/acme/project"
```

### Compliant (multi‑stage; labels only on final)
```Dockerfile
FROM golang:1.22 AS build
# build stuff...

FROM alpine:3.19
LABEL org.opencontainers.image.source="https://github.com/acme/project"
LABEL org.opencontainers.image.revision="abcd1234"
COPY --from=build /out/app /usr/local/bin/app
```

## 10. Best Practices (informative)
* Apply all required labels in the final stage.
* Prefer OCI label keys such as:
  * `org.opencontainers.image.title`
  * `org.opencontainers.image.description`
  * `org.opencontainers.image.url`
  * `org.opencontainers.image.source`
  * `org.opencontainers.image.version`
  * `org.opencontainers.image.revision`
  * `org.opencontainers.image.licenses`
* Centralize label values via build args or fragments.

## 11. Known Edge Cases
* Anonymous stages not referenced by name may still be checked; label the final stage.
* Generated Dockerfiles must inject required labels on the final stage only.

## 12. Suppression / Overrides
Use inline suppression sparingly:
```Dockerfile
# hadolint ignore=DL3049
```
Adjust policy by changing `require-label` in configuration instead.

## 13. Remediation Checklist
* [ ] Define required label keys in configuration or CLI flags.
* [ ] Add `LABEL <key>=<value>` to the final stage.
* [ ] Add labels to intermediate stages if they produce standalone images.
* [ ] Validate value format with regex constraints where applicable.

Our repository enforces DL3049 through the docker-lint rule set. When configured with `require-label`, the linter reports missing labels, and CI runs `go test` to verify rule behavior.
