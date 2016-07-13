package org.semanticweb.fbench.query;


/**
 * Data class for a SPARQL query
 * 
 * @author as
 *
 */
public class Query {

	protected String query;				// the query itself
	protected String queryFile;			// the type of the query (defines the file)
	protected int number;				// the number of this query in queries file
	
	public Query(String query, String queryFile, int number) {
		super();
		this.query = query;
		this.queryFile = queryFile;
		this.number = number;
	}

	public String getQuery() {
		return query;
	}

	public String getType() {
		return queryFile;
	}

	public int getNumber() {
		return number;
	}
	
	public String getIdentifier() {
		return queryFile + "_" + number;
	}
	
	@Override
	public String toString() {
		return query;
	}
	

	@Override
	public int hashCode() {
		final int prime = 31;
		int result = 1;
		result = prime * result + number;
		result = prime * result + ((queryFile == null) ? 0 : queryFile.hashCode());
		return result;
	}

	@Override
	public boolean equals(Object obj) {
		if (this == obj)
			return true;
		if (obj == null)
			return false;
		if (getClass() != obj.getClass())
			return false;
		Query other = (Query) obj;
		if (number != other.number)
			return false;
		if (queryFile == null) {
			if (other.queryFile != null)
				return false;
		} else if (!queryFile.equals(other.queryFile))
			return false;
		return true;
	}
}
