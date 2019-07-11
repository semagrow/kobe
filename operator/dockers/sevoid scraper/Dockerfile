#FROM java:7-jdk
FROM xqdocker/ubuntu-openjdk:jdk-7

MAINTAINER Yiannis Mouchakis <gmouchakis@iit.demokritos.gr>

RUN apt-get update && apt-get install -y git && apt-get install -y maven --allow-unauthenticated && apt-get clean && rm -rf /var/lib/apt/lists/*

RUN git clone https://github.com/semagrow/sevod-scraper.git && \
    cd sevod-scraper && rm pom.xml
COPY pom.xml /sevod-scraper


WORKDIR /sevod-scraper 
RUN ls    
RUN mvn clean install
COPY log4j.properties /sevod-scraper/target/classes/

ENTRYPOINT ["/sevod-scraper/sevod-scraper.sh"]
