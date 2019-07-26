#FROM java:7-jdk
FROM xqdocker/ubuntu-openjdk:jdk-8

MAINTAINER Yiannis Mouchakis <gmouchakis@iit.demokritos.gr>

RUN apt-get update && apt-get install -y git && apt-get install -y maven --allow-unauthenticated && apt-get clean && rm -rf /var/lib/apt/lists/*

RUN git clone https://github.com/semagrow/sevod-scraper.git && \
    cd sevod-scraper && rm pom.xml
COPY pom.xml /sevod-scraper


WORKDIR /sevod-scraper 
RUN ls    
RUN mvn clean install

WORKDIR /sevod-scraper/assembly/target

RUN tar xzvf sevod-scraper-1.0-SNAPSHOT-dist.tar.gz

RUN mkdir /sevod-scraper/output
RUN mkdir /sevod-scraper/input

WORKDIR /sevod-scraper/input

COPY log4j.properties /sevod-scraper/target/classes/

COPY rdf2rdf-1.0.1-2.3.1.jar /sevod-scraper/
COPY ont-converter.jar /sevod-scraper/

COPY init.sh /sevod-scraper/

RUN ["chmod", "+x", "/sevod-scraper/init.sh"]

ENTRYPOINT ["/sevod-scraper/init.sh"]

#CMD sleep 3600
#CMD /sevod-scraper/assembly/target/bin/sevod-scraper.sh rdfdump /dump/cold-2010-complete.rdf.nt http://hi.com -spovj  out.n3 && tail -f /dev/null
#ENTRYPOINT ["/sevod-scraper/assembly/target/bin/sevod-scraper.sh"]
