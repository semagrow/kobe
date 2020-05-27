#!/bin/bash

cd /kobe/input

FNAME="/kobe/output/metadata.ttl"

touch $FNAME

cat *.ttl | grep "^@prefix" | sort | uniq >> $FNAME
cat *.ttl | grep -v "^@prefix" >> $FNAME

cp /kobe-temp/repository.ttl /kobe/output/repository.ttl

