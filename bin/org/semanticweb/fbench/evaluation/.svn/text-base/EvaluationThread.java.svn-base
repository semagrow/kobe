package org.semanticweb.fbench.evaluation;

import org.apache.log4j.Logger;
import org.semanticweb.fbench.query.Query;
import org.semanticweb.fbench.report.EarlyResultsMonitor;
import org.semanticweb.fbench.report.ReportStream;



/**
 * Helper thread to allow for timeout handling while executing queries
 * 
 * @author as
 */
public class EvaluationThread extends Thread {
	
	public static Logger log = Logger.getLogger(EvaluationThread.class);
	
	protected Evaluation evaluator;
	protected Query query;
	protected ReportStream report;
	protected EarlyResultsMonitor earlyResults;
	protected int run;
	
	private boolean finished;
	private boolean error;
		
	public EvaluationThread(Evaluation evaluator, Query query, ReportStream report, EarlyResultsMonitor earlyResults, int run) {
		super();
		this.evaluator = evaluator;
		this.query = query;
		this.report = report;
		this.run = run;
		this.earlyResults = earlyResults;
		this.finished = false;
	}
	
	public boolean isFinished() {
		return this.finished;
	}
	
	public boolean isError() {
		return this.error;
	}
	
	@Override
	public void run() {
		try {
			error = false;
			log.info("Evaluation of query " + query.getIdentifier() + " in thread " + Thread.currentThread().getName());
			report.beginQueryEvaluation(query, run);
			long start = System.currentTimeMillis();
			earlyResults.nextQuery(query, start);
			int numberOfResults = evaluator.runQuery(query, run);
			long duration = System.currentTimeMillis() - start;
			report.endQueryEvaluation(query, run, duration, numberOfResults);
			log.info(query.getIdentifier() + " (#" + run + ", duration: " + duration + "ms, results " + numberOfResults + ")");
			
		} catch (IllegalMonitorStateException e) { 
			// reporting is done in evaluation (finished is still false)
			//log.info("Execution of query " + query.getIdentifier() + " resulted in timeout.");
			log.debug("Thread " + Thread.currentThread().getName() + " lost monitor. Timeout occurred. Thread will close.");
			error = true;
			return;
		} catch (Exception e) {
			report.endQueryEvaluation(query, run, -2, -1);
			log.error("Error executing query " + query.getIdentifier() + " (" + e.getClass().getSimpleName() + "): " + e.getMessage());
			log.debug("Exception details:", e);
			error = true;
		}
		this.finished = true;
		synchronized (Evaluation.class) {
			Evaluation.class.notify();	
		}
	}

}
