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
- etcd >= 3.4.0
- gRPC >= 1.32.0



## Build

```bash
git clone https://github.com/craftslab/metalbeat.git

cd metalbeat
make build
```



## Run

```bash
./metalbeat --config-file="config.yml"
```



## Docker

```bash
git clone https://github.com/craftslab/metalbeat.git

cd metalbeat
docker build --no-cache -f Dockerfile -t craftslab/metalbeat:latest .
docker run -it craftslab/metalbeat:latest ./metalbeat --config-file="config.yml"
```



## Usage

```bash
usage: metalbeat --config-file=CONFIG-FILE [<flags>]

Metal Beat

Flags:
  --help                     Show context-sensitive help (also try --help-long
                             and --help-man).
  --version                  Show application version.
  --config-file=CONFIG-FILE  Config file (.yml)
```



## Settings

*metalbeat* parameters can be set in the directory [config](https://github.com/craftslab/metalbeat/blob/master/config).

An example of configuration in [config.yml](https://github.com/craftslab/metalbeat/blob/master/config/config.yml):

```yaml
metadata:
  name: metalbeat
spec:
  - name: metalmetrics
    role: worker
    node:
      container:
        image: craftslab/metalmetrics
        expose:
          - 8080
        command:
          - ./metalmetrics
          - --config-file="config.yml"
          - --output-file="output.json"
  - name: metalflow
    role: master
    node:
      bare:
        host: 127.0.0.1
        port: 2379
```



## Design

![design](design.png)



## License

Project License can be found [here](LICENSE).



## Reference

- [Protocol Buffers](https://developers.google.com/protocol-buffers/docs/proto3)
- [etcd](https://etcd.io/docs/)
- [gRPC in Go](https://grpc.io/docs/languages/go/)
