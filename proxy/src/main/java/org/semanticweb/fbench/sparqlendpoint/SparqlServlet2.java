package org.semanticweb.fbench.sparqlendpoint;

import info.aduna.lang.FileFormat;
import info.aduna.lang.service.FileFormatServiceRegistry;

import java.io.BufferedOutputStream;
import java.io.BufferedReader;
import java.io.File;
import java.io.FileOutputStream;
import java.io.IOException;
import java.io.InputStreamReader;
import java.io.OutputStream;
import java.net.URLDecoder;
import java.util.ArrayList;
import java.util.HashSet;
import java.util.LinkedList;
import java.util.List;
import java.util.concurrent.Executor;
import java.util.concurrent.Executors;

import javax.servlet.ServletInputStream;
import javax.servlet.ServletOutputStream;
import javax.servlet.http.HttpServlet;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;

import org.apache.log4j.Logger;
import org.openrdf.http.server.ProtocolUtil;
import org.openrdf.query.BooleanQuery;
import org.openrdf.query.GraphQuery;
import org.openrdf.query.GraphQueryResult;
import org.openrdf.query.Query;
import org.openrdf.query.QueryLanguage;
import org.openrdf.query.QueryResultUtil;
import org.openrdf.query.TupleQuery;
import org.openrdf.query.TupleQueryResult;
import org.openrdf.query.resultio.BooleanQueryResultWriter;
import org.openrdf.query.resultio.BooleanQueryResultWriterFactory;
import org.openrdf.query.resultio.BooleanQueryResultWriterRegistry;
import org.openrdf.query.resultio.TupleQueryResultWriter;
import org.openrdf.query.resultio.TupleQueryResultWriterFactory;
import org.openrdf.query.resultio.TupleQueryResultWriterRegistry;
import org.openrdf.repository.Repository;
import org.openrdf.repository.RepositoryConnection;
import org.openrdf.repository.RepositoryException;
import org.openrdf.repository.sail.SailRepository;
import org.openrdf.rio.RDFWriter;
import org.openrdf.rio.RDFWriterFactory;
import org.openrdf.rio.RDFWriterRegistry;
import org.openrdf.sail.NotifyingSailConnection;
import org.openrdf.sail.SailException;
import org.openrdf.sail.nativerdf.NativeStore;
import org.openrdf.sail.nativerdf.NativeStoreConnection;




public class SparqlServlet2 extends HttpServlet {

	private static final long serialVersionUID = 2627590629243739807L;
	
	protected static Repository repo = null;
	protected static Executor executor = Executors.newCachedThreadPool();
	protected static LinkedList<QueryRequest> queryRequestQueue = new LinkedList<QueryRequest>();
	protected static int nextRequestId = 1;
	protected static HashSet<QueryRequest> activeQueries = new HashSet<QueryRequest>();
	
	protected List<WorkerThread> workers = new ArrayList<WorkerThread>();
	protected int nWorkers;
	protected int idleWorkers;
	
	public SparqlServlet2() {
		initializeRepository();
		initializeWorkerThreads();
	}
	
	private static final Logger log = Logger.getLogger(SparqlServlet2.class);
	
	
	@Override
	protected void doPost(HttpServletRequest req, HttpServletResponse resp){
		
		try {
//			ServletInputStream input = req.getInputStream();
//			InputStreamReader in = new InputStreamReader(input);
//			BufferedReader reader = new BufferedReader(in);
//			String tmp;
//			while ((tmp = reader.readLine()) != null){
//				query = query + tmp;
//			}
			
			String query;
			if (req.getParameter("q") != null){
				query = req.getParameter("q");
			}
			else if (req.getParameter("query") != null){
				query = req.getParameter("query");				
			}
			else {
				error(resp, 400, "You provided no query.");
				return;
			}
			
//			query = query.substring(6);
			query = URLDecoder.decode(query, "ISO-8859-1");

            handleQuery(query, req, resp, resp.getOutputStream());
		} 
		catch (Exception e) {
			log.error("Error: ", e);
			error(resp, 500, "Unexpected error: " + e.getMessage());
		}
	}
	
