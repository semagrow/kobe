package org.semanticweb.fbench.report;

import java.util.List;

import org.semanticweb.fbench.query.Query;
import org.semanticweb.fbench.query.QueryType;



/**
 * Combined report stream for rdf and csv reporting. Uses the delegate pattern. 
 * 
 * @see RdfReportStream
 * @see CsvReportStream2
 * 
 * @author as
 *
 */
public class CsvRdfReportStream implements ReportStream {

	
	protected RdfReportStream rdf = new RdfReportStream();
	protected CsvReportStream2 cvs = new CsvReportStream2();
	
	public void beginEvaluation(String dataConfig, List<String> querySet,
			int numberOfQueries, int numberOfRuns) {
		cvs.beginEvaluation(dataConfig, querySet, numberOfQueries, numberOfRuns);
		rdf.beginEvaluation(dataConfig, querySet, numberOfQueries, numberOfRuns);
	}
	public void beginQueryEvaluation(Query query, int run) {
		cvs.beginQueryEvaluation(query, run);
		rdf.beginQueryEvaluation(query, run);
	}
	public void beginRun(int run, int totalNumberOfRuns) {
		cvs.beginRun(run, totalNumberOfRuns);
		rdf.beginRun(run, totalNumberOfRuns);
	}
	public void close() throws Exception {
		cvs.close();
		rdf.close();
	}
	public void endEvaluation(long duration) {
		cvs.endEvaluation(duration);
		rdf.endEvaluation(duration);
	}
	public void endQueryEvaluation(Query query, int run, long duration,
			int numberOfResults) {
		cvs.endQueryEvaluation(query, run, duration, numberOfResults);
		rdf.endQueryEvaluation(query, run, duration, numberOfResults);
	}
	public void endRun(int run, int totalNumberOfRuns, long duration) {
		cvs.endRun(run, totalNumberOfRuns, duration);
		rdf.endRun(run, totalNumberOfRuns, duration);
	}
	public void initializationBegin() {
		cvs.initializationBegin();
		rdf.initializationBegin();
	}
	public void initializationEnd(long duration) {
		cvs.initializationEnd(duration);
		rdf.initializationEnd(duration);
	}
	public void open() throws Exception {
		cvs.open();
		rdf.open();
	}
	public void reportDatasetLoadTime(String id, String name, String location, String type, long duration) {
		cvs.reportDatasetLoadTime(id, name, location, type, duration);
		rdf.reportDatasetLoadTime(id, name, location, type, duration);
	}
}
