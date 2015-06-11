@echo off
cd ..\..
SET baseDir=C:\Data\repositories

echo Starting local SPARQL endpoints

REM DBPedia351 => 10000
start startSparqlEndpoint.bat %baseDir%\native-storage.dbpedia351 10000
start startSparqlEndpoint.bat %baseDir%\native-storage.KEGG 10001
start startSparqlEndpoint.bat %baseDir%\native-storage.chEBI 10002
start startSparqlEndpoint.bat %baseDir%\native-storage.drugbank 10003
