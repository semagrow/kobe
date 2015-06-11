package org.semanticweb.fbench.evaluation;

import org.apache.log4j.Logger;
import org.openrdf.query.BindingSet;
import org.openrdf.query.QueryEvaluationException;
import org.openrdf.query.QueryLanguage;
import org.openrdf.query.TupleQuery;
import org.openrdf.query.TupleQueryResult;
import org.openrdf.repository.Repository;
import org.openrdf.repository.RepositoryConnection;
import org.openrdf.repository.RepositoryException;
import org.semanticweb.fbench.Config;
import org.semanticweb.fbench.misc.TimedInterrupt;
import org.semanticweb.fbench.misc.Utils;
import org.semanticweb.fbench.query.Query;
import org.semanticweb.fbench.report.VoidReportStream;



/**
 * Implementation for query evaluation based on original Sesame Build
 * 
 * @author as
 */
public class SesameEvaluation extends Evaluation {

	public static Logger log = Logger.getLogger(SesameEvaluation.class);
	
	protected Repository sailRepo;
	protected RepositoryConnection conn;
	
	public SesameEvaluation() {
		super();
	}
	
	@Override
	public void finish() throws Exception {
		log.debug("Trying to close connection, interrupt time is 10000");
		new TimedInterrupt().run( new Runnable() {
			@Override
			public void run() {
				try {
					conn.close();
					log.info("Connection successfully closed.");
				} catch (RepositoryException e) {
					log.error("Error closing conenction: " + e.getMessage());
				}						
			}
		}, 10000);	
		sailRepo.shutDown();
	}

	@Override
	public void initialize() throws Exception {
		log.info("Performing Sesame Initialization...");
		
		sailRepo = SesameRepositoryLoader.loadRepositories(report);
		if (!Config.getConfig().isFill())
			conn = sailRepo.getConnection();
		
		log.info("Sesame Repository successfully initialized.");
	}

	@Override
	public int runQuery(Query q, int run) throws Exception {
		TupleQuery query = conn.prepareTupleQuery(QueryLanguage.SPARQL, q.getQuery());
		TupleQueryResult res = (TupleQueryResult) query.evaluate();
		int resCounter = 0;
		
		try {
			while(res.hasNext()){
				if (isInterrupted())
					throw new QueryEvaluationException("Thread has been interrupted.");
				BindingSet bindings = res.next();
				resCounter++;
				earlyResults.handleResult(bindings, resCounter);			
			}
		} finally {
			res.close();
		}
		return resCounter;
	}

	@Override
	public int runQueryDebug(Query q, int run, boolean showResult) throws Exception {
		TupleQuery query = conn.prepareTupleQuery(QueryLanguage.SPARQL, q.getQuery());
		TupleQueryResult res = (TupleQueryResult) query.evaluate();
		int resCounter = 0;
				
		while(res.hasNext()){
			if (isInterrupted())
				throw new QueryEvaluationException("Thread has been interrupted.");
			BindingSet bindings = res.next();
			if (showResult){
				// TODO use logging
				System.out.println(bindings);
			}
			
			resCounter++;
			earlyResults.handleResult(bindings, resCounter);	
		}
		
		return resCounter;
	}

		
	protected boolean isInterrupted() {
		return Thread.interrupted();
	}

	@Override
	public void reInitialize() throws Exception {
		log.info("Reinitializing repository and connection due to error in past results.");
		
		if (conn.isOpen())  {
			try {
				log.debug("Trying to close connection, interrupt time is 10000");
				Utils.closeConnectionTimeout(conn, 10000);
				
			} catch (Exception e) { ; /*ignore*/ }
		}
		
		if (sailRepo!=null) {
			boolean shutDown = Utils.shutdownRepositoryTimeout(sailRepo, 20*60*1000);	// timeout of 20 min for shutdown
			sailRepo = null;
			if (shutDown)
				log.debug("Repository shut down.");
			else
				log.debug("Failed to shut down repository within 20 minutes!");
		}
		
		
//		log.info("Waiting for 2 minutes to give server time for reinitialization");
//		Thread.sleep(120000);
		log.debug("loading repositories from scratch.");
		System.gc();
		sailRepo = SesameRepositoryLoader.loadRepositories(new VoidReportStream());
		conn = sailRepo.getConnection();
		log.debug("reinitialize done.");
	}
}
