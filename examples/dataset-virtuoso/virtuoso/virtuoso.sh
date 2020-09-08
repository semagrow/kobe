#!/bin/bash
set -e
cd /data
if [ "$ISTIO_USE" == "YES" ]; then
    until curl --head localhost:15021 ; do echo "Waiting for Sidecar" ; sleep 3 ; done ; echo "Sidecar available" ; sleep 10 ;
fi

mkdir -p dumps

if [ "$FORCE_LOAD" ] &&  [ -d "/kobe/dataset/$DATASET_NAME" ]; then
    echo "removing old files completely"   
    rm -r /kobe/dataset/$DATASET_NAME
fi

if  [ "$DATASET_NAME" ] && [ ! -f "/kobe/dataset/$DATASET_NAME/.data_loaded" ]  ;
then

    if [ "$DOWNLOAD_URL" ]; then
      echo "starting data downloading"
      mkdir -p toLoad
      cd toLoad
      wget $DOWNLOAD_URL
      tar xvf *.tar*
      cd ..
      echo "finished downloading"
    fi
    
    echo "starting data loading"
    pwd="dba"
    graph="http://localhost:8890/DAV"

    if [ "$DBA_PASSWORD" ]; then pwd="$DBA_PASSWORD" ; fi
    if [ "$DEFAULT_GRAPH" ]; then graph="$DEFAULT_GRAPH" ; fi
        
    virtuoso-t +wait +configfile /virtuoso.ini
    
    isql-v -U dba -P "$pwd" exec="ld_dir_all('toLoad', '*', '$graph');"
    
    cores=$(nproc --all)
    loaders=$(awk  'BEGIN { rounded = sprintf("%.0f", '$cores'/2.5); print rounded }')
    
   for ((n=1;n<=$loaders;n++)); do
      echo Starting RDF loader $n 
      isql-v -U dba -P "$pwd" exec="rdf_loader_run();" &
    done

    wait
    isql-v -U dba -P "$pwd" exec="checkpoint;"
    
    isql-v -U dba -P "$pwd" -K

    #save the dump and the database file in the nfs directory
    touch /data/.data_loaded

    if [ "$DATASET_NAME" ]; then
	cd /kobe/dataset/
    	mkdir -p $DATASET_NAME
    	cd $DATASET_NAME
    	mkdir -p database
    	mkdir -p dump
    	cp -r /data/toLoad/*/* dump/
    	cp /var/lib/virtuoso/db/virtuoso.db database/
    	touch /kobe/dataset/$DATASET_NAME/.data_loaded
    fi   
	 
    if [ "$DOWNLOAD_URL" ]; then
	rm -rf /data/toLoad/*
    fi
    
    echo "finished loading"
    
    while [ -f "/usr/local/virtuoso-opensource/var/lib/virtuoso/db/virtuoso.lck" ]
    do
      echo "waiting for virtuoso to close"
      sleep 10
    done

else 
    if [ "$DATASET_NAME" ]; then
        echo "copying old database.db"
    	cp  /kobe/dataset/$DATASET_NAME/database/virtuoso.db /var/lib/virtuoso/db/virtuoso.db
    fi 
fi

virtuoso-t +wait +foreground +configfile /virtuoso.ini

