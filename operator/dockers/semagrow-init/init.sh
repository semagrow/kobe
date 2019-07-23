#!/bin/bash

FILE=/sevod-scraper/output/$DATASET_NAME.nt
if [ ! -f "$FILE" ] || [ "$INITIALIZE" ] ; then
    for file in *.rdf
    do
    # java -jar /sevod-scraper/rdf2rdf-1.0.1-2.3.1.jar ${file} ${file}.nt
    java -jar /sevod-scraper/ont-converter.jar -i ${file} -if rdf -o ${file}.nt -of nt
    #echo $?
    if [ $? != 0 ] ; then
        rm ${file}.nt && echo "removing it now "
    else 
        echo "not removing"
    fi
    done

    ls *.nt | xargs cat | sort -k 2 > $DATASET_NAME.nt

    /sevod-scraper/assembly/target/bin/sevod-scraper.sh rdfdump "$DATASET_NAME.nt" "$DATASET_ENDPOINT" -sp  "/sevod-scraper/output/$DATASET_NAME.ttl"

    # java -jar /sevod-scraper/rdf2rdf-1.0.1-2.3.1.jar "/sevod-scraper/output/$DATASET_NAME.ttl" "$DATASET_NAME.nt"
    java -jar /sevod-scraper/ont-converter.jar -i "/sevod-scraper/output/$DATASET_NAME.ttl" -if ttl -o "/sevod-scraper/output/$DATASET_NAME.nt" -of nt -v
fi
