package org.semanticweb.fbench.report;

import org.openrdf.query.BindingSet;
import org.semanticweb.fbench.query.Query;

/**
 * An early results monitor to demonstrate possibilities. Prints for each early result
 * the time after which it occurs. Note that this obviously influences performance and
 * adds some overhead to the overall query evaluation.
 * 
 * @author as
 *
 */
public class PrintEarlyResultsMonitor implements EarlyResultsMonitor {

	protected Query currentQuery;
	protected long currentQueryEvalStart;
	
	
	public PrintEarlyResultsMonitor() {
		
	}
	
	@Override
	public void handleResult(BindingSet bindings, int resultNumber) {

		long duration = System.currentTimeMillis() - currentQueryEvalStart;
		System.out.println("Result " + resultNumber + " ready after " + duration + "ms");
		
	}

	@Override
	public void nextQuery(Query q, long queryEvalStart) {
		this.currentQuery = q;
		this.currentQueryEvalStart = queryEvalStart;
	}

	@Override
	public void queryDone() {
		;	// no-op		
	}

	@Override
	public void close() throws Exception {
		;	// no-op			
	}

	@Override
	public void init() throws Exception {
		;	// no-op			
	}

}
