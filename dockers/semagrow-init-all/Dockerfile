FROM xqdocker/ubuntu-openjdk:jdk-8

RUN mkdir /kobe-temp

ADD init.sh /script/init.sh

RUN ["chmod", "+x", "/script/init.sh"]

ENTRYPOINT ["/script/init.sh"]
