# Golinks Server

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
$ git clond https://github.com/haostudio/golinks
$ make deps
$ make golinks
```

#### Run

```
$ ./build/golinks
```
