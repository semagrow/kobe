@echo off
cd ..\..
echo Running the Single Native Benchmark

mkdir suites\cross-domain\result

call runEval.bat suites\cross-domain\config\singleNative-config.prop
move result\loadTimes.csv suites\cross-domain\result\singleNative-load.csv
move result\result.csv suites\cross-domain\result\singleNative-result.csv
move result\result.nt suites\cross-domain\result\singleNative-result.nt
