# mobymatze/containerssh-authserver

**containerssh-authserver** is a configurable authentication server application designed to work with [ContainerSSH](https://containerssh.io/) as a webhook backend. It implements the [ContainerSSH authentication API](https://containerssh.io/v0.5/reference/api/authconfig) and allows user-specific Docker container profiles to be defined in a simple YAML configuration file.

With this server, you can enable per-user images, shell commands, bind-mounts, and network connections for your ContainerSSH deployment.

## How to use this image

Pull the image via `docker pull mobymatze/containerssh-authserver:latest`, and run a container; see the following example with a config file mounted at `/var/run/authserver/config.yaml`:

```sh
docker run -d \
  --name authserver \
  -p 5000:5000 \
  -v /path/to/your/config.yaml:/var/run/authserver/config.yaml:ro \
  mobymatze/containerssh-authserver:latest
```

## Configuration

The application can be configured via a YAML configuration file or environment variables.

### Configuration YAML format

By default, the application looks for a `config.yaml` file in `/var/run/authserver/`. You can change this path by setting the `AUTHSERVER_CONFIG_PATH` environment variable.

Example `config.yaml`:

```yaml
app:
  port: 5000
  logLevel: info
  authServer:
    users:
      - johndoe:
          publicKeys:
            - "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAA... user@host"
          image: alpine:3.21
          shellCommand: ["/bin/sh"]
          binds:
            - "/host/path:/container/path:ro"
          networks:
            - "my-network"
```

### Environment variables

Settings can also be overridden using environment variables. Use `__` (double underscore) as a separator for nested keys.

- `APP__PORT`: The port the server listens on (default: `5000`).
- `APP__LOGLEVEL`: Log level (e.g., `debug`, `info`, `warn`, `error`).
- `AUTHSERVER_CONFIG_PATH`: Path to the directory containing `config.yaml`.

## Docker Compose Example

This example shows how to run `containerssh-authserver` alongside ContainerSSH.

```yaml
services:
  authserver:
    image: mobymatze/containerssh-authserver
    restart: on-failure
    environment:
      AUTHSERVER_CONFIG_PATH: "/var/run/authserver"
    volumes:
      - "./config.yaml:/var/run/authserver/config.yaml:ro"
    networks:
      - backend

  containerssh:
    image: containerssh/containerssh:v0.5.2
    ports:
      - "2222:2222"
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock:ro
      - ./containerssh-config.yaml:/etc/containerssh/config.yaml:ro
    networks:
      - backend

networks:
  backend:
    driver: bridge
```

## Source & License

* Source: [https://github.com/matzefriedrich/containerssh-authserver](https://github.com/matzefriedrich/containerssh-authserver)
* License: [Apache License 2.0](https://github.com/matzefriedrich/containerssh-authserver/blob/main/LICENSE).
