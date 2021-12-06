FROM alpine:3.7

RUN mkdir /kobe-temp

ADD init.sh /script/init.sh

RUN ["chmod", "+x", "/script/init.sh"]

ENTRYPOINT ["/script/init.sh"]
