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
		return "../queries/resources/" + queryFile;
	}

	
	/**
	 * load the queries from a queries file located at the path obtained by 
	 * {@link Config#getQuerySetPath()}.
	 * 
	 * @param querySetPath
	 * @return
	 * 			a list of queries for the query type
	 * @throws IOException
	 */

	public static List<Query> loadQueries(String querySetPath) throws IOException {
		ArrayList<Query> res = new ArrayList<Query>();

		File directory = new File(querySetPath);
		File[] listOfFiles = directory.listFiles();
		Arrays.sort(listOfFiles);
		
		for (File file : listOfFiles) {
			FileReader fin = new FileReader(file);
			BufferedReader in = new BufferedReader(fin);
			String tmp;
			String tmpQuery = "";

			while ((tmp = in.readLine()) != null){
				if (tmp.equals("")){
					if (!tmpQuery.equals(""))
						res.add(new Query(tmpQuery, file.getName(), -1));
					tmpQuery = "";
				}
				else {
					tmpQuery = tmpQuery + tmp + "\n";
				}
			}
			if (!tmpQuery.equals(""))
				res.add(new Query(tmpQuery, file.getName(), -1));
		}

		return res;
	}
}
