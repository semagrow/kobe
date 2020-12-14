#!/bin/sh

if [ "$DATASET_NAME" ]
then
  mkdir /kobe/dataset/$DATASET_NAME
  cd /kobe/dataset/$DATASET_NAME
    
  if [ ! -f "data_downloaded" ]
  then
    if [ "$DOWNLOAD_URL" ]
    then
    
      mkdir dump
      mkdir toLoad
      cd toLoad
    
      echo "starting data downloading"
      wget $DOWNLOAD_URL
      echo "finished downloading"
      
      tar xzvf *.tar.gz
      
      cd ..
      
      cp -r toLoad/*/* dump/
      rm -rf toLoad
      
      touch data_downloaded
      
    else
      echo "download url not specified"
    fi
  else
    echo "data already downloaded"
  fi
else
  echo "dataset name not specified"
fi