FROM postgis/postgis

RUN apt-get update \
 && apt-get install --assume-yes curl wget \
 && apt-get --assume-yes --purge autoremove \
 && apt-get clean \
 && rm -rf /var/lib/apt/lists/*

COPY loaddata.sh /docker-entrypoint-initdb.d/

