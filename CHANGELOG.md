# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [v0.3.2] - 2026-04-22

### Changed

* Upgrades Go version to 1.26.2 [#46](https://github.com/matzefriedrich/containerssh-authserver/pull/46)
* Bumps `github.com/gofiber/contrib/v3/zerolog` from v1.0.1 to v1.0.2 [#46](https://github.com/matzefriedrich/containerssh-authserver/pull/46)
* Bumps `github.com/matzefriedrich/parsley` from v1.3.2 to v1.3.3 [#46](https://github.com/matzefriedrich/containerssh-authserver/pull/46)
* Bumps `github.com/rs/zerolog` from v1.35.0 to v1.35.1 [#46](https://github.com/matzefriedrich/containerssh-authserver/pull/46)
* Bumps indirect dependencies to their latest versions. [#46](https://github.com/matzefriedrich/containerssh-authserver/pull/46)


## [v0.3.1] - 2026-04-13

### Added

* Added **`DOCKERHUB.md`** documentation for Docker image usage and configuration.

### Changed

* Upgraded the **ContainerSSH** image version to **v0.6** and updated volume mount permissions for the Docker socket in the sample Docker Compose configuration.
* Bumps `golang.org/x/crypto` from 0.49.0 to 0.50.0.
* Bumps **indirect dependencies** to their latest versions.

### Fixed

* Improved **public-key authentication test** by including a required `Host` header.
* Ignores configuration loading error in `configuration_module.go`.


## [v0.3.0] - 2026-04-12

### Changed

* Upgraded the HTTP server stack to **Fiber v3** and updated request handling to match the new API. [#42](https://github.com/matzefriedrich/containerssh-authserver/pull/42)
* Added **graceful shutdown support** to server startup. [#42](https://github.com/matzefriedrich/containerssh-authserver/pull/42)
* Bumps `github.com/rs/zerolog` from 1.34.0 to 1.35.0 [#43](https://github.com/matzefriedrich/containerssh-authserver/pull/43)


### Fixed

* Improved **SSH public-key authentication** error handling by surfacing profile lookup, key parsing, and key mismatch failures more clearly. [#42](https://github.com/matzefriedrich/containerssh-authserver/pull/42)
* Enhanced authentication and profile-loading **structured logging** for better observability during SSH requests. [#42](https://github.com/matzefriedrich/containerssh-authserver/pull/42)


## [v0.2.7] - 2026-03-12

### Changed

* Bumps `github.com/matzefriedrich/parsley` from 1.3.0 to 1.3.2 [#37](https://github.com/matzefriedrich/containerssh-authserver/pull/37)
* Bumps `golang.org/x/crypto` from 0.46.0 to 0.48.0 [#38](https://github.com/matzefriedrich/containerssh-authserver/pull/38)
* Bumps `golang` from 1.25-alpine to 1.26-alpine [#39](https://github.com/matzefriedrich/containerssh-authserver/pull/39)
* Bumps `github.com/gofiber/fiber/v2` from 2.52.10 to 2.52.12 [#40](https://github.com/matzefriedrich/containerssh-authserver/pull/40)
* Bumps Go version from 1.25.7 to 1.26.1 [#41](https://github.com/matzefriedrich/containerssh-authserver/pull/41)


## [v0.2.6] - 2025-01-08

### Changed

* Bumps `golang.org/x/crypto` from 0.44.0 to 0.45.0 [#31](https://github.com/matzefriedrich/containerssh-authserver/pull/31)
* Bumps `github.com/gofiber/fiber/v2` from 2.52.9 to 2.52.10 [#32](https://github.com/matzefriedrich/containerssh-authserver/pull/32)
* Upgrades Go version to `1.25.5`
* Bumps `golang.org/x/crypto` from 0.45.0 to 0.46.0


## [v0.2.5] - 2025-11-17

### Changed

* Bumps `golang.org/x/crypto` from 0.43.0 to 0.44.0 [#27](https://github.com/matzefriedrich/containerssh-authserver/pull/27)
* Updates the Go version to 1.25.4 [#29](https://github.com/matzefriedrich/containerssh-authserver/pull/29)
* Bumps indirect dependencies to their latest versions [#30](https://github.com/matzefriedrich/containerssh-authserver/pull/30)


## [v0.2.4] - 2025-11-16

### Changed

* Bumps `github.com/docker/docker` from 28.5.1+incompatible to 28.5.2+incompatible [#26](https://github.com/matzefriedrich/containerssh-authserver/pull/26)


## [v0.2.3] - 2025-11-04

### Changed

* Bumps `github.com/spf13/viper` from 1.20.1 to 1.21.0 [#19](https://github.com/matzefriedrich/containerssh-authserver/pull/19)
* Bumps `golang.org/x/crypto` from 0.41.0 to 0.43.0 [#22](https://github.com/matzefriedrich/containerssh-authserver/pull/22/)
* Bumps `github.com/docker/docker` from 28.4.0+incompatible to 28.5.1+incompatible [#23](https://github.com/matzefriedrich/containerssh-authserver/pull/23)
* Bumps `github.com/matzefriedrich/parsley` from 1.1.4 to 1.3.0 [#25](https://github.com/matzefriedrich/containerssh-authserver/pull/25)


## [v0.2.2] - 2025-09-09

### Changed

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