	@Override
	protected void doGet(HttpServletRequest req, HttpServletResponse resp){

		try {
			ServletOutputStream outputStream = resp.getOutputStream();	
			
			// XXX document!!! requestCount get parameter is served!
			if (req.getParameter("requestCount") != null) {
				handleCountRequest(req, resp, outputStream);
			}
			else if (req.getParameter("q") != null){
				String query = URLDecoder.decode(req.getParameter("q"), "ISO-8859-1");
				handleQuery(query, req, resp, outputStream);
			}
			else if (req.getParameter("query") != null){
				String query = URLDecoder.decode(req.getParameter("query"), "ISO-8859-1");
				handleQuery(query, req, resp, outputStream);
			}
			else {
				error(resp, 400, "You provided no query.");
			}
		} 
		catch (IOException e) {
			log.error("Error: ", e);
		}

	}
	
	
	@Override
	public void destroy() {
		try {
			log.info("Shutting down repository and closing connection.");
			repo.shutDown();
			log.info("Repository successfully shut down.");
		} catch (RepositoryException e) {
			log.error("Error while shutting down repository.", e);
		}
		
		// TODO stop worker threads
		
		super.destroy();
	}
	
	private void handleQuery(String query, HttpServletRequest req, HttpServletResponse resp, ServletOutputStream outputStream) {
		
		// stupid hook for DARQ (TODO remove)
		if (query.startsWith("BASE")) {
			log.debug("BASE mentioned in query, use HOOK:");
			if (query.contains("<>"))
				query = query.replace("<>", "<http://example.org>");
		}
		
		log.trace("Server retrieved query: " + query);
		QueryRequest qr = null;
		
		// additional delay to be configured via argument at startup
		if (StartJettySparqlEndpoint.requestDelay>0) {
    		try {
				Thread.sleep(StartJettySparqlEndpoint.requestDelay);
			} catch (InterruptedException e1) {
				;
			}
		}
		
		try {
			int currentRequest;
			
			synchronized (queryRequestQueue) {
				currentRequest = nextRequestId++;
				qr = new QueryRequest(currentRequest, query, req, resp, outputStream);
				queryRequestQueue.addLast(qr);
				queryRequestQueue.notify();
			}
					
			synchronized( activeQueries) {
				activeQueries.add(qr);
			}
			
			synchronized (resp) {
				try {
					// XXX
					// the while is just cheating to avoid deadlock
					// there is still some deadlock cause hidden anywhere!!!
					// somehow this response thread does not get notified
					// the access to queryrequest should be synched through resp
					while (!qr.isDone())
						resp.wait(500);				
				} catch (InterruptedException e) {
					log.warn("Request " + currentRequest + " was interrupted.");
					error(resp, 500, "Request was interrupted.");
					// TODO check if threads working on this are running, if yes, kill them
					// XXX are parallel child threads killed as well?
				} 
			}
			
		} catch (Exception e) {
			log.error("Unexpected exception occured.", e);
			error(resp, 500, "Unexpected error: " + e.getMessage());
		} finally {
			
			synchronized( activeQueries) {
				activeQueries.remove(qr);
			}
		}
	}
	
	protected void error(HttpServletResponse resp, int errorCode, String message) {
		
		try {
			log.error("Error (" + errorCode + "): " + message);
			resp.sendError(errorCode, message);
		} catch (IllegalStateException e) {
			log.warn("Error message could not be send. Stream is committed: " + message);
		} catch (IOException e) {
			log.error("Error message could not be sent", e);
		}
	}
	
