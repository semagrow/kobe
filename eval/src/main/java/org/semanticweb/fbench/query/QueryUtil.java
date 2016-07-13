package org.semanticweb.fbench.query;

import java.io.*;
import java.util.ArrayList;
import java.util.Arrays;
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
	 * @param queryFile
	 * @return
	 * 		the location of the query configuration for the specified type
	 */
	public static String getQueryLocation(String queryFile) {
		return "config/queries/" + queryFile;
	}

	/**
	 * queries are expected to be located in all files under directory config/queries/%querySet%/
	 * @param querySet
	 * @return
	 * 		the location of the query configuration for the specified type
	 */
	public static String getQueryDirectory(String querySet) {
		return "config/queries/" + querySet;
	}
	
	
	/**
	 * load the queries from a queries file located at the path obtained by 
	 * {@link #getQueryDirectory(QueryFile)}.
	 * 
	 * @param queryFile
	 * @return
	 * 			a list of queries for the query type
	 * @throws FileNotFoundException
	 * @throws IOException
	 */

	public static List<Query> loadQueries(String queryFile) throws FileNotFoundException, IOException {
		ArrayList<Query> res = new ArrayList<Query>();

		File directory = new File(getQueryDirectory(queryFile));
		File[] listOfFiles = directory.listFiles();
		Arrays.sort(listOfFiles);

		for (File file : listOfFiles) {
			FileReader fin = new FileReader(file);
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
		}

		return res;
	}
}
