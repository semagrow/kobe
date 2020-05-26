package org.semanticweb.fbench.report;

import java.io.BufferedOutputStream;
import java.io.File;
import java.io.FileOutputStream;

import org.openrdf.query.BindingSet;
import org.semanticweb.fbench.query.Query;

/**
 * Print results to file "results/results.out"
 * 
 * @author Andreas Schwarte
 *
 */
public class FileEarlyResultsMonitor implements EarlyResultsMonitor {

	protected BufferedOutputStream bout = null;
	
	
	@Override
	public void handleResult(BindingSet bindings, int resultNumber) {
		try {
			bout.write( (bindings.toString() + "\r\n").getBytes("UTF-8"));
		} catch (Exception e) {
			throw new RuntimeException(e);
		}
	}	

	@Override
	public void nextQuery(Query q, long queryEvalStart) {
		; // no-op		
	}

	@Override
	public void queryDone() {
		; // no-op
		
	}

	@Override
	public void close() throws Exception {
		bout.flush();
		bout.close();		
	}

	@Override
	public void init() throws Exception {
		
		File file = new File("result/results.out");
		bout = new BufferedOutputStream( new FileOutputStream(file));
	}

}
