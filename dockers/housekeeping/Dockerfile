FROM alpine

RUN apk update && apk upgrade && apk add bash

ADD init.sh /script/init.sh
RUN ["chmod", "+x", "/script/init.sh"]

ENTRYPOINT ["/script/init.sh"]
