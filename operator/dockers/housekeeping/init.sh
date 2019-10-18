#!/bin/bash

for i in `seq 0 $N`; 
do 
  DATASET_NAME="DATASET_NAME_$i"
  DATASET=${!DATASET_NAME}
  cp /kobe/$DATASET/$FEDERATOR_NAME/* "/kobe/temp-$FEDERATION_NAME"
  echo "copying file to temp folder "
done


