package org.semanticweb.fbench.query;

import java.io.BufferedReader;
import java.io.FileNotFoundException;
import java.io.FileReader;
import java.io.IOException;
import java.util.ArrayList;
import java.util.List;



/**
 * Utility class for queries.
 * 
 * @author as
 *
 */
public class QueryUtil {

	
	/**
	 * queries are expected to be located in a file at config\queries\%queryType%
	 * @param queryType
	 * @return
	 * 		the location of the query configuration for the specified type
	 */
	public static String getQueryLocation(String queryFile) {
		return "config/queries/" + queryFile;
	}
	
	
	/**
	 * load the queries from a queries file located at the path obtained by 
	 * {@link #getQueryLocation(QueryType)}. 
	 * 
	 * Expected format:
	 *  - Queries are SPARQL queries in String format
	 *  - queries are allowed to span several lines
	 *  - a query is intepreted to be finished if an empty line occurs
	 *  
	 *  Ex:
	 *  
	 *  QUERY1 ...
	 *   Q1 cntd
	 *   
	 *  QUERY2
	 * 
	 * @param queryType
	 * @return
	 * 			a list of queries for the query type
	 * @throws FileNotFoundException
	 * @throws IOException
	 */
	public static List<Query> loadQueries(String queryFile) throws FileNotFoundException, IOException {
		ArrayList<Query> res = new ArrayList<Query>();
		FileReader fin = new FileReader(getQueryLocation(queryFile));
		BufferedReader in = new BufferedReader(fin);
		String tmp;
		String tmpQuery = "";
		int nQuery=0;
		while ((tmp = in.readLine()) != null){
			if (tmp.equals("")){
				if (!tmpQuery.equals(""))
					res.add(new Query(tmpQuery, queryFile, ++nQuery));
				tmpQuery = "";
			}
			else {
				tmpQuery = tmpQuery + tmp + "\n";
			}
		}
		if (!tmpQuery.equals(""))
			res.add(new Query(tmpQuery, queryFile, ++nQuery));
		return res;
	}
}
