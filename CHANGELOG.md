# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [v0.2.1] - 2025-09-09

- Bumps `github.com/gofiber/fiber/v2` from 2.52.8 to 2.52.9 [#8](https://github.com/matzefriedrich/containerssh-authserver/pull/8)
- Bumps `golang.org/x/crypto` from 0.39.0 to 0.41.0 [#10](https://github.com/matzefriedrich/containerssh-authserver/pull/10)


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
