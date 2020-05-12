FROM maven:3-jdk-8

MAINTAINER Antonis Troumpoukis <antru@iit.demokritos.gr>

RUN git clone https://github.com/semagrow/sevod-scraper.git \
 && cd sevod-scraper \
 && git checkout -b devel \
 && mvn clean package \
 && cd assembly/target \
 && tar xzvf sevod-scraper-1.0-SNAPSHOT-dist.tar.gz

RUN apt-get update \
 && apt-get install --assume-yes serdi raptor2-utils \
 && apt-get --assume-yes --purge autoremove \
 && apt-get clean \
 && rm -rf /var/lib/apt/lists/*

WORKDIR /sevod-scraper/

RUN mkdir /sevod-scraper/input
RUN mkdir /sevod-scraper/output

ENTRYPOINT ["/sevod-scraper/init.sh"]
