FROM alpine:3.7

RUN mkdir /input
RUN mkdir /output

COPY init.sh /init.sh
COPY repository.ttl /repository.ttl
COPY metadata.ttl /metadata.ttl

ENTRYPOINT ["/bin/sh", "/init.sh"]
