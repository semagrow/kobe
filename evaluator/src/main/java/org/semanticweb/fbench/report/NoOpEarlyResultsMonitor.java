package org.semanticweb.fbench.report;

import org.openrdf.query.BindingSet;
import org.semanticweb.fbench.query.Query;

/**
 * A simple early results monitor that performs no operations.
 * 
 * @author as
 */
public class NoOpEarlyResultsMonitor implements EarlyResultsMonitor {

	@Override
	public void handleResult(BindingSet bindings, int resultNumber) {
		; 	// no-op		
	}

	@Override
	public void nextQuery(Query q, long queryEvalStart) {
		; 	// no-op		
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
