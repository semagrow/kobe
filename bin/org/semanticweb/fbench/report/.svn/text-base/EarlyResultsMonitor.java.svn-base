package org.semanticweb.fbench.report;

import org.openrdf.query.BindingSet;
import org.semanticweb.fbench.query.Query;


/**
 * Interface that allows measuring early results processing.
 * 
 * Note that when using this class you might introduce some overhead in the
 * overall query evaluation. This depends on your implementation.
 * 
 * Can be configured using config setting "earlyResultsMonitorClass".
 * 
 * @author as
 *
 */
public interface EarlyResultsMonitor {

	/**
	 * Perform any initialization here, e.g. open files
	 * 
	 * @throws Exception
	 */
	public void init() throws Exception;
	
	/**
	 * Perform any clean up operations here, e.g. close files.
	 * 
	 * @throws Exception
	 */
	public void close() throws Exception;
	
	
	/**
	 * Callback which can be used to measure early results performance.
	 * 
	 * Use the maintained queryEvalStart argument from nextQuery() to 
	 * compute the duration at which early result %resultNumber% 
	 * is ready
	 * 
	 * <code>
	 * long timeSinceStart = System.currentTimeMillis() - queryEvalStart;
	 * </code>
	 * 
	 * @param bindings
	 * @param resultNumber
	 */
	public void handleResult(BindingSet bindings, int resultNumber);
	
	
	/**
	 * Define the next query that is to be handled. This method is called
	 * each time a new query evaluation start. In particular note that the
	 * evaluation start time is specified as parameter. This time can be used
	 * in handleResult to compute the runtime at which early result x is ready.
	 * 
	 * @param q
	 * @param queryEvalStart
	 * 				the evaluation start time in milliseconds
	 */
	public void nextQuery(Query q, long queryEvalStart);
	
	
	/**
	 * Callback when the current query is processed completely, can be used
	 * to write back results (e.g. data points of type <earlyResult, afterDuration>)
	 * to some file.
	 */
	public void queryDone();
}
