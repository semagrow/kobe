#FROM java:7-jdk
#FROM alpine
FROM xqdocker/ubuntu-openjdk:jdk-8

#RUN apk update && apk upgrade && apk add bash
RUN mkdir /kobe-temp
COPY ont-converter.jar /kobe-temp
ADD init.sh /script/init.sh
RUN ["chmod", "+x", "/script/init.sh"]

ENTRYPOINT ["/script/init.sh"]