	private void handleCountRequest(HttpServletRequest req, HttpServletResponse resp, ServletOutputStream outputStream) {
		
		log.info("Incoming request for count of past SPARQL queries");
		
		int count;
		synchronized (queryRequestQueue) {
			count = nextRequestId-1;
			nextRequestId = 1;
		}
		
		log.info("Reporting #requests= " + count + ", resetting counter to 0");
		try {
			resp.setStatus(200);
			outputStream.print(count);
			outputStream.flush();
			outputStream.close();
		} catch (Exception e) {
			log.error("Error while sending count report.", e);
			error(resp, 500, "Error while sending count report: " + e.getMessage() );
		}
		
		
	}
	
	
	protected QueryRequest getNextQueryRequestSynch() {
		
		synchronized (queryRequestQueue) {
			
			if (queryRequestQueue.isEmpty())
				return null;
			
			return queryRequestQueue.removeFirst();
		}
	}
	
	
	protected void initializeWorkerThreads() {
		
		int nWorkers = StartJettySparqlEndpoint.nWorkerThreads;
		
		log.info("Initializing " + nWorkers + " worker threads .. Delay is " + (StartJettySparqlEndpoint.requestDelay<=0?"disabled":StartJettySparqlEndpoint.requestDelay+"ms"));
				
		for (int i=0; i<nWorkers; i++) {
			WorkerThread t = new WorkerThread();
			try {
				t.init();
			} catch (Exception e) {
				log.fatal("Error initializing worker thread: ", e);
				System.exit(1);
			}
			t.start();
			
			workers.add(t);
		}
		
		if (log.isDebugEnabled()) {
			log.info("Registering IdleStatusMonitor. Total number of workers: " + nWorkers);
			new IdleWorkersMonitor().start();
		}
	}
	
	
	protected void initializeRepository() {
		
		repo = getRepository();
		
		log.info("Calling initialize on repository .. this may take some minutes.");
		try {
			repo.initialize();
		} catch (RepositoryException e) {
			log.fatal("Error initializing repository.");
			throw new RuntimeException(e);
		}

		log.info("Repository successfully initialized.");
		
	}
	
	
	
	protected SailRepository getRepository() {
		SailRepository res = null;
		
		String type = StartJettySparqlEndpoint.repositoryType;
		String loc = StartJettySparqlEndpoint.repositoryLocation;
		
		if (type==null)
			type = "native";
		
		if (type.equals("native")) {
			log.info("Initializing instance with native repository at " + loc);
			res = new SailRepository( new NativeStore(new File(loc) ));
			
//			res = new SailRepository( getNativeStore(new File(loc), "spoc,psoc") );
			log.info("Repository initialized.");
		} else {
			throw new RuntimeException("Type not supported yet: " + type);
		}
		
		return res;
	}
	
	
	/**
     * Get a Native Store with better shutdown behaviour when any active
     * connection objects aren't properly closed Introduced since changing
     * graceful shutdown default timeout of 20 seconds is buggy in sesame (issue
     * Tracker http://www.openrdf.org/issues/browse/SES-673)
     */
    protected static NativeStore getNativeStore(File file, String indices)  {
        
    	return new NativeStore(file, indices) {
            ArrayList<NativeStoreConnection> activeCon = new ArrayList<NativeStoreConnection>();

            @Override
            protected NotifyingSailConnection getConnectionInternal() throws SailException {
                NativeStoreConnection con = (NativeStoreConnection) super.getConnectionInternal();
                activeCon.add(con);
                return con;
            }

            @Override
            public void shutDown() throws SailException {
                for (NativeStoreConnection con : activeCon) {
                    con.close();
                }
                super.shutDown();
            }
        };
    }
    
   
    
    protected class QueryRequest {
    	
    	public final String query;
    	public final int requestID;
    	public final HttpServletRequest req;
    	public final HttpServletResponse resp;
    	public final ServletOutputStream outputStream;
    	public boolean done;
		
    	public QueryRequest(int requestID, String query, 
				HttpServletRequest req, HttpServletResponse resp,
				ServletOutputStream outputStream) {
			super();
			this.query = query;
			this.requestID = requestID;
			this.req = req;
			this.resp = resp;
			this.outputStream = outputStream;
			this.done = false;
		}
    	
    	public boolean isDone() {
    		synchronized (this) {
    			return done;
    		}
    	}
    	
    	public void done() {
    		synchronized (this) {
    			done = true;
    		}
    	}
    }
    
    
    
    protected class WorkerThread extends Thread {
    	
    	protected RepositoryConnection conn = null;
    	private final Logger log = Logger.getLogger(SparqlServlet2.class);
    	private BufferedOutputStream traceTo = null;
    	
    	public WorkerThread() {
    	}
    	
