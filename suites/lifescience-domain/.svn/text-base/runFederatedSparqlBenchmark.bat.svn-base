@echo off
cd ..\..
echo Running the Federated SPARQL Benchmark

mkdir suites\lifescience-domain\result

call runEval.bat suites\lifescience-domain\config\federated-sparql-config.prop
move result\loadTimes.csv suites\lifescience-domain\result\federated-sparql-load.csv
move result\result.csv suites\lifescience-domain\result\federated-sparql-result.csv
move result\result.nt suites\lifescience-domain\result\federated-sparql-result.nt
