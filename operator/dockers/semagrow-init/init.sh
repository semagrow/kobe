#!/bin/bash

for file in *.rdf
do
   java -jar /sevod-scraper/rdf2rdf-1.0.1-2.3.1.jar ${file} ${file}.nt
done

for file in *.nt
do
   /sevod-scraper/assembly/target/bin/sevod-scraper.sh rdfdump ${file} http://testingsss.com -sp  /sevod-scraper/output/${file}.n3
done 
