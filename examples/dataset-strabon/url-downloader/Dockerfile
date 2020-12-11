FROM alpine:3.7

RUN apk update && apk upgrade && apk add wget

RUN mkdir -p /kobe/dataset

COPY init.sh /init.sh

ENTRYPOINT ["/bin/sh", "/init.sh"]
