package org.semanticweb.fbench.proxy;

import org.eclipse.jetty.server.Connector;
import org.eclipse.jetty.server.Server;
import org.eclipse.jetty.server.nio.SelectChannelConnector;
import org.eclipse.jetty.webapp.WebAppContext;


/**
 * This executable starts a Jetty server that acts as proxy for our
 * use case scenario. 
 * 
 * In particular the configuration registers the MyAsyncProxyServlet
 * which forwards HTTP requests to the URL specified within the 
 * parameter.
 * 
 * Example:
 * 
 * http://localhost:2000/http://myEndpoint.com:80/sparql
 * 
 * In the above example the request is forwarded to 
 * http://myEndpoint.com:80/sparql. Note that in the general case
 * this URL is urlEncoded.
 * 
 * Usage:
 * 
 * startProxy
 * startProxy <Port>
 * startProxy <Port> <RequestDelay>
 * startProxy <Port> <RequestDelay> <RequestHandler> 
 * 
 * Params:
 * <RequestHandler> - the fully qualified class implementing a RequestHanlder
 * 						default: {@link DelayRequestHandler}
 * <RequestDelay> - the delay in ms (e.g. used in {@link DelayRequestHandler}
 * 						default: 100
 * <Port> - the port of the jetty server, default is 2000
 * 
 * @author as
 *
 */
public class StartJettyProxy {
	
	/**
	 * Default request delay in ms
	 */
	public static final long DEFAULT_REQUEST_DELAY = 100;	
	
	public static String requestHandler = null;
	public static long requestDelay = DEFAULT_REQUEST_DELAY;
	
	
	/**
	 * @param args
	 */
	public static void main(String[] args) throws Exception {
		
		int port = 2000;
		String host = "localhost";
		
		requestHandler = DelayRequestHandler.class.getCanonicalName();
		
		// request handler is specified
		if (args.length==1) {
			port = Integer.parseInt(args[0]);
		} 
		
		else if (args.length==2) {
			port = Integer.parseInt(args[0]);
			requestDelay = Long.parseLong(args[1]);
		}
		
		else if (args.length==3) {
			port = Integer.parseInt(args[0]);
			requestDelay = Long.parseLong(args[1]);
			requestHandler = args[2];
		}
		
		else {
			printHelpAndExit();
		}
		
		if (System.getProperty("log4j.configuration")==null)
			System.setProperty("log4j.configuration", "file:config/log4j-proxy.properties");
		
		System.setProperty("org.mortbay.io.nio.MAX_SELECTS", "50000");
        System.setProperty("org.mortbay.io.nio.BUSY_KEY", "10");
        System.setProperty("org.mortbay.io.nio.BUSY_PAUSE", "100");
			
		Server server = new Server();
        Connector connector = new SelectChannelConnector();
        connector.setPort(port);
        connector.setHost(host);
        connector.setMaxIdleTime(1*60*60);		// 3600s=1h
        server.addConnector(connector);

       
        WebAppContext wac = new WebAppContext();
        wac.setContextPath("/");
        wac.setWar("config/jetty/proxy/");
        server.setHandler(wac);
        server.setStopAtShutdown(true);

        server.start();	
	}

	
	protected static void printHelpAndExit() {
		System.out.println("Usage: \n" +
				"\tstartProxy\n" +
				"\tstartProxy <Port>\n" +
				"\tstartProxy <Port> <RequestDelay>\n" +
				"\tstartProxy <Port> <RequestDelay> <RequestHandler>\n");
		System.exit(1);
	}
}
