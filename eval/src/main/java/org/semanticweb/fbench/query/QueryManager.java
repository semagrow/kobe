package org.semanticweb.fbench.query;

import java.io.FileNotFoundException;
import java.io.IOException;
import java.util.ArrayList;
import java.util.List;

import org.semanticweb.fbench.Config;



/**
 * Query Manager that is responsible for loading all queries.
 * 
 * @author as
 */
public class QueryManager {

	
	private static QueryManager instance = null;
	
	public static QueryManager getQueryManager() {
		if (instance==null)
			throw new RuntimeException("QueryManager not initialized. Call Config.load() first.");
		return instance;
	}	
	
	public static void initialize() throws FileNotFoundException, IOException {
		instance = new QueryManager();
		instance.init();
	}
	
	protected ArrayList<Query> queries;
	
	private QueryManager() {
		this.queries = new ArrayList<Query>();
	}
	
	/**
	 * initialize and load the queries as specified by {@link Config#getQuerySetPath()}.
	 *
	 * @throws IOException
	 */
	private void init() throws IOException {
		queries.addAll( QueryUtil.loadQueries(Config.getConfig().getQuerySetPath()) );
	}
	
	/**
	 * @return
	 * 		the initialized queries
	 */
	public List<Query> getQueries() {
		return queries;
	}
	
	
	/**
	 * @param queryType
	 * @return
	 * 		the initialized queries corresponding to queryType
	 */
	public List<Query> getQueries(String queryFile) {
		ArrayList<Query> res = new ArrayList<Query>();
		for (Query q : queries)
			if (q.getType().equals(queryFile))
				res.add(q);
		return res;
	}
}
