@echo off
cd ..\..
echo Running the Cross Benchmark

mkdir suites\SPLENDID\result

call runEval.bat suites\SPLENDID\crossDomain-config.prop
move result\loadTimes.csv suites\SPLENDID\result\SPLENDID-cross-load.csv
move result\result.csv suites\SPLENDID\result\SPLENDID-cross-result.csv
move result\result.nt suites\SPLENDID\result\SPLENDID-cross-result.nt
