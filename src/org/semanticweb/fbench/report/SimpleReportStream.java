package org.semanticweb.fbench.report;

import java.util.List;

import org.semanticweb.fbench.query.Query;



/**
 * Simple reportstream implementation that prints to stdout.
 * 
 * @author as
 *
 */
public class SimpleReportStream implements ReportStream {

	
	
	public SimpleReportStream() {
		
	}
	
	
	@Override
	public void beginEvaluation(String dataConfig, List<String> querySet,
			int numberOfQueries, int numberOfRuns) {
		System.out.println("[EVAL] - Begin eval on dataConfig " + dataConfig + ", number of Queries " + numberOfQueries);
		
	}

	@Override
	public void beginQueryEvaluation(Query query, int run) {
		System.out.println("[EVAL] - executing query " + query.getIdentifier() + " (#" + run + ") ...");
	}

	@Override
	public void endEvaluation(long duration) {
		System.out.println("[EVAL] - End of evaluation. Duration " + duration + "ms");		
	}

	@Override
	public void endQueryEvaluation(Query query, int run, long duration,
			int numberOfResults) {
		System.out.println("[EVAL] - " + query.getIdentifier() + " (#" + run + ", duration: " + duration + "ms, results " + numberOfResults + ")");
	}

	@Override
	public void beginRun(int run, int totalNumberOfRuns) {
		System.out.println("[EVAL] - begin run " + run + " / " + totalNumberOfRuns );
	}


	@Override
	public void endRun(int run, int totalNumberOfRuns, long duration) {
		System.out.println("[EVAL] - end run " + run + " / " + totalNumberOfRuns + ". Duration " + duration + "ms");
	}
	
	@Override
	public void close() throws Exception {
		;		
	}
	
	@Override
	public void open() throws Exception {
		;		
	}

	@Override
	public void initializationBegin() {
		System.out.println("[INIT] Begin of initialization.");		
	}


	@Override
	public void initializationEnd(long duration) {
		System.out.println("[INIT] End of initialization. Duration: " + duration + "ms");
		
	}


	@Override
	public void reportDatasetLoadTime(String id, String name, String location, String type, long duration) {
		System.out.println("[INIT] Added dataset: id=" +id + ", name=" + name + ", location=" + location + ", type=" + type + ", duration=" + duration + "ms");
	}

}
