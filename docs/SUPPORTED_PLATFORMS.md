<!-- file: docs/SUPPORTED_PLATFORMS.md -->
<!-- (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com -->

# Supported Platforms

The following GOOS/GOARCH combinations are verified in continuous integration:

| GOOS    | GOARCH |
|---------|--------|
| linux   | amd64  |
| linux   | arm64  |
| windows | amd64  |
| windows | arm64  |
| darwin  | amd64  |
| darwin  | arm64  |

> **Note**
> The Go toolchain does not currently provide support for `arm64be`, so builds for that architecture are unavailable.
