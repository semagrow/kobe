#!/bin/bash

echo Running the Life Science Benchmark

result_dir=suites/SPLENDID/result

cd ../..
mkdir $result_dir

./runEval.sh suites/SPLENDID/lifeScience-config.prop
mv result/loadTimes.csv $result_dir/SPLENDID-lifeScience-load.csv
mv result/result.csv $result_dir/SPLENDID-lifeScience-result.csv
mv result/result.nt $result_dir/SPLENDID-lifeScience-result.nt
