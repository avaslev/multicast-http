FROM alpine:latest
LABEL Maintainer="Alexey Vasilyev <alex@onlamp.ru>" \
        Description="Container with multicast http service"

ARG builder='true'

RUN apk update \
        && apk add --no-cache musl-dev go git

# Configure Go
ENV GOROOT /usr/lib/go
ENV GOPATH /go
ENV GOBIN /go/bin
ENV PATH $GOBIN:$PATH

ENV MULTICAST_HTTP_K8S_POD_PORT 80
ENV MULTICAST_HTTP_DEBUG 0

COPY . $GOPATH
WORKDIR $GOPATH

# Build and cleanup
RUN if [ ${builder} == 'true' ]; then \
        go install main.go \
        && ls $GOPATH | grep -v bin | xargs rm -rf \
        && rm -rf $GOPATH/.[^.]*\
        && apk update \
        && apk del go musl-dev git \
;fi

EXPOSE 80

CMD ["main"]