# Golinks Server

## Usage

### Run with docker

```
$ docker run -v \
  /path/to/datadir:/opt/golinks/datadir \
  -p 8000:8000 \
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
