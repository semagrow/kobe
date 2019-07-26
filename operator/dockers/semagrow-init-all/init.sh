#!/bin/bash


cd /kobe/input
find -name '*.nt' | xargs cat > /kobe/output/metadata.nt 
java -jar /kobe-temp/ont-converter.jar -i /kobe/output/metadata.nt -if nt -o /kobe/output/metadata.ttl -of ttl
rm /kobe/output/metadata.nt


