#!/bin/bash

mkdir -p /valid
for i in `seq 0 $N`; 
do 
  DATASET_NAME="DATASET_NAME_$i"
  DATASET=${!DATASET_NAME}
  cp "/kobe/$DATASET/$FEDERATOR_NAME/$DATASET.nt" /valid 
  echo "hi"
done

#cd /kobe

mkdir -p /kobe/"$FEDERATION_NAME"

mkdir -p /meta

cd /valid
find -name '*.nt' | grep -v dump | xargs cat > /meta/metadata.nt 

cp /meta/metadata.nt /kobe/"$FEDERATION_NAME"/



