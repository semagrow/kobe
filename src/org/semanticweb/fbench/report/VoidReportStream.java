package org.semanticweb.fbench.report;

import java.util.List;

import org.semanticweb.fbench.query.Query;

public class VoidReportStream implements ReportStream {

	@Override
	public void beginEvaluation(String dataConfig, List<String> querySet,
			int numberOfQueries, int numberOfRuns) {
		
	}

	@Override
	public void beginQueryEvaluation(Query query, int run) {
		
	}

	@Override
	public void beginRun(int run, int totalNumberOfRuns) {
		
	}

	@Override
	public void close() throws Exception {
		
	}

	@Override
	public void endEvaluation(long duration) {
		
	}

	@Override
	public void endQueryEvaluation(Query query, int run, long duration,
			int numberOfResults) {
	
	}

	@Override
	public void endRun(int run, int totalNumberOfRuns, long duration) {
		
	}

	@Override
	public void initializationBegin() {

	}

	@Override
	public void initializationEnd(long duration) {
		
	}

	@Override
	public void open() throws Exception {
		
	}

	@Override
	public void reportDatasetLoadTime(String id, String name, String location,
			String type, long duration) {
		
	}

}
