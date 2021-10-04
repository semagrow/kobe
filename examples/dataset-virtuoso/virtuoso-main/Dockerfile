FROM openlink/virtuoso-opensource-7

RUN apt-get update && apt-get install -y curl

RUN mkdir -p /kobe/dataset

COPY init.sh /init.sh

ENTRYPOINT /init.sh
