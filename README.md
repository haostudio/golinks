# GOLINKS

[![MIT license](https://img.shields.io/badge/License-MIT-blue.svg)](https://github.com/haostudio/golinks/blob/master/LICENSE)
![go-build](https://github.com/haostudio/golinks/workflows/go-build/badge.svg)
![docker-master](https://github.com/haostudio/golinks/workflows/docker-master/badge.svg)
![docker-release](https://github.com/haostudio/golinks/workflows/docker-release/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/haostudio/golinks)](https://goreportcard.com/report/github.com/haostudio/golinks)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/haostudio/golinks)

Golinks is a open-sourced short link redirect service in Golang, under MIT-license.

![index](./images/index.png?raw=true "Index")

## Feature Highlights

- Golinks styled (`go/mylink`) short link redirect
- Parameter substitution (`go/mylink/{VAR} -> https://mylink.com/{VAR}/mypage`)
- HTTP basic authentication
- Separate namespaces for multiple organizations
- Out-of-box solution with public docker image `haosutdio/golinks`
- Tracing with Jaeger

## Get started

### Run with docker

```sh
$ docker run -v \
  /path/to/datadir:/opt/golinks/datadir \
  -p 8000:8000 \
  haostudio/golinks
```

### Run with source code

#### Build binary

```sh
# Clone repo
$ git clone https://github.com/haostudio/golinks
# Install dependencies
$ make deps
# Build binary
$ make golinks
# Run
$ ./build/golinks
```

## Advanced options

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

For more options, checkout the [configuration guide](wiki/docs/install.md)

## Demo

### Edit link

![edit_link](./images/edit_link.png?raw=true "Edit Link")

### Register organization

![create_org](./images/create_org.png?raw=true "Create Org")

### Add user

![create_user](./images/create_user.png?raw=true "Create User")
