# Golinks Server

[![MIT license](https://img.shields.io/badge/License-MIT-blue.svg)](https://github.com/haostudio/golinks/blob/master/LICENSE)
![Go](https://github.com/haostudio/golinks/workflows/Go/badge.svg)
![Docker](https://github.com/haostudio/golinks/workflows/Docker/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/haostudio/golinks)](https://goreportcard.com/report/github.com/haostudio/golinks)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/haostudio/golinks)

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
authentication and stores all the links in a shared namespace.

```
$ docker run -v \
  /path/to/datadir:/opt/golinks/datadir \
  -p 8000:8000 \
  -e AUTHPROVIDER_NOAUTH_ENABLED=true \
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
