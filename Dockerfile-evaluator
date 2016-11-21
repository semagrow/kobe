FROM maven:3-jdk-7

MAINTAINER Yiannis Mouchakis <gmouchakis@iit.demokritos.gr>

COPY ./eval/docker/run.sh .
COPY . /kobe
RUN chmod -v +x run.sh

RUN cd kobe && \
    mvn clean package && \
    mvn dependency:copy-dependencies

CMD /run.sh
