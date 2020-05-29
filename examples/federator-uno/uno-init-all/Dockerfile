FROM alpine:3.7

RUN mkdir /uno
RUN mkdir /uno/input
RUN mkdir /uno/output

COPY init.sh /uno/init.sh

ENTRYPOINT ["/bin/sh", "/uno/init.sh"]
