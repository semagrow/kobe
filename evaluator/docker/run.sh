#!/bin/bash

touch /kobe/run.prop

if [ "$TIMEOUT" ]; then
    echo timeout = "$TIMEOUT" >> /kobe/run.prop
fi

if [ "$EVAL_RUNS" ]; then
    echo evalRuns = "$EVAL_RUNS" >> /kobe/run.prop
fi

if [ "$EXPERIMENT" ]; then
    echo experimentName = "$EXPERIMENT" >> /kobe/run.prop
fi

cd /kobe/
sh runEval.sh "$ENDPOINT" /kobe/run.prop
cat /kobe/result/result.csv
