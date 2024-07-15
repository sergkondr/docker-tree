# docker-tree

[//]: # ([![GitHub release]&#40;https://img.shields.io/github/release/sergkondr/docker-tree.svg&#41;]&#40;https://github.com/sergkondr/docker-tree/releases/latest&#41;)
[//]: # ([![Go Report Card]&#40;https://goreportcard.com/badge/github.com/sergkondr/docker-tree&#41;]&#40;https://goreportcard.com/report/github.com/sergkondr/docker-tree&#41;)
[![License: MIT](https://img.shields.io/badge/License-MIT%202.0-blue.svg)](https://github.com/wagoodman/dive/blob/main/LICENSE)

This command shows the directory tree of a Docker image, like the 'tree' command. 
Provide the image name and an optional tag or digest to view the file structure inside the image. 
You can also specify a directory to see the file tree relative to this directory.

This is not a replacement for the amazing [Dive](https://github.com/wagoodman/dive) utility, but it works as a Docker plugin, so you might find it simpler and more convenient
Think of this app mainly as an attempt to understand how Docker images work and how to create Docker plugins. However, it does work, and I hope you find it useful.

### Install
```
cp ./docker-tree ~/.docker/cli-plugins/docker-tree
```

### Usage
```shell
➜ docker tree alpine:3.20 /etc/ssl
3.20: Pulling from library/alpine
a258b2a6b59a: Pull complete
Digest: sha256:b89d9c93e9ed3597455c90a0b88a8bbb5cb7188438f70953fede212a0c4394e0
Status: Downloaded newer image for alpine:3.20
docker.io/library/alpine:3.20
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
