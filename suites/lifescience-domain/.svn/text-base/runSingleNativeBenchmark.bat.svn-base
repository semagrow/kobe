@echo off
cd ..\..
echo Running the Single Native Benchmark

mkdir suites\lifescience-domain\result

call runEval.bat suites\lifescience-domain\config\singleNative-config.prop
move result\loadTimes.csv suites\lifescience-domain\result\singleNative-load.csv
move result\result.csv suites\lifescience-domain\result\singleNative-result.csv
move result\result.nt suites\lifescience-domain\result\singleNative-result.nt
