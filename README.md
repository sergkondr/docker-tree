# docker-tree

[![License: MIT](https://img.shields.io/badge/License-MIT%202.0-blue.svg)](https://github.com/sergkondr/docker-tree/blob/main/LICENSE)
[![GitHub release](https://img.shields.io/github/release/sergkondr/docker-tree.svg)](https://github.com/sergkondr/docker-tree/releases/latest)
[![Go Report Card](https://goreportcard.com/badge/github.com/sergkondr/docker-tree)](https://goreportcard.com/report/github.com/sergkondr/docker-tree)
[![Go](https://github.com/sergkondr/docker-tree/actions/workflows/go.yml/badge.svg)](https://github.com/sergkondr/docker-tree/actions/workflows/go.yml)
[![CodeQL](https://github.com/sergkondr/docker-tree/actions/workflows/codeql.yml/badge.svg)](https://github.com/sergkondr/docker-tree/actions/workflows/codeql.yml)

This command shows the directory tree of a Docker image, like the 'tree' command.
Provide the image name and an optional tag or digest to view the file structure inside the image.
You can also specify a directory to see the file tree relative to this directory.

This is not a replacement for the amazing [Dive](https://github.com/wagoodman/dive) utility, but it works as a Docker
plugin, so you might find it simpler and more convenient
Think of this app mainly as an attempt to understand how Docker images work and how to create Docker plugins. However,
it does work, and I hope you find it useful.

### Install

```
mv ./docker-tree ~/.docker/cli-plugins/docker-tree
```

### Usage

```shell
➜ docker tree alpine:3.20 /etc/ssl
processing image: alpine:3.20
ssl/
├── cert.pem
├── certs/
│   └── ca-certificates.crt
├── ct_log_list.cnf
├── ct_log_list.cnf.dist
├── openssl.cnf
├── openssl.cnf.dist
└── private/
```
