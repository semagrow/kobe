FROM alpine:3.7

RUN apk add --no-cache wget

RUN mkdir -p /kobe/dataset

COPY init.sh /init.sh

ENTRYPOINT ["/bin/sh", "/init.sh"]
