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
      touch dump/dump.nt
      
      mkdir tmp
      cd tmp
      
      echo "starting data downloading"
      wget $DOWNLOAD_URL
      echo "finished downloading"
      
      tar xzvf *.tar.gz
      
      for FILE in `ls */*.*`
      do
        rapper -i guess -o ntriples $FILE >> ../dump/dump.nt
      done
      
      cd ..
      rm -rf tmp
      
      cp /config.json .
      sed -i 's/DATASET_NAME/'$DATASET_NAME'/g' config.json
      
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
