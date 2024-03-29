# metalbeat

[![Actions Status](https://github.com/craftslab/metalbeat/workflows/CI/badge.svg?branch=master&event=push)](https://github.com/craftslab/metalbeat/actions?query=workflow%3ACI)
[![Docker](https://img.shields.io/docker/pulls/craftslab/metalbeat)](https://hub.docker.com/r/craftslab/metalbeat)
[![Go Report Card](https://goreportcard.com/badge/github.com/craftslab/metalbeat)](https://goreportcard.com/report/github.com/craftslab/metalbeat)
[![License](https://img.shields.io/github/license/craftslab/metalbeat.svg?color=brightgreen)](https://github.com/craftslab/metalbeat/blob/master/LICENSE)
[![Tag](https://img.shields.io/github/tag/craftslab/metalbeat.svg?color=brightgreen)](https://github.com/craftslab/metalbeat/tags)



## Introduction

*metalbeat* is an agent of *[metalflow](https://github.com/craftslab/metalflow/)* written in Go.



## Prerequisites

- Go >= 1.15.0
- etcd == 3.3.25
- gRPC == 1.26.0



## Build

```bash
git clone https://github.com/craftslab/metalbeat.git

cd metalbeat
make build
```



## Run

```bash
./metalbeat --config-file="config.yml" --host-addr="127.0.0.1"
```



## Docker

```bash
git clone https://github.com/craftslab/metalbeat.git

cd metalbeat
docker build --no-cache -f Dockerfile -t craftslab/metalbeat:latest .
docker run craftslab/metalbeat:latest /metalbeat --config-file="/config.yml" --host-addr="127.0.0.1"
```



## Usage

```
usage: metalbeat --config-file=CONFIG-FILE --host-addr=HOST-ADDR [<flags>]

Metal Beat

Flags:
  --help                     Show context-sensitive help (also try --help-long
                             and --help-man).
  --version                  Show application version.
  --config-file=CONFIG-FILE  Config file (.yml)
  --host-addr=HOST-ADDR      Host address
```

```bash
# Fetch host address via ifconfig
ifconfig -a | grep inet | grep -v 127.0.0.1 | grep -v 172.17.0.1 | grep -v inet6 | awk '{print $2}' | tr -d "addr:"
```



## Settings

*metalbeat* parameters can be set in the directory [config](https://github.com/craftslab/metalbeat/blob/master/config).

An example of configuration in [config.yml](https://github.com/craftslab/metalbeat/blob/master/config/config.yml):

```yaml
apiVersion: v1
kind: agent
metadata:
  name: metalbeat
spec:
  etcd:
    host: 127.0.0.1
    port: 2379
```



## Design

![design](design.png)



## Etcd

- Agent

```
key: /metalflow/agent/{HOST}/register
val: metalbeat
```

- Master

```
key: /metalflow/worker/{HOST}/dispatch
val: {COMMAND}
```



## License

Project License can be found [here](LICENSE).



## Reference

- [etcd](https://etcd.io/docs/)
- [go-kit](https://github.com/go-kit/kit)
- [go-zero](https://github.com/tal-tech/go-zero)
- [gRPC](https://grpc.io/docs/languages/go/)
- [protobuf](https://developers.google.com/protocol-buffers/docs/proto3)
