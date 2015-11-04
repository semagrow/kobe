
i=$1

rm config/queries/bigrdfdataB
ln -s B$i config/queries/bigrdfdataB

./runEval.sh doc/fedx/BigRdfB.prop | tee run-fedx-B$i.txt
sleep 10

#mv lib/fedxhibiscus/fedxhibiscus.jab lib/fedxhibiscus/fedxhibiscus.jar
#./runEval.sh doc/fedx-hibiscus/BigRdfC.prop | tee run-fedx-hibiscus-C$i.txt
#sleep 10
#mv lib/fedxhibiscus/fedxhibiscus.jar lib/fedxhibiscus/fedxhibiscus.jab

./runEval.sh suites/semagrow-reactive/bigrdfB/bigrdfB-config.prop | tee run-SG-B$i.txt
sleep 10

#./runEval.sh suites/semagrow-hibiscus/bigrdfB/bigrdfB-config.prop | tee run-SG-hibiscus-B$i.txt
#sleep 10
