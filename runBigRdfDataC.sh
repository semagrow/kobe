
i=$1

rm config/queries/bigrdfdataC
ln -s C$i config/queries/bigrdfdataC

./runEval.sh doc/fedx/BigRdfC.prop | tee run-fedx-C$i.txt
sleep 10

#mv lib/fedxhibiscus/fedxhibiscus.jab lib/fedxhibiscus/fedxhibiscus.jar
#./runEval.sh doc/fedx-hibiscus/BigRdfC.prop | tee run-fedx-hibiscus-C$i.txt
#sleep 10
#mv lib/fedxhibiscus/fedxhibiscus.jar lib/fedxhibiscus/fedxhibiscus.jab

rm semagrow
ln -s semagrow-reactive semagrow
./runEval.sh suites/semagrow-reactive/bigrdfC/bigrdfC-config.prop | tee run-SG-C$i.txt
sleep 10

rm semagrow
ln -s semagrow-hibiscus semagrow
./runEval.sh suites/semagrow-hibiscus/bigrdfC/bigrdfC-config.prop | tee run-SG-hibiscus-C$i.txt
sleep 10
