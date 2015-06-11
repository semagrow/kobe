@echo off
cd ..\..
echo Running the Federated SPARQL Benchmark

mkdir suites\SP2B-10M\result

call runEval.bat suites\SP2B-10M\config\federated-sparql-config.prop
move result\loadTimes.csv suites\SP2B-10M\result\federated-sparql-load.csv
move result\result.csv suites\SP2B-10M\result\federated-sparql-result.csv
move result\result.nt suites\SP2B-10M\result\federated-sparql-result.nt
