@Echo off
java -Xmx1024m -Dlog4j.configuration=file:config\log4j-sparql.properties -cp bin;lib\fedbench.jar;lib\jetty7\*;lib\log4j\*;lib\slf4j\*;lib\sesame\* org.semanticweb.fbench.sparqlendpoint.StartJettySparqlEndpoint %*
exit