    	private void initTraceTo() throws IOException{
    		if (log.isTraceEnabled()) {
    			String fileName = "request_" + StartJettySparqlEndpoint.port + "_" + Thread.currentThread().getName() + ".req";
        		File file = new File("logs/req/" + fileName);
        		file.getParentFile().mkdirs();
        		log.trace("Tracing result of request to " + file.getAbsolutePath());
        		traceTo = new BufferedOutputStream( new FileOutputStream(file, true));	// append
        		traceTo.write(("RepoLoc: " + StartJettySparqlEndpoint.repositoryLocation + "\r\n\r\n").getBytes("UTF-8"));
    		}
    	}
    	
    	@Override
    	public void run() {
    		
    		try {
				initTraceTo();
			} catch (IOException e) {
				throw new RuntimeException(e);
			}
			
    		while (true) {
    			
    			QueryRequest qr = getNextQueryRequestSynch();
				
    			try {
					if (qr!=null) {
						if (log.isTraceEnabled())
							log.trace("Processing request " + qr.requestID);
						else if (qr.requestID%10==1)
							log.info("Status Information: Current request is " + qr.requestID);	// log every 10 statement
						
						try {
							processQuery(qr.query, qr.requestID, qr.req, qr.resp, qr.outputStream);
						} catch (Exception e) {
							log.error("Error handling query " + qr.requestID , e);
							
							error(qr.resp, 500, "Error handling query " + qr.requestID + ": " + e.getMessage());
							
						} finally {
							qr.done();
			    			synchronized (qr.resp) {
			    				qr.resp.notify();
			    			}
						}
						
		    			if (log.isTraceEnabled())
		    				log.trace("Processing request " + qr.requestID + " done.");
		    			continue;	// check if there are more requests before wait()
		    		}
					   			
					// wait on signal on the queue => push
					synchronized (queryRequestQueue) {
						try {
							if (!queryRequestQueue.isEmpty())
								continue;	// check too guarantee correctness
							log.trace("Waiting on notify call on queryRequestQueue. Falling asleep ...");
							idleWorkers++;
							queryRequestQueue.wait();
							idleWorkers--;
						} catch (InterruptedException e) {
							// ignore
						}
					}
    			} catch (Exception e) {
    				if (qr!=null) {
    					log.error("Unexpected exception occured. Rescheduling query task.", e);
    				
    					synchronized (queryRequestQueue) {
    						queryRequestQueue.addLast(qr);
    					}
    				}
    			}
    			log.trace("Awake again, will check if new requests are available.");
    		}
    		
    	}
    	
    	public void init() throws Exception {
    		conn = repo.getConnection();
    	}
    	
