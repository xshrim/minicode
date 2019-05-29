# Docker环境部署

## 目录

<!-- TOC -->

- [Docker环境部署](#docker环境部署)
    - [目录](#目录)
    - [前置准备](#前置准备)
    - [Docker环境](#docker环境)
    - [环境部署](#环境部署)

<!-- /TOC -->

## 前置准备

Docker环境为独立的部署环境, 区块链集群, MongoDB集群, 事件驱动和监控服务均已部署到Docker镜像中(简单起见并未使用多容器协同), 只需启动容器并进行相应端口映射和卷挂载即可.

## Docker环境

1. Docker安装 (略)

2. 镜像下载

```bash
docker pull xshrim/cmdev:v1.0
```

## 环境部署

启动容器:

```bash
docker run -itd -p 2222:22 -p 8515:8515 -p 8516:8516 -p 9927:9927 -p 9928:9928 -p 9929:9929 -p 8000:8000 --name backend xshrim/cmdev:v1.0
```

[注]: 建议先将容器内/data目录复制到主机, 然后进行适配修改再重新启动容器, 挂载本地目录到容器中:

```bash
docker run -itd -p 2222:22 -p 8515:8515 -p 8516:8516 -p 9927:9927 -p 9928:9928 -p 9929:9929 -p 8000:8000 -v ~/data:/data --name backend xshrim/cmdev:v1.0
```
