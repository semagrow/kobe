# !/bin/sh

# collect all jars
for jar in `ls lib/*.jar lib/*/*.jar`; do path=$path:$jar; done

java -Dlog4j.configuration=file:config/log4j.properties -cp .$path org.semanticweb.fbench.proxy.StartJettyProxy $*
