@Echo off
java -Dlog4j.configuration=file:config\log4j-proxy.properties -cp lib\fedbench.jar;lib\jetty7\*;lib\log4j\* org.semanticweb.fbench.proxy.StartJettyProxy %*