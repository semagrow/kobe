package org.semanticweb.fbench.evaluation;

import java.io.FileOutputStream;
import java.io.OutputStream;
import java.util.List;

import org.semanticweb.fbench.LogUtils;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.openrdf.query.BindingSet;
import org.openrdf.query.QueryLanguage;
import org.openrdf.query.QueryResultHandlerException;
import org.openrdf.query.TupleQuery;
import org.openrdf.query.TupleQueryResultHandler;
import org.openrdf.query.TupleQueryResultHandlerException;
import org.openrdf.query.resultio.TupleQueryResultWriter;
import org.openrdf.query.resultio.TupleQueryResultWriterFactory;
import org.openrdf.query.resultio.text.csv.SPARQLResultsCSVWriterFactory;
import org.semanticweb.fbench.Config;
import org.semanticweb.fbench.query.Query;

import java.util.concurrent.CountDownLatch;

public class SesameEvaluationReactive extends SesameEvaluation {

	public static Logger log = LoggerFactory.getLogger(SesameEvaluationReactive.class);
	protected class MyHandler implements TupleQueryResultHandler {
		
		private int resCounter = 0;
		private TupleQueryResultHandler innerHandler = null;
                public CountDownLatch latch = new CountDownLatch(1);
		
		public MyHandler() { 
			latch = new CountDownLatch(1);	
		}
		
		public MyHandler(TupleQueryResultHandler handler) { 
			innerHandler = handler;
		}
		
		@Override
		public void endQueryResult()
				throws TupleQueryResultHandlerException {
			// TODO Auto-generated method stub
			if (innerHandler != null)
				innerHandler.endQueryResult();
			log.info(LogUtils.getCurrTime() + " [" + LogUtils.getQueryID() + "] Query evaluation End");
			latch.countDown();
            log.debug("countDown " + latch.getCount());
		}

		@Override
		public void handleBoolean(boolean arg0)
				throws QueryResultHandlerException {
			// TODO Auto-generated method stub
			if (innerHandler != null)
				innerHandler.handleBoolean(arg0);
		}

		@Override
		public void handleLinks(List<String> arg0)
				throws QueryResultHandlerException {
			// TODO Auto-generated method stub
			if (innerHandler != null)
				innerHandler.handleLinks(arg0);
		}

		@Override
		public void handleSolution(BindingSet arg0)
				throws TupleQueryResultHandlerException {
			if (resCounter == 0) {
				log.info(LogUtils.getCurrTime() + " [" + LogUtils.getQueryID() + "] Got first result");
			}
			// TODO Auto-generated method stub
			resCounter++;
			earlyResults.handleResult(arg0, resCounter);
			if (innerHandler != null)
				innerHandler.handleSolution(arg0);;
		}

		@Override
		public void startQueryResult(List<String> arg0)
				throws TupleQueryResultHandlerException {
			// TODO Auto-generated method stub
			resCounter = 0;
			if (innerHandler != null)
				innerHandler.startQueryResult(arg0);
		} 
			
		public int getCount() { return resCounter; }
	}
	
	@Override
	public int runQuery(Query q, int run) throws Exception {
		TupleQuery query = conn.prepareTupleQuery(QueryLanguage.SPARQL, q.getQuery());
		
		final int resCounter = 0; 
		
		MyHandler handler = new MyHandler();
		
		query.evaluate(handler);
		handler.latch.await();
		return handler.getCount();
	}

	@Override
	public int runQueryDebug(Query q, int run, boolean showResult) throws Exception {
		TupleQuery query = conn.prepareTupleQuery(QueryLanguage.SPARQL, q.getQuery());

		
		final int resCounter = 0; 
		
		MyHandler handler = new MyHandler();
		
		TupleQueryResultWriter writer = null;
		
		boolean writerStarted = false;
		
		if (showResult) {
			OutputStream results = new FileOutputStream(Config.getConfig().getBaseDir()+ "/result/" + q.getIdentifier() + "_" + run + ".csv");
			TupleQueryResultWriterFactory factory = new SPARQLResultsCSVWriterFactory();
			writer = factory.getWriter(results);
			handler = new MyHandler(writer);
		}
		
		query.evaluate(handler);
        handler.latch.await();
		return handler.getCount();
	}

	
}