    	private void processQuery(String query, int reqId, HttpServletRequest req, HttpServletResponse resp, ServletOutputStream out)	{
            
    		 		
    		OutputStream outputStream = out;
            try {	
            	
            	if (log.isTraceEnabled()) {
            		String _log = "Query (#" + reqId + "): \r\n" + query;
            		log.trace(_log);
            		
            		// trace result to file '%port%_request%reqId%
            		traceTo.write(_log.getBytes("UTF-8"));
            		traceTo.write("\r\n\r\n".getBytes("UTF-8"));
            		traceTo.flush();
            		outputStream = new TracingOutputStream(out, traceTo);
            	}
            		        	
            	query = query.trim();
            	
    	        Query result;
    	        if (query.startsWith("SELECT"))
    	            result = conn.prepareTupleQuery(QueryLanguage.SPARQL, query);
    	        else if (query.startsWith("CONSTRUCT"))
    	            result = conn.prepareGraphQuery(QueryLanguage.SPARQL, query);
    	        else if (query.startsWith("ASK"))
    	            result = conn.prepareBooleanQuery(QueryLanguage.SPARQL, query);
    	        else
    	        	result = conn.prepareQuery(QueryLanguage.SPARQL, query);
    	            
    	            
    	        if (result instanceof BooleanQuery){
    	            BooleanQuery bQuery = (BooleanQuery) result;
    	            boolean res = bQuery.evaluate();

    	            FileFormatServiceRegistry<? extends FileFormat, ?> registry = BooleanQueryResultWriterRegistry.getInstance();
    	            BooleanQueryResultWriterFactory qrWriterFactory = (BooleanQueryResultWriterFactory)ProtocolUtil.getAcceptableService(req, resp, registry);

    	            resp.setStatus(HttpServletResponse.SC_OK);
    	            BooleanQueryResultWriter qrWriter = qrWriterFactory.getWriter(outputStream);
    	            qrWriter.write(res);
    	        }
    	        else if (result instanceof TupleQuery){
    	            
    	            TupleQuery tQuery = (TupleQuery)result;
    	            TupleQueryResult res = tQuery.evaluate();
                       
                    FileFormatServiceRegistry<? extends FileFormat, ?> registry = TupleQueryResultWriterRegistry.getInstance();
                    TupleQueryResultWriterFactory qrWriterFactory = (TupleQueryResultWriterFactory)ProtocolUtil.getAcceptableService(req, resp, registry);
                    
                    resp.setContentType("application/sparql-results+xml");
                    resp.setStatus(HttpServletResponse.SC_OK);
                    TupleQueryResultWriter qrWriter = qrWriterFactory.getWriter(outputStream);
                    QueryResultUtil.report(res, qrWriter);
                }
    	        else if (result instanceof GraphQuery){
    	            GraphQuery gQuery = (GraphQuery)result;
    	            GraphQueryResult res = gQuery.evaluate();

    	            FileFormatServiceRegistry<? extends FileFormat, ?> registry = RDFWriterRegistry.getInstance();
    	            RDFWriterFactory qrWriterFactory = (RDFWriterFactory)ProtocolUtil.getAcceptableService(req, resp, registry);

    	            resp.setStatus(HttpServletResponse.SC_OK);
    	            resp.setContentType("application/x-trig");
    	            RDFWriter qrWriter = qrWriterFactory.getWriter(outputStream);
    	            
    	            QueryResultUtil.report(res, qrWriter);
    	        }
    	        
    	        outputStream.flush();
    	        resp.flushBuffer();
    	        outputStream.close();
 
    	    }        
    	    catch (Exception e) {

    			log.error("Error occured while processing the query. \nQuery:" + query, e);
    			error(resp, 500, "Error occured while processing the query. <p>" + query + "<p>" + e.getClass().getSimpleName() + ": " + e.getMessage());
    	    	
    	    } 
    	}
    }
    
    
    protected class TracingOutputStream extends OutputStream {

    	private final OutputStream out;
    	private final OutputStream traceTo;
    	private int written = 0;
    	
    	public TracingOutputStream(OutputStream out, OutputStream traceTo) {
    		this.out = out;
    		this.traceTo = traceTo;
    	}
    	
		@Override
		public void write(int b) throws IOException {
			out.write(b);	
			traceTo.write(b);
			written++;
			flushIf();			
		}
    	
		@Override
		public void flush() throws IOException {
			out.flush();
			traceTo.flush();
		}
		
		@Override
		public void write(byte[] b, int off, int len) throws IOException {
			out.write(b, off, len);
			traceTo.write(b, off, len);
			written += len;
			flushIf();
		}
		
		@Override
		public void write(byte[] b) throws IOException {
			out.write(b);
			traceTo.write(b);
			written += b.length;
			flushIf();
		}
		
		@Override
		public void close() throws IOException {
			traceTo.flush();	
			// do not close traceTo at is still needed!
			out.close();
		}
		
		
		protected void flushIf() throws IOException {
			if (written>4096) {
				flush();
				written=0;
			}
		}
    }
    
    
    protected class IdleWorkersMonitor extends Thread {
    	
    	@Override
    	public void run() {
    		
    		while (!Thread.interrupted()) {
    			
    			int _idle;
    			int req;
    			StringBuilder sb = new StringBuilder();
    			synchronized (queryRequestQueue) {
    				_idle = idleWorkers;
    				req = queryRequestQueue.size();
    				for (QueryRequest q : activeQueries)
        				sb.append(q.requestID).append(";");
    			}
    			    			
    			System.out.println("Worker Status: " + _idle + " idle, requests in queue: " + req + ", active requests: " + sb.toString());
    			
    			try {
					Thread.sleep(5000);
				} catch (InterruptedException e) {
					// ignore
				}
    		}
    	}
    }
    
}
