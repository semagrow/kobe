# !/bin/sh

mainclass=org.semanticweb.fbench.FederationEval

# set current dir and bin/ as classpath
classpath=.:bin/

# add all jars to classpath
for jar in `ls lib/*.jar lib/*/*.jar`; do classpath=$classpath:$jar; done

# set logging options
logging="-Dlog4j.configuration=file:config/log4j.properties"

java -Xmx1560m $logging -cp $classpath $mainclass $*

# move result files to the suite's result directory
config_dir=${1%/*}
result_dir=$config_dir/result
for res in result\\*; do
  if [ -e $res ]; then
    dst=$result_dir/${res##*\\}
    mv $res $dst
  fi
done
