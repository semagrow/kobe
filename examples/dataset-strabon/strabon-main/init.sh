#!/bin/sh

service postgresql stop

if [ "$DATASET_NAME" ]
then
  mkdir -p /kobe/dataset/$DATASET_NAME
  cd /kobe/dataset/$DATASET_NAME
    
  if [ -f "data_loaded" ]
  then
    rm -rf /var/lib/postgresql/9.4/main
    cp -r postgis/main /var/lib/postgresql/9.4
    chown -R postgres:postgres /var/lib/postgresql/9.4/main
  else
    echo "data not already loaded"
  fi
else
  echo "dataset name not specified"
fi

service postgresql start

/usr/local/bin/rocket.sh