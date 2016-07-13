package org.semanticweb.fbench.report;

import java.io.BufferedReader;
import java.io.BufferedWriter;
import java.io.FileWriter;
import java.io.IOException;
import java.io.InputStreamReader;
import java.net.HttpURLConnection;
import java.net.URL;
import java.util.List;

import org.apache.log4j.Logger;
import org.semanticweb.fbench.evaluation.SesameSparqlEvaluation.RepoInformation;
import org.semanticweb.fbench.query.Query;


/**
 * Reports the query statistics for each endpoint, i.e. the number of HTTP requests sent 
 * to participating endpoints per query.
 * 
 * Results are printed to result\sparql_stats.csv
 * 
 * Format:
 * Query;endpointId1-requests;..;endpointIdN-requests;total;
 * 
 * Note: this report only works for local sparql servers, i.e. only SparqlServlet2 can handle
 * the requests.
 * 
 * @author as
 *
 */
public class SparqlQueryRequestReport {
	
	public static Logger log = Logger.getLogger(SparqlQueryRequestReport.class);

	protected List<RepoInformation> repoInformation;
	protected BufferedWriter bout;
	
	
	public SparqlQueryRequestReport() {
		
	}
	
	
	public void init(List<RepoInformation> repoInformation) throws IOException {
		this.repoInformation = repoInformation;
		bout = new BufferedWriter( new FileWriter("result/sparql_stats.csv"));
		
		bout.write("query;");
		for (RepoInformation r : repoInformation)
			bout.write(r.id + ";");
		bout.write("total;");
		bout.write("\r\n");
		
	}
	
	
	/**
	 * asks each endpoint for request counter stat. This is done by sending a get request 
	 * to http://endpoint/sparql?requestCount=true
	 * 
	 * @param query
	 */
	public void handleQuery(Query query) throws IOException {
		log.info("Requesting request count statistics from endpoints");
		bout.write(query.getIdentifier() + ";");
		
		int total = 0;
		for (RepoInformation r : repoInformation) {
			
			try {
				// make a get request to this url, and local server will return the count
				URL url = new URL(r.url + "?requestCount=true");
				HttpURLConnection conn = (HttpURLConnection) url.openConnection();
				BufferedReader bin = new BufferedReader( new InputStreamReader( conn.getInputStream() ));
				
				String _counter = bin.readLine();		// the first line is an integer representing the value
				bin.close();
				
				int requestCount = Integer.parseInt(_counter);
				total += requestCount;
				
				log.debug("## " + r.id + " => " + requestCount + " requests");
				bout.write(requestCount + ";");
			} catch (Exception e) {
				log.warn("Exception (" + e.getClass() + "):" + e.getMessage());
				bout.write("-1;");
			}
		}
		bout.write(total + ";");
		bout.write("\r\n");
		bout.flush();
	}

	public void finish() throws IOException {
		bout.flush();
		bout.close();
	}

}
