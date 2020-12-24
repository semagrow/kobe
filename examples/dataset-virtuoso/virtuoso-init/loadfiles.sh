#!/bin/bash

pwd=`cat /settings/dba_password`
graph="http://localhost:8890/DAV"

if [ "$DBA_PASSWORD" ]; then pwd="$DBA_PASSWORD" ; fi
if [ "$DEFAULT_GRAPH" ]; then graph="$DEFAULT_GRAPH" ; fi
      
isql -U dba -P "$pwd" exec="ld_dir_all('toLoad', '*', '$graph');"
      
cores=$(nproc --all)
loaders=$(awk  'BEGIN { rounded = sprintf("%.0f", '$cores'/2.5); print rounded }')

for ((n=1;n<=$loaders;n++)); do
  echo Starting RDF loader $n 
  isql -U dba -P "$pwd" exec="rdf_loader_run();" &
done

wait
isql -U dba -P "$pwd" exec="checkpoint;"
