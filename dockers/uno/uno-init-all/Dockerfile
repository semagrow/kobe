FROM alpine:3.7

WORKDIR /uno

COPY init.sh /uno/init.sh

RUN ["chmod", "+x", "/uno/init.sh"]

ENTRYPOINT ["/uno/init.sh"]
