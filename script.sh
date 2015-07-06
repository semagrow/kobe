rm semagrow
ln -s semagrow-hibiscus semagrow
./runEval.sh suites/semagrow-hibiscus/bigrdfC/bigrdfC-config.prop | tee bigrdfC-hibiscus.txt
rm semagrow
ln -s semagrow-tbss semagrow
./runEval.sh suites/semagrow-tbss/bigrdfC/bigrdfC-config.prop | tee bigrdfC-tbss.txt
