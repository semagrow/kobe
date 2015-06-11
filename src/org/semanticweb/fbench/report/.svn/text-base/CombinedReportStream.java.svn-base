package org.semanticweb.fbench.report;

import java.util.List;

import org.semanticweb.fbench.query.Query;



/**
 * Combined report stream for simple and cvs reporting. Uses the delegate pattern. 
 * 
 * @see SimpleReportStream
 * @see CsvReportStream2
 * 
 * @author as
 *
 */
public class CombinedReportStream implements ReportStream {

	
	protected SimpleReportStream simple = new SimpleReportStream();
	protected CsvReportStream2 csv = new CsvReportStream2();
	
	public void beginEvaluation(String dataConfig, List<String> querySet,
			int numberOfQueries, int numberOfRuns) {
		csv.beginEvaluation(dataConfig, querySet, numberOfQueries, numberOfRuns);
		simple.beginEvaluation(dataConfig, querySet, numberOfQueries, numberOfRuns);
	}
	public void beginQueryEvaluation(Query query, int run) {
		csv.beginQueryEvaluation(query, run);
		simple.beginQueryEvaluation(query, run);
	}
	public void beginRun(int run, int totalNumberOfRuns) {
		csv.beginRun(run, totalNumberOfRuns);
		simple.beginRun(run, totalNumberOfRuns);
	}
	public void close() throws Exception {
		csv.close();
		simple.close();
	}
	public void endEvaluation(long duration) {
		csv.endEvaluation(duration);
		simple.endEvaluation(duration);
	}
	public void endQueryEvaluation(Query query, int run, long duration,
			int numberOfResults) {
		csv.endQueryEvaluation(query, run, duration, numberOfResults);
		simple.endQueryEvaluation(query, run, duration, numberOfResults);
	}
	public void endRun(int run, int totalNumberOfRuns, long duration) {
		csv.endRun(run, totalNumberOfRuns, duration);
		simple.endRun(run, totalNumberOfRuns, duration);
	}
	public void initializationBegin() {
		csv.initializationBegin();
		simple.initializationBegin();
	}
	public void initializationEnd(long duration) {
		csv.initializationEnd(duration);
		simple.initializationEnd(duration);
	}
	public void open() throws Exception {
		csv.open();
		simple.open();
	}
	public void reportDatasetLoadTime(String id, String name, String location, String type, long duration) {
		csv.reportDatasetLoadTime(id, name, location, type, duration);
		simple.reportDatasetLoadTime(id, name, location, type, duration);
	}
}
