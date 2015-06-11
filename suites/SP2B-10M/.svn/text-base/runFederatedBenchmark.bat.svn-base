@echo off
cd ..\..
echo Running the Federated Benchmark

mkdir suites\SP2B-10M\result

call runEval.bat suites\SP2B-10M\config\federated-config.prop
move result\loadTimes.csv suites\SP2B-10M\result\federated-load.csv
move result\result.csv suites\SP2B-10M\result\federated-result.csv
move result\result.nt suites\SP2B-10M\result\federated-result.nt
