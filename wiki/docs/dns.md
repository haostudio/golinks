# Setup go path

It is extremely convenient to have the `go/mylink` setup. Here we provide a few
approaches to the `go/` path setup.

## Browser extension

The [Requestly](https://www.requestly.in/) extension allows you to redirect the
`go` domain to the domain/IP of your `golinks` server. Simply create a rule that
replaces `go` with your server domain.

![requestly](img/requestly.png)

## DNS Entry

You could setup a DNS entry if you are running in a corporate network.

## The `/etc/hosts` file

This will only work for you and your computer. Simply add an entry that points
to you server IP.

```
127.0.0.1 go
```
