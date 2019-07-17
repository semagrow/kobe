#FROM java:7-jdk
FROM xqdocker/ubuntu-openjdk:jdk-7

MAINTAINER Yiannis Mouchakis <gmouchakis@iit.demokritos.gr>

RUN mkdir /input
RUN mkdir /output
WORKDIR /input    
COPY rdf2rdf-1.0.1-2.3.1.jar /
ADD convert.sh /convert.sh
RUN chmod +x /convert.sh
#RUN java -jar rdf2rdf-1.0.1-2.3.1.jar -.rdf HALLO.nt

#ENTRYPOINT ["java","-jar","rdf2rdf-1.0.1-2.3.1.jar","input/-.rdf","/output/hallo.nt"]


#CMD ["/bin/bash", "/convert.sh"]

ENTRYPOINT ["/convert.sh"]
