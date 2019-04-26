Multicast http service
======================

Simple functional `Multicast http` service will send your request asynchronously to multiple addresses without guaranteeing delivery.


## Environment Variables

| Variable                                   | Default Value              | Description                                     |
| ------------------------------------------ | -------------------------- | ----------------------------------              |
| `MULTICAST_HTTP_HEADER`                    |                            | Additional header                               |
| `MULTICAST_HTTP_HEADER_VALUE`              |                            | Additional header value                         |
| `MULTICAST_HTTP_HOSTS`                     |                            | List hosts comma separated: `https://ya.ru,goole.com, 127.0.0.1:5000`   |
