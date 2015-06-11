package org.semanticweb.fbench.sparqlendpoint;

import java.io.File;
import java.io.IOException;

import org.apache.log4j.Logger;
import org.eclipse.jetty.server.Connector;
import org.eclipse.jetty.server.Server;
import org.eclipse.jetty.server.nio.SelectChannelConnector;
import org.eclipse.jetty.webapp.WebAppContext;
import org.semanticweb.fbench.misc.TimedInterrupt;
import org.semanticweb.fbench.misc.Utils;


public class StartJettySparqlEndpoint {
	
	private static Logger log;
	
	public static String repositoryType = null;
	public static String repositoryLocation = null;
	public static int port = 0;
	public static int nWorkerThreads = 8;
	public static int requestDelay = -1;	// a request delay which is added to each request, -1=disabled
	protected static Server server = null;
	protected static File pidFile = null;
	
	/**
	 * @param args
	 */
	public static void main(String[] args) throws Exception {
		
		if (args.length!=2 && args.length!=3 && args.length!=4)
			printHelpAndExit();
		
		repositoryType = "native";
		repositoryLocation = args[0];
			
		port = Integer.parseInt(args[1]);
		if (args.length>2)
			requestDelay = Integer.parseInt(args[2]);	// see SparqlServlet2.handleQuery	
		
		if (args.length==4)
			nWorkerThreads = Integer.parseInt(args[3]);
		String host = "localhost";
		
		// the id for the log4j log file => logs/sparql_%port%.log
		System.setProperty("serverId", Integer.toString(port));
		
		if (System.getProperty("log4j.configuration")==null)
			System.setProperty("log4j.configuration", "file:config/log4j-sparql.properties");
		log = Logger.getLogger(StartJettySparqlEndpoint.class);
		
		
		writePIDFile();	// write a file of %PID%.pid such that process can be killed if it does not terminate
		
        new GracefullShutdownThread().start();
        
        /*
         * Hopefully these settings fix some rarely occurring bugs:
         * http://wiki.eclipse.org/Jetty/Feature/JVM_NIO_Bug
         */
        System.setProperty("org.mortbay.io.nio.MAX_SELECTS", "50000");
        System.setProperty("org.mortbay.io.nio.BUSY_KEY", "10");
        System.setProperty("org.mortbay.io.nio.BUSY_PAUSE", "100");
        
		server = new Server();
        Connector connector = new SelectChannelConnector();
        connector.setPort(port);
        connector.setHost(host);
        server.addConnector(connector);

       
        WebAppContext wac = new WebAppContext();
        wac.setContextPath("/");
        wac.setWar("config/jetty/sparql/");
        server.setHandler(wac);
        server.setStopAtShutdown(true);

        server.start();		
        
	}

	
	protected static void printHelpAndExit() {
		System.out.println("Usage: \n" +
				"\tstartSparqlEndpoint <RepositoryLocation> <Port>\n" +
				"\tstartSparqlEndpoint <RepositoryLocation> <Port> <Delay>\n" +
				"\tstartSparqlEndpoint <RepositoryLocation> <Port> <Delay> <WorkerThreads>\n");
		System.exit(1);
	}
	
	
	protected static void writePIDFile() throws IOException {
		
		long pid = Utils.getPID();
		pidFile = new File(pid + ".pid");
		pidFile.createNewFile();
	}
	
	protected static class GracefullShutdownThread extends Thread {
		
		private File cFile = new File("_shutdown");
		
		@Override
		public void run() {
			
			log.info("Started GracefullShutdownThread...");
			while (true) {
				
				try {
					Thread.sleep(250);
				} catch (InterruptedException e) {
					
				}
				if (cFile.exists()) {
					log.info("Gracefull shutdown request ... ");
					try {
						boolean success = new TimedInterrupt().run( new Runnable() {
							@Override
							public void run() {
								try {
									server.stop();
								} catch (Exception e) {
									log.error("Error closing conenction.", e);
									exit(1);
								}						
							}
						}, 30000);
					
						if (!success)
							exit(1);
										
					} catch (Exception e) {
						log.error("Stopping the server failed.", e);
						System.exit(1);
					}
					exit(0);
				}
			}
		}
		
		
		protected void exit(int code) {
			if (pidFile!=null)
				pidFile.delete();
			System.exit(code);
		}
	}
}
