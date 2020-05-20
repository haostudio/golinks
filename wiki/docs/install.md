# Installl

## Get started

### Run with docker

```sh
$ docker run -v \
  /path/to/datadir:/opt/golinks/datadir \
  -p 8000:8000 \
  haostudio/golinks
```

!!! Tip "[http://localhost:8000](http://localhost:8000)"

### Run with source code

#### Build binary

```sh
# Clone repo
$ git clone https://github.com/haostudio/golinks
# Install dependencies
$ make deps
# Build binary
$ make wiki golinks
# Run
$ ./build/golinks
```

!!! Note
	Make `wiki` before golinks because we compile the wiki into `golinks` binary,
	as well as the html templates. The `golinks` binary will be everything you
	need to run the server.

## Advanced

### Run in NoAuth mode

By default, `golinks` supports multiple organizations with HTTP basic
authentication. The links of different organizations are stored in different
namespaces. Running `golinks` in **NoAuth** mode disables HTTP basic
authentication and stores all the links in a shared namespace of the **default org**.

```sh
$ docker run -v \
  /path/to/datadir:/opt/golinks/datadir \
  -p 8000:8000 \
  -e AUTHPROVIDER_NOAUTH_ENABLED=true \
  -e AUTHPROVIDER_NOAUTH_DEFAULTORG=my_org \ # optional
  haostudio/golinks
```

### Enable static wiki site

```sh
$ docker run -v \
  /path/to/datadir:/opt/golinks/datadir \
  -p 8000:8000 \
  -e HTTP_GOLINKS_WIKI=true \
  haostudio/golinks
```

!!! Tip "[http://localhost:8000/wiki](http://localhost:8000/wiki)"

### More options

golinks uses [conf](https://github.com/popodid/conf) for configurations
and currently supports two sources, environment variables and `golinks_config.yaml` file.

For more advanced configurations, checkout the [conf](https://github.com/popodid/conf),
the [config struct](https://github.com/haostudio/golinks/blob/master/cmd/golinks/main.go#L32),
and the [sample config](https://github.com/haostudio/golinks/blob/master/configs/local.yaml).

| Environment Variable / Yaml Path                                        | Type   | Default                             | Usage                                         |
| ----------------------------------------------------------------------- | ------ | ----------------------------------- | --------------------------------------------- |
| `PORT` / `Port`                                                         | int    | `8000`                              | Listening port                                |
| `HTTP_GOLINKS_WIKI` / `Http.Golinks.Wiki`                               | bool   | `false`                             | Serve wiki                                    |
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
