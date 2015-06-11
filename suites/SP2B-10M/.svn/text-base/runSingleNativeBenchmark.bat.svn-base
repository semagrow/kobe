@echo off
cd ..\..
echo Running the Single Native Benchmark

mkdir suites\SP2B-10M\result

call runEval.bat suites\SP2B-10M\config\singleNative-config.prop
move result\loadTimes.csv suites\SP2B-10M\result\singleNative-load.csv
move result\result.csv suites\SP2B-10M\result\singleNative-result.csv
move result\result.nt suites\SP2B-10M\result\singleNative-result.nt
