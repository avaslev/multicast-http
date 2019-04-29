Multicast http service
======================

[![Docker Pulls](https://img.shields.io/docker/pulls/avaslev/multicast-http.svg)](https://hub.docker.com/r/avaslev/multicast-http)
[![Docker Layers](https://images.microbadger.com/badges/image/avaslev/multicast-http.svg)](https://microbadger.com/images/avaslev/multicast-http)

Simple functional `Multicast http` service will send your request asynchronously to multiple addresses without guaranteeing delivery.

## Docker Images

Overview:

* Image are based on Alpine Linux
* [Docker Hub](https://hub.docker.com/r/avaslev/multicast-http)

## Environment Variables

| Variable                                   | Default Value              | Description                                     |
| ------------------------------------------ | -------------------------- | ----------------------------------              |
| `MULTICAST_HTTP_HEADER`                    |                            | Additional header                               |
| `MULTICAST_HTTP_HEADER_VALUE`              |                            | Additional header value                         |
| `MULTICAST_HTTP_HOSTS`                     |                            | List hosts comma separated: `https://ya.ru, goole.com, 127.0.0.1:5000`   |
| `MULTICAST_HTTP_K8S_POD_LABEL`             |                            | Mandatory for kubernetes. Pod label epm. `"app: symfony"`|
| `MULTICAST_HTTP_K8S_POD_PORT`              | `80`                       | Avalible pod port                               |


## Getting started on Kubernetes

Multicast http service can be added to any existing Kubernetes cluster. Here you can find RBAC and deployment file.

>If you want to add `Pod` to multicast, don't forget to add  `MULTICAST_HTTP_K8S_POD_LABEL` and `MULTICAST_HTTP_K8S_POD_PORT` if needed

```
kubectl apply -f https://raw.githubusercontent.com/avaslev/multicast-http/master/kubernetes/multicast-http.yaml
```