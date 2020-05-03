FROM golang:1.13.7-alpine

# Add workdir
RUN mkdir -p /opt/golinks
WORKDIR /opt/golinks

# Install dependeicies
RUN apk add make git curl

# Build binary
COPY . .
RUN make install-packr
RUN make tidy golinks

FROM alpine:latest

ARG GOLINKS_CONFIG

EXPOSE 80
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /opt/golinks/
COPY --from=0 /opt/golinks/build/golinks .

ADD configs/${GOLINKS_CONFIG}.yaml ./golinks_config.yaml

ENTRYPOINT ["./golinks"]
