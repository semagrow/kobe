#!/bin/sh

./scripts/dataConfig.sh $1 > dataConfig.ttl
shift

mainclass=org.semanticweb.fbench.FederationEval

# add all jars to classpath
for jar in `ls target/*.jar  target/dependency/*.jar`
do classpath=$classpath:$jar; done

java -Xmx1560m -cp $classpath $mainclass $*

rm dataConfig.ttl
