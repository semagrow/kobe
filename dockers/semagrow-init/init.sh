#!/bin/bash

#!/bin/bash

#
# environment variables
#
# Notice: kobe controller places dump files in /sevod-scraper/input
# and expects the output metadata file to be placed in /sevod-scraper/output
# DATASET_NAME and DATASET_ENDPOINT are initialized by kobe controller
#

INPUT="/sevod-scraper/input"
OUTPUT="/sevod-scraper/output"
DUMP="/tmp/dump"
TEMP="/tmp/temp"

SEVOD_SCRAPER="/sevod-scraper/assembly/target/bin/sevod-scraper.sh"


#
# if the metadata file already exists, simply change the endpoint
# and don't call sevod scraper
#

if [ -f $OUTPUT/$DATASET_NAME.ttl ]
then
  sed -i \
    's|void:sparqlEndpoint <[^ ]\+>|void:sparqlEndpoint <'$DATASET_ENDPOINT'>|g' \
    $OUTPUT/$DATASET_NAME.ttl
  exit
fi

#
# temporary directories
#

mkdir -p $DUMP
mkdir -p $TEMP

#
# convert to ntriples format all dump files
# we currently support .rdf, .owl, .n3, and .ttl files
#

for FILE in $INPUT/*.rdf
do
  rapper $FILE > $TEMP/`uuidgen`.nt
done

for FILE in $INPUT/*.owl
do
  rapper $FILE > $TEMP/`uuidgen`.nt
done

for FILE in $INPUT/*.n3
do
  serdi  $FILE > $TEMP/`uuidgen`.nt
done

for FILE in $INPUT/*.ttl
do
  serdi  $FILE > $TEMP/`uuidgen`.nt
done

cp *.nt $TEMP

#
# merge all converted ntriples files into a single ntriples file
#

cat $TEMP/*.nt > $DUMP/$DATASET_NAME.nt

#
# run sevod scraper tool for extracting metadata in rdfdump mode
#

$SEVOD_SCRAPER --rdfdump \
        -i $DUMP/$DATASET_NAME.nt \
        -e $DATASET_ENDPOINT \
        -o $OUTPUT/$DATASET_NAME.ttl

