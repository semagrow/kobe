@echo off
cd ..\..
echo Running the Federated Benchmark

mkdir suites\lifescience-domain\result

call runEval.bat suites\lifescience-domain\config\federated-config.prop
move result\loadTimes.csv suites\lifescience-domain\result\federated-load.csv
move result\result.csv suites\lifescience-domain\result\federated-result.csv
move result\result.nt suites\lifescience-domain\result\federated-result.nt
