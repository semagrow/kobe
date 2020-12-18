#!/bin/sh

/usr/local/bin/rocket.sh &

until nc -z localhost 8080
do
  echo 'waiting for strabon to start'
  sleep 2
done

sleep 5

if [ "$DATASET_NAME" ]
then
  mkdir -p /kobe/dataset/$DATASET_NAME
  cd /kobe/dataset/$DATASET_NAME
    
  if [ -f "data_downloaded" ]
  then
    if [ ! -f "data_loaded" ]
    then
    
      service postgresql start
        
      for file in `ls dump`
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
        
        curl 'http://localhost:8080/Strabon/Store' \
          -H 'Accept: text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8' \
          -H 'Accept-Language: en-US,en;q=0.5' --compressed \
          -H 'Content-Type: application/x-www-form-urlencoded' \
          -H 'Origin: http://localhost:8080' \
          -H 'Connection: keep-alive' \
          -H 'Referer: http://localhost:8080/Strabon/Store' \
          -H 'Upgrade-Insecure-Requests: 1' \
          --data 'view=HTML&graph=&format='$format'&data=&url=file%3A%2F%2F%2Fkobe%2Fdataset%2F'$DATASET_NAME'%2Fdump%2F'$file'&fromurl=Store+from+URI'
        sleep 2
      done
      
      service postgresql stop
      
      mkdir postgis
      cp -r /var/lib/postgresql/9.4/main postgis
      
      touch data_loaded
      
    else
      echo "data already loaded"
    fi
  else
    echo "data not available"
  fi
else
  echo "dataset name not specified"
fi
