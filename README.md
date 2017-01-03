# Ishi - Ad hoc HTTP Repeater

Ishi is ad hoc HTTP repeater/reverse proxy for development environment. It is simple and doesn't require file based configuration.

## Motivation

[Dinghy] has made it easier to develop web applications on macOS using Docker. It has a DNS server and nginx-proxy that resolves virtual hosts, so that we can develop multiple web services with one Docker machine.

When developing web services for mobile, I have struggled with testing the site running on dinghy from the iPhone. I have often wired the Mac and iPhone to the same LAN, set up a DNS server on the Mac, connected iPhone to the DNS server...

I felt this setting is very complicated. Like [`stone`] command, which is a simple TCP repeater, much more easier way to connect iPhone and web applications on Dinghy was needed. Unfortunately, `stone` can repeat traffics at TCP level, but it doesn't rewrite` Host` in the HTTP header. So I made `ishi`. Ishi means a stone in Japanese.

[`stone`]: http://manpages.org/stone
[Dinghy]: https://github.com/codekitchen/dinghy

## Usage

```
Usage:
  ishi [-l=<port>] <upstream>
  ishi -h | --help
  ishi --version

Arguments:
  upstream  Upstream host.

Options:
  -h --help             Show help.
  --version             Show version.
  -l --listen=<port>    Specify port to listen.
```

## Examples

Forwarding requests on 0.0.0.0:8000 to 192.168.99.100:80.

```
ishi 192.168.99.100
```

Forwarding requests on 0.0.0.0:80 to myapp.docker:80

```
ishi --listen 80 myapp.docker
```

Forwarding requests on 0.0.0.0:80 to myapp.docker:443

```
ishi --listen 80 https://myapp.docker
```

## Installation

To install Ishi, please use `go get`:

```
go get github.com/suin/ishi
```

## How Does It Works

1. Ishi starts to listen on Desktop (192.168.3.2:8000)
* Mobile device (192.168.3.2) on same LAN send HTTP request to `http://192.168.3.2:8000`
* Ishi fowards the request to Docker container. On the same time, Ishi overwrite `Host` header to `app.docker` from `192.168.3.2:8000` so that reverse proxy can foward the request to another container.
* The Docker container responds.
* Ishi forwards the response to the mobile device

![](https://raw.github.com/suin/ishi/master/image.png)
