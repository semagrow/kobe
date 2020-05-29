FROM xqdocker/ubuntu-openjdk:jdk-8

RUN mkdir /sevod-scraper
RUN mkdir /sevod-scraper/output
RUN mkdir /sevod-scraper/input

WORKDIR /sevod-scraper/input

COPY init.sh /sevod-scraper/

RUN ["chmod", "+x", "/sevod-scraper/init.sh"]

ENTRYPOINT ["/sevod-scraper/init.sh"]
