Requirements for queries:

We denote the minimal set of data sources for each query in the following:

 q_simple: DBpedia

 q1: DBpedia, NYTimes
 q2: DBPedia, Drugbank
 q3: Drugbank, KEGG
 q4:  DBPedia, LinkedMDB, NYTimes
 queries: DBPedia, LinkedMDB, NYTimes
   
Known Issues: 

* The drugbank SPARQL endpoint at high load times has very poor response times,
   all queries to be evaluated there do not work properly in such case. 
* The total evaluation time depends on the load of the server => Sometimes query
   response time is quite slow. If this is the case, try again later   
