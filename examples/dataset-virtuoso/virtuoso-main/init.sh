#!/bin/bash

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

if [ "$USE_ISTIO" == "YES" ]
then
  until curl --head localhost:15021
  do echo "Waiting for Sidecar"
    sleep 3
  done
  echo "Sidecar available"
fi

/virtuoso-entrypoint.sh
