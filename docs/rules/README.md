# Lint Rules

The following Hadolint-compatible rules are implemented:

- [DL1001](DL1001.md) - Avoid inline ignore pragmas.

- [DL3000](DL3000.md) - Use absolute WORKDIR.
- [DL3001](DL3001.md) - Avoid irrelevant shell commands like `ssh` or `vim`.
- [DL3002](DL3002.md) - Last USER should not be root.
- [DL3007](DL3007.md) - Avoid using implicit or `latest` tags.
- [DL3008](DL3008.md) - Pin versions in apt-get install.
- [DL3009](DL3009.md) - Delete the APT lists after installing packages.

- [DL3010](DL3010.md) - Use ADD for extracting archives into an image.
- [DL3013](DL3013.md) - Pin versions in pip.
- [DL3014](DL3014.md) - Use the -y switch for apt-get install.
- [DL3015](DL3015.md) - Require `--no-install-recommends` with apt-get install.
- [DL3018](DL3018.md) - Pin versions in apk add.
- [DL3019](DL3019.md) - Use --no-cache with apk add.

- [DL3020](DL3020.md) - Use COPY instead of ADD for files and folders.
- [DL3021](DL3021.md) - COPY with more than 2 arguments requires the last argument to end with /.

- [DL3040](DL3040.md) - dnf clean all missing after dnf command.
- [DL3041](DL3041.md) - Avoid dnf upgrade or update in Dockerfiles.


- [DL3042](DL3042.md) - Combine consecutive RUN instructions that use the same package manager.
- [DL3043](DL3043.md) - Specify OS version tag for base images.
- [DL3044](DL3044.md) - Specify version with dnf/microdnf install.

- [DL3046](DL3046.md) - Avoid apk upgrade in Dockerfiles.
- [DL3047](DL3047.md) - Clean apk cache after installing packages.
- [DL3048](DL3048.md) - Invalid Label Key


- [DL3050](DL3050.md) - Superfluous label(s) present.

- [DL3060](DL3060.md) - `yarn cache clean` missing after `yarn install`.


- [DL4000](DL4000.md) - `MAINTAINER` is deprecated. Use `LABEL maintainer` instead.
- [DL4001](DL4001.md) - Either use Wget or Curl but not both.
- [DL4003](DL4003.md) - Multiple CMD instructions found. Only the last CMD takes effect.
- [DL4004](DL4004.md) - Avoid multiple ENTRYPOINT instructions.
- [DL4005](DL4005.md) - Use SHELL to change the default shell.

