# Golinks Server

[![MIT license](https://img.shields.io/badge/License-MIT-blue.svg)](https://github.com/haostudio/golinks/blob/master/LICENSE)
![go-build](https://github.com/haostudio/golinks/workflows/go-build/badge.svg)
![docker-master](https://github.com/haostudio/golinks/workflows/docker-master/badge.svg)
![docker-release](https://github.com/haostudio/golinks/workflows/docker-release/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/haostudio/golinks)](https://goreportcard.com/report/github.com/haostudio/golinks)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/haostudio/golinks)

![index](./images/index.png?raw=true "Index")

## Usage

### Run with docker

```
$ docker run -v \
  /path/to/datadir:/opt/golinks/datadir \
  -p 8000:8000 \
  haostudio/golinks
```

### Run in `NoAuth` mode

By default, `golinks` supports multiple organizations with HTTP basic
authentication. The links of different organizations are stored in different
namespaces. Running `golinks` in `NoAuth` mode disables HTTP basic
authentication and stores all the links in a shared namespace of `default_org`.

```
$ docker run -v \
  /path/to/datadir:/opt/golinks/datadir \
  -p 8000:8000 \
  -e AUTHPROVIDER_NOAUTH_ENABLED=true \
  -e AUTHPROVIDER_NOAUTH_DEFAULTORG=my_org \ # optional
  haostudio/golinks
```

### Run with source code

#### Build binary

```
$ git clone https://github.com/haostudio/golinks
$ make deps
$ make golinks
```

#### Run

```
$ ./build/golinks
```

## Configuration

golinks uses [conf](https://github.com/popodid/conf) for configurations
and currently supports two sources, environment variables and `golinks_config.yaml` file.

For more advanced configurations, checkout the [conf](https://github.com/popodid/conf),
the [config struct](https://github.com/haostudio/golinks/blob/master/cmd/golinks/main.go#L32),
and the [sample config](https://github.com/haostudio/golinks/blob/master/configs/local.yaml).

### Useful Configurations

| Environment Variable /Yaml Path                                         | Type   | Default                             | Usage                                         |
| ----------------------------------------------------------------------- | ------ | ----------------------------------- | --------------------------------------------- |
| `PORT` / `Port`                                                         | int    | `8000`                              | Listening port                                |
| `LOG_LEVEL` / `Log.Level`                                               | int    | `6`                                 | Maximum log level (`1`~`6`)                   |
| `LOG_STDOUT_ENABLED` / `Log.Stdout.Enabled`                             | bool   | `true`                              | Log to stdout                                 |
| `LOG_STDOUT_WITHCOLOR` / `Log.Stdout.Enabled`                           | bool   | `true`                              | Log to stdout with color                      |
| `METRICS_SAMPLERATE` / `Metrics.SampleRate`                             | float  | `1`                                 | Tracing/Metrics sampling rate                 |
| `METRICS_JAEGER_ENABLED` / `Metrics.Jaeger.Enabled`                     | bool   | `false`                             | Enable tracing with jaeger                    |
| `METRICS_JAEGER_AGENTENDPOINT` / `Metrics.Jaeger.AgentEndpoint`         | string | `localhost:6831`                    | jaeger agent endpoint                         |
| `METRICS_JAEGER_COLLECTORENDPOINT` / `Metrics.Jaeger.CollectorEndpoint` | string | `http://localhost:14268/api/traces` | jaeger collector endpoint                     |
| `METRICS_JAEGER_ENABLED` / `Metrics.Jaeger.Enabled`                     | bool   | `false`                             | Enable tracing with jaeger                    |
| `AUTHPROVIDER_NOAUTH_ENABLED` / `AuthProvider.NoAuth.Enabled`           | bool   | `false`                             | Run in NoAuth mode                            |
| `AUTHPROVIDER_NOAUTH_DEFAULTORG` / `AuthProvider.NoAuth.DefaultOrg`     | string | `_no_org_`                          | The default org namespace used in NoAuth mode |

## Demo

#### Edit link

![edit_link](./images/edit_link.png?raw=true "Edit Link")

#### Register organization

![create_org](./images/create_org.png?raw=true "Create Org")

#### Add user

![create_user](./images/create_user.png?raw=true "Create User")
