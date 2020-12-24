FROM openlink/virtuoso-opensource-7

RUN apt-get update && apt-get -y install netcat

RUN mkdir -p /kobe/dataset

COPY init.sh /init.sh
COPY loadfiles.sh /loadfiles.sh

ENTRYPOINT /init.sh 
