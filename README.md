Multicast http service
======================

[![Docker Pulls](https://img.shields.io/docker/pulls/avaslev/multicast-http.svg)](https://hub.docker.com/r/avaslev/multicast-http)
[![Docker Stars](https://img.shields.io/docker/stars/avaslev/multicast-http.svg)](https://hub.docker.com/r/avaslev/multicast-http)
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
| `MULTICAST_HTTP_HOSTS`                     |                            | List hosts comma separated: `https://ya.ru,goole.com, 127.0.0.1:5000`   |
