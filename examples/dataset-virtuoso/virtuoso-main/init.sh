#!/bin/sh

cd /database

if [ "$DATASET_NAME" ]
then    
  if [ -f "/kobe/dataset/$DATASET_NAME/data_loaded" ]
  then
    cp /kobe/dataset/$DATASET_NAME/database/* .
  else
    echo "data not already loaded"
  fi
else
  echo "dataset name not specified"
fi

/virtuoso-entrypoint.sh
