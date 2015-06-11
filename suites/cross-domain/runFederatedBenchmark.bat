@echo off
cd ..\..
echo Running the Federated Benchmark

mkdir suites\cross-domain\result

call runEval.bat suites\cross-domain\config\federated-config.prop
move result\loadTimes.csv suites\cross-domain\result\federated-load.csv
move result\result.csv suites\cross-domain\result\federated-result.csv
move result\result.nt suites\cross-domain\result\federated-result.nt
