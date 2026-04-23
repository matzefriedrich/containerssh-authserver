[![CI](https://github.com/matzefriedrich/containerssh-authserver/actions/workflows/go.yml/badge.svg)](https://github.com/matzefriedrich/containerssh-authserver/actions/workflows/go.yml)
[![Docker Image CI](https://github.com/matzefriedrich/containerssh-authserver/actions/workflows/release.yml/badge.svg)](https://github.com/matzefriedrich/containerssh-authserver/actions/workflows/release.yml)
[![Coverage Status](https://coveralls.io/repos/github/matzefriedrich/containerssh-authserver/badge.svg?branch=main)](https://coveralls.io/github/matzefriedrich/containerssh-authserver?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/matzefriedrich/containerssh-authserver)](https://goreportcard.com/report/github.com/matzefriedrich/containerssh-authserver)
![License](https://img.shields.io/github/license/matzefriedrich/containerssh-authserver)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/matzefriedrich/containerssh-authserver)
![GitHub Release](https://img.shields.io/github/v/release/matzefriedrich/containerssh-authserver?include_prereleases)

# containerssh-authserver

**containerssh-authserver** is a configurable authentication server application designed to work with [ContainerSSH](https://containerssh.io/) as a webhook backend, by implementing the [ContainerSSH authentication API](https://containerssh.io/v0.5/reference/api/authconfig) . It allows user-specific Docker container profiles to be defined in a simple YAML configuration file, enabling per-user images, shell commands, bind-mounts and network connections.


## Prerequisites

* Go 1.26 or newer (for building from source)
* Docker 20.10+ (or compatible), with Docker Compose (for the demo), Docker API version 1.41 
* openssl, ssh-keygen

## Quick start

The repository comes with an example Docker Compose stack that needs little configuration. Run the following commands to generate the required key material; generated key files are stored in the `keys` directory:

```sh
docker/generate-keys.sh
```

### Assign an SSH public key to a user

Find the public key for the `johndoe` demo account in `keys/johndoe.pem.pub` and add it to the `publicKeys` list of the `johndoe` user in the authserver configuration file in `docker/services/authserver/config.yaml`. Of couse, you can add more users at will.

### Configure per-user containers

The `johndoe` demo account uses the `alpine:3.21` image, without any additional bind mounts or network connections.

### Understand the webhook backend configuration

The `containerssh` service comes with a minimal configuration file (see `docker/services/containerssh/config.yaml` that defines the listening port, the backend URLs for the authentication webhook, and the per-user container configuration.

## Start the demo stack

Use the following command to build and run `containerssh-authserver` in conjunction with `containerssh`:

```sh
docker compose -f docker/docker-compose.yml up --build
```

Once started, you can connect to `containerssh` as `johndoe` using the generated private key and get a shell to a container as configured, for instance:

```sh
ssh -i docker/keys/johndoe.pem -p 2222 johndoe@localhost 
```

### SSH: Too Many Authentication Failures

If you have many keys loaded in `ssh-agent`, SSH may attempt to authenticate with all of them before using the key you specify with `-i`. ContainerSSH limits the number of authentication attempts, which can cause the connection to fail before the correct key is tried.

To prevent SSH from offering all agent keys, use the `IdentitiesOnly=yes` option:

```bash
ssh -o IdentitiesOnly=yes -i docker/keys/johndoe.pem -p 2222 johndoe@localhost
```

This tells SSH to use **only the explicitly specified identity file** and ignore any keys loaded in `ssh-agent`.

### Generate a bcrypt password (for password authentication)

Example: Run the following command to generate a bcrypt password hash for the `johndoe` user:

```sh
docker run --rm httpd:2.4 htpasswd -Bbn johndoe topsecret
```

See `docker/services/authserver/config.yaml` file; find the user configuration and update the `secret` property.
