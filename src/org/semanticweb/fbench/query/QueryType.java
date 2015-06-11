package org.semanticweb.fbench.query;

/**
 * Definition of all query types and corresponding query files.
 * 
 * @author as
 *
 */
public enum QueryType {

	SIMPLE("Simple SPARQL query", "simple"),
	
	CROSSDOMAIN("Cross domain queries", "cross-domain"),
	
	LIFESCIENCEDOMAIN("Life science domain queries", "lifescience-domain"),
	
	SP2B("SP2B queries", "SP2B"),
	
	LINKEDDATA("Linked Data queries", "linked-data"),
	
	CUSTOM1("Custom SPARQL query", "custom1");
	
	
	private String desc;
	private String fileName;
	private QueryType(String desc, String fileName) {
		this.desc = desc;
		this.fileName = fileName;
	}
	
	public String getFileName() {
		return this.fileName;
	}
	
	public String toString() {
		return desc;
	}
	
	public void setFileName(String fileName) {
		this.fileName = fileName;
	}
}
