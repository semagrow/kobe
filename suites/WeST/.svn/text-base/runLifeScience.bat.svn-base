@echo off
cd ..\..
echo Running the Life Science Benchmark

mkdir suites\SPLENDID\result

call runEval.bat suites\SPLENDID\lifeScience-config.prop
move result\loadTimes.csv suites\SPLENDID\result\SPLENDID-lifeScience-load.csv
move result\result.csv suites\SPLENDID\result\SPLENDID-lifeScience-result.csv
move result\result.nt suites\SPLENDID\result\SPLENDID-lifeScience-result.nt
