FROM alpine:3.14

RUN apk add --no-cache raptor2

RUN mkdir -p /kobe/dataset

COPY init.sh /init.sh
COPY config.json /config.json

ENTRYPOINT ["/bin/sh", "/init.sh"]
