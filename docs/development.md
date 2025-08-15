# Development

Common tasks are managed with [`make`](../Makefile):

| Target | Description |
| ------ | ----------- |
| `all` | Run `clean`, `lint`, `test`, and `build` |
| `bump_version` | Increment the minor version and tag the current commit |
| `bump_version/major` | Increment the major version and tag the current commit |
| `lint` | Run static analysis |
| `test` | Run unit and integration tests |
| `build` | Build the docker-lint binary |
| `tidy` | Update Go module dependencies |
| `clean` | Remove build artifacts |

(c) 2025 Asymmetric Effort, LLC. <scaldwell@asymmetric-effort.com>
[<img src="img/asymmetric-effort.png" alt="Asymmetric Effort logo" width="60" height="60">](https://asymmetric-effort.com/)
