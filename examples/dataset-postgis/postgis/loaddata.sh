#!/bin/bash

#
# Waiting for services to start
#

while ! pg_isready 
do
  echo "$(date) - waiting for database to start"
  sleep 3
done

#until curl --head localhost:15000
#do
#  echo "Waiting for Sidecar"
#  sleep 3
#done

#sleep 10

#
# Downloading the dataset
#

mkdir /kobe
mkdir /kobe/dataset

echo "Downloading the dataset"

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
    tar xzvf *.tar.gz
    cd ..
    
    cp -r toLoad/*/* dump/
    rm -rf toLoad
    
    touch /kobe/dataset/$DATASET_NAME/.data_loaded
  fi
fi

#
# Loading the dataset
#

echo "Loading the dataset"

cd /kobe/dataset/$DATASET_NAME/dump

createdb -U postgres $DATASET_NAME

cat *.sql | psql -U postgres -d $DATASET_NAME

echo "done"
