# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [v0.2.5] - 2025-11-17

* Bumps `golang.org/x/crypto` from 0.43.0 to 0.44.0 [#27](https://github.com/matzefriedrich/containerssh-authserver/pull/27)
* Updates the Go version to 1.25.4 [#29](https://github.com/matzefriedrich/containerssh-authserver/pull/29)
* Bumps indirect dependencies to their latest versions [#30](https://github.com/matzefriedrich/containerssh-authserver/pull/30)


## [v0.2.4] - 2025-11-16

* Bumps `github.com/docker/docker` from 28.5.1+incompatible to 28.5.2+incompatible [#26](https://github.com/matzefriedrich/containerssh-authserver/pull/26)


## [v0.2.3] - 2025-11-04

* Bumps `github.com/spf13/viper` from 1.20.1 to 1.21.0 [#19](https://github.com/matzefriedrich/containerssh-authserver/pull/19)
* Bumps `golang.org/x/crypto` from 0.41.0 to 0.43.0 [#22](https://github.com/matzefriedrich/containerssh-authserver/pull/22/)
* Bumps `github.com/docker/docker` from 28.4.0+incompatible to 28.5.1+incompatible [#23](https://github.com/matzefriedrich/containerssh-authserver/pull/23)
* Bumps `github.com/matzefriedrich/parsley` from 1.1.4 to 1.3.0 [#25](https://github.com/matzefriedrich/containerssh-authserver/pull/25)


## [v0.2.2] - 2025-09-09

- Bumps `github.com/gofiber/fiber/v2` from 2.52.8 to 2.52.9 [#8](https://github.com/matzefriedrich/containerssh-authserver/pull/8)
- Bumps `golang.org/x/crypto` from 0.39.0 to 0.41.0 [#10](https://github.com/matzefriedrich/containerssh-authserver/pull/10)
- Updates the Go version from 1.24.0 to 1.25
- Bumps `golang` from 1.24-alpine to 1.25-alpine [#14](https://github.com/matzefriedrich/containerssh-authserver/pull/14)
- Bumps `github.com/stretchr/testify` from 1.10.0 to 1.11.1 [#15](https://github.com/matzefriedrich/containerssh-authserver/pull/15)
- Bumps `github.com/docker/docker` from 28.2.2+incompatible to 28.4.0+incompatible [#16](https://github.com/matzefriedrich/containerssh-authserver/pull/16)
- Bumps `github.com/matzefriedrich/parsley` from 1.0.13 to 1.1.3 [#17](https://github.com/matzefriedrich/containerssh-authserver/pull/17)


## [v0.2.0] - 2025-06-07

### Fixes

- Bumps `github.com/docker/docker` from version `20.10` to `28.2`, resolving known CVEs

- Removed the `containerssh@0.5.2` package reference, and introduced type shims to ensure compatibility between ContainerSSH and the latest Docker SDK types.


## [v0.1.1] - 2025-06-06

### Added

- **API service**: Implements the ContainerSSH authentication and configuration server API using `gofiber/fiber/v2` and `matzefriedrich/parsley`, with public key authentication support

- **Docker backend**: Supports per-user container configurations, including image selection, shell commands, bind mounts, and network settings

- **Configuration**: User and container profile configuration via YAML files

- **Deployment**: Sample Docker Compose stack
