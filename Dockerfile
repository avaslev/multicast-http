FROM golang:alpine
LABEL Maintainer="Alexey Vasilyev <alex@onlamp.ru>" \
      Description="Container with multicast http service"

ADD ./src /go/src/app
WORKDIR /go/src/app

EXPOSE 80

CMD ["go", "run", "main.go"]