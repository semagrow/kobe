package org.semanticweb.fbench.sparqlendpoint;

import info.aduna.lang.FileFormat;
import info.aduna.lang.service.FileFormatServiceRegistry;

import java.io.BufferedReader;
import java.io.File;
import java.io.IOException;
import java.io.InputStreamReader;
import java.net.URLDecoder;
import java.util.ArrayList;
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
import org.semanticweb.fbench.misc.TimedInterrupt;


@Deprecated 
public class SparqlServlet extends HttpServlet {

	private static final long serialVersionUID = 2627590629243739807L;
	
	protected static Repository repo = null;
	protected static Executor executor = Executors.newCachedThreadPool();
	protected static boolean initialized = false;
	
	public SparqlServlet() {
		initializeRepository();
	}
	
	private static final Logger log = Logger.getLogger(SparqlServlet.class);
	
	@Override
	protected void doPost(HttpServletRequest req, HttpServletResponse resp){
		//logger.log(Level.INFO, "POST Request to SPARQL Endpoint Servlet");
		//logger.setLevel(Level.OFF);
		try {
			ServletInputStream input = req.getInputStream();
			InputStreamReader in = new InputStreamReader(input);
			BufferedReader reader = new BufferedReader(in);
			String query = "";
			String tmp;
			while ((tmp = reader.readLine()) != null){
				query = query + tmp;
			}
			
			query = query.substring(6);
			query = URLDecoder.decode(query, "ISO-8859-1");
			
			ServletOutputStream outputStream = resp.getOutputStream();
            processQuery(query, req, resp, outputStream);
			outputStream.flush();
			outputStream.close();
			reader.close();
		} 
		catch (IOException e) {
			e.printStackTrace();
		}
	}
	
	@Override
	protected void doGet(HttpServletRequest req, HttpServletResponse resp){
		//logger.log(Level.INFO, "GET Request to SPARQL Endpoint Servlet");
		//logger.setLevel(Level.OFF);
		
		try {
			ServletOutputStream outputStream = resp.getOutputStream();		
			if (req.getParameter("query") != null){
				String query = req.getParameter("query");
				processQuery(query, req, resp, outputStream);
			}
			else {
				resp.setStatus(501);
				outputStream.println("You provided no query");
			}
			outputStream.flush();
			outputStream.close();
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
		super.destroy();
	}
	
	private void processQuery(String query, HttpServletRequest req, HttpServletResponse resp, ServletOutputStream outputStream) {
		processQuerySynch(query, req, resp, outputStream);
	}
	
	
	private void processQuerySynch(String query, HttpServletRequest req, HttpServletResponse resp, ServletOutputStream outputStream)	{
          
		RepositoryConnection conn = null;
        try {	 
        	
        	if (!isInitialized()) {
        		resp.setStatus(503);
				outputStream.print("Error occured while processing the query: repository is not initialized.");
				outputStream.flush();
				
				// changed msc: TRY to transmit a result in addition
	            FileFormatServiceRegistry<? extends FileFormat, ?> registry = BooleanQueryResultWriterRegistry.getInstance();
				BooleanQueryResultWriterFactory qrWriterFactory = (BooleanQueryResultWriterFactory)ProtocolUtil.getAcceptableService(req, resp, registry);
	            BooleanQueryResultWriter qrWriter = qrWriterFactory.getWriter(outputStream);
	            qrWriter.write(false);
	            
	            return;
        	}
        	
        	conn = repo.getConnection();
        	        	
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
	        conn.close();	// check if this is blocking, if yes, make it asynchronous

	    }        
	    catch (Exception e) {
	    	try {
	    		resp.setStatus(501);
				outputStream.print("Error occured while processing the query. " + e.getClass().getSimpleName() + ": " + e.getMessage());
				outputStream.flush();
				
			} catch (IOException e1) {
				// ignore
			} catch (Exception e2) {
				// ignore
			}
			
			log.error("Error occured while processing the query. Trying to close the connection asynchronously.");
			final RepositoryConnection _conn = conn;
	    	executor.execute( new Runnable() {
				@Override
				public void run() {
					
					new TimedInterrupt().run( new Runnable() {
						@Override
						public void run() {
							try {
								_conn.close();
								System.gc();
							} catch (RepositoryException e) {
								log.error("Error closing conenction.", e);
							}						
						}
					}, 10000);
					System.gc();					
				}
			});

	    } finally {

//	    	final RepositoryConnection _conn = conn;
//	    	executor.execute( new Runnable() {
//				@Override
//				public void run() {
//					
//					new TimedInterrupt().run( new Runnable() {
//						@Override
//						public void run() {
//							try {
//								_conn.close();
//								//	log.info("Connection successfully closed.");
//								System.gc();
//							} catch (RepositoryException e) {
//								log.error("Error closing conenction.", e);
//							}						
//						}
//					}, 1000);
//					System.gc();					
//				}
//			});
	    	
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
		
		initialized = true;
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
			res = new SailRepository( getNativeStore(new File(loc), "spoc,psoc") );
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
    
    
    private void setInitialize(boolean flag) {
    	synchronized (SparqlServlet.class) {
    		initialized = flag;
    	}
    }
    
    private boolean isInitialized() {
    	synchronized (SparqlServlet.class) {
    		return initialized;
    	}
    }
}
