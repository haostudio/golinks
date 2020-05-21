# Intro

[![MIT license](https://img.shields.io/badge/License-MIT-blue.svg)](https://github.com/haostudio/golinks/blob/master/LICENSE)
![go-build](https://github.com/haostudio/golinks/workflows/go-build/badge.svg)
![docker-master](https://github.com/haostudio/golinks/workflows/docker-master/badge.svg)
![docker-release](https://github.com/haostudio/golinks/workflows/docker-release/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/haostudio/golinks)](https://goreportcard.com/report/github.com/haostudio/golinks)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/haostudio/golinks)

Golinks is a open-sourced short link redirect service in Golang, under MIT-license.

![index](img/index.png)

## Feature Highlights

- Golinks styled (`go/mylink`) short link redirect
- Parameter substitution (`go/mylink/{VAR} -> https://mylink.com/{VAR}/mypage`)
- HTTP basic authentication
- Separate namespaces for multiple organizations
- Out-of-box solution with public docker image `haosutdio/golinks`
- Tracing with Jaeger
