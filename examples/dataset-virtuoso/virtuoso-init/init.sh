#!/bin/bash

export VIRTUOSO_INI_FILE="/virtuoso.ini"

if [ "$DATASET_NAME" ]
then
  mkdir -p /kobe/dataset/$DATASET_NAME
    
  if [ -f "/kobe/dataset/$DATASET_NAME/data_downloaded" ]
  then
    if [ ! -f "/kobe/dataset/$DATASET_NAME/data_loaded" ]
    then
      /virtuoso-entrypoint.sh &
      
      until nc -z localhost 1111
      do
        echo 'waiting for virtuoso to start'
        sleep 2
      done

      sleep 5
      
      echo "starting data loading"

      mkdir /database/toLoad
      cp /kobe/dataset/$DATASET_NAME/dump/* /database/toLoad
      
      /loadfiles.sh
      
      rm -rf /database/toLoad
      cp -R /opt/virtuoso-opensource/database /kobe/dataset/$DATASET_NAME
      
      touch /kobe/dataset/$DATASET_NAME/data_loaded
      
      echo "data loaded"
    else
      echo "data already loaded"
    fi
  else
    echo "data not available"
  fi
else
  echo "dataset name not specified"
fi
