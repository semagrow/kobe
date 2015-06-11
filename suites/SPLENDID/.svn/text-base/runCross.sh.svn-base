#!/bin/bash

echo Running the Cross Benchmark

result_dir=suites/SPLENDID/result

cd ../..
mkdir $result_dir

./runEval.sh suites/SPLENDID/crossDomain-config.prop
mv result/loadTimes.csv $result_dir/SPLENDID-cross-load.csv
mv result/result.csv $result_dir/SPLENDID-cross-result.csv
mv result/result.nt $result_dir/SPLENDID-cross-result.nt
