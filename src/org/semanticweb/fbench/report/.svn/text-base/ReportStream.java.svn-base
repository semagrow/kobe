package org.semanticweb.fbench.report;

import java.util.List;

import org.semanticweb.fbench.query.Query;



/**
 * Interface for any ReportStream implementation.
 * 
 * @author as
 *
 */
public interface ReportStream {

	/**
	 * called to open the report stream, e.g. to open a file stream.
	 * 
	 * @throws Exception
	 */
	public void open() throws Exception;
	
	/**
	 * called to close a report stream when operations are done, e.g. to close a file stream.
	 * @throws Exception
	 */
	public void close() throws Exception;
	
	
	/**
	 * called at the beginning of any evaluation.
	 * 
	 * @param dataConfig
	 * 			the location of the dataConfig that is used
	 * @param dataSet
	 * 			the query types that are executed
	 * @param numberOfQueries
	 * 			the number of queries that are executed
	 * @param numberOfRuns
	 * 			the number of evaluation runs
	 */
	public void beginEvaluation(String dataConfig, List<String> querySet, int numberOfQueries, int numberOfRuns);

	/**
	 * called at the beginning of a new run.
	 * 
	 * @param run
	 * 			the number of the run
	 * @param totalNumberOfRuns
	 * 			the total number of runs
	 */
	public void beginRun(int run, int totalNumberOfRuns);

	/**
	 * called before the query is being executed.
	 * 
	 * @param query
	 * 			the query that will be executed
	 * @param run
	 * 			the run
	 */
	public void beginQueryEvaluation(Query query, int run);
	
	
	/**
	 * called when the query evaluation is done.
	 * 
	 * @param query
	 * 			the query that has been executed
	 * @param run
	 * 			the run
	 * @param duration
	 * 			the duration in ms
	 * @param numberOfResults
	 * 			the number of results
	 */
	public void endQueryEvaluation(Query query, int run, long duration, int numberOfResults);
	
	/**
	 * called of the end of a run.
	 * 
	 * @param run
	 * 			the run
	 * @param totalNumberOfRuns
	 * 			the total number of runs
	 * @param duration
	 * 			the run duration in ms
	 */
	public void endRun(int run, int totalNumberOfRuns, long duration);
	
	
	/**
	 * called at the end of the overall evaluation.
	 * 
	 * @param duration
	 * 			the duration in ms
	 */
	public void endEvaluation(long duration);
	
	/**
	 * called at the beginning of the initialization, i.e. before any evaluation starts.
	 */
	public void initializationBegin();
	
	/**
	 * called to report a dataset that is being loaded, can be used to report load times during init.
	 * 
	 * @param id
	 * 			a unique identifier for this dataset
	 * @param name
	 * 			the data set name
	 * @param location
	 * 			the locations used in the dataset
	 * @param type
	 * 			the type of the dataset
	 * @param duration
	 * 			the load time in ms
	 */
	public void reportDatasetLoadTime(String id, String name, String location, String type, long duration);
	
	
	/**
	 * called when the initialization is done.
	 * 
	 * @param duration
	 * 			duration in ms.
	 */
	public void initializationEnd(long duration);
	

}
