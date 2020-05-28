#!/bin/bash

until nc -z localhost 8080
do
  echo 'waiting for strabon to start'
  sleep 2
done

until curl --head localhost:15000
do
  echo "Waiting for Sidecar"
  sleep 3
done

sleep 10

mkdir /kobe
mkdir /kobe/dataset

if [ "$DATASET_NAME" ] && [ ! -f "/kobe/dataset/$DATASET_NAME/.data_loaded" ]
then
  
  if [ "$FORCE_LOAD" ] && [ -d "/kobe/dataset/$DATASET_NAME" ]
  then
    echo "removing old files completely"   
    rm -r /kobe/dataset/$DATASET_NAME
  fi
  
  if [ "$DOWNLOAD_URL" ]
  then
    
    cd /kobe/dataset/
    mkdir -p $DATASET_NAME
    cd $DATASET_NAME
    mkdir -p dump
    mkdir -p toLoad
    cd toLoad
  
    echo "starting data downloading"
    wget $DOWNLOAD_URL
    echo "finished downloading"
    tar xvf *.tar*
    cd ..
    
    cp -r toLoad/*/* dump/
    rm -rf toLoad
    
    touch /kobe/dataset/$DATASET_NAME/.data_loaded
  fi
fi

for file in `ls /kobe/dataset/$DATASET_NAME/dump`
do
  
  ext="${file##*.}"
  
  case "$ext" in
    "rdf")
      format='RDF/XML'
      ;;
    "nt")
      format='N-Triples'
      ;;
    "ttl")
      format='Turtle'
      ;;
    "n3")
      format='N3'
      ;;
    *)
      format='N-Triples';
      ;; 
  esac
  
  curl 'http://localhost:8080/strabon/Store' \
     -H 'Accept: text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8' \
     -H 'Accept-Language: en-US,en;q=0.5' --compressed \
     -H 'Content-Type: application/x-www-form-urlencoded' \
     -H 'Origin: http://localhost:8080' \
     -H 'Connection: keep-alive' \
     -H 'Referer: http://localhost:8080/strabon/Store' \
     -H 'Upgrade-Insecure-Requests: 1' \
     --data 'view=HTML&graph=&format='$format'&data=&url=file%3A%2F%2F%2Fkobe%2Fdataset%2F'$DATASET_NAME'%2Fdump%2F'$file'&fromurl=Store+from+URI'
  sleep 2
done
