#!/bin/sh

FNAME="/sevod-scraper/output/"$DATASET_NAME".nt"

touch $FNAME

echo "<http://"$DATASET_NAME"> <http://fluidops.org/config#store> \"SPARQLEndpoint\" ." >> $FNAME
echo "<http://"$DATASET_NAME"> <http://fluidops.org/config#SPARQLEndpoint> \""$DATASET_ENDPOINT"\" ." >> $FNAME
echo "<http://"$DATASET_NAME"> <http://fluidops.org/config#supportsASKQueries> \"false\" ." >> $FNAME
