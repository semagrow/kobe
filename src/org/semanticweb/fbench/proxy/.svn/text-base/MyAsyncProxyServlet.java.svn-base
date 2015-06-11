package org.semanticweb.fbench.proxy;

import java.io.IOException;
import java.io.UnsupportedEncodingException;
import java.net.MalformedURLException;
import java.net.URL;
import java.net.URLDecoder;

import javax.servlet.ServletException;
import javax.servlet.ServletRequest;
import javax.servlet.ServletResponse;

import org.apache.log4j.Logger;
import org.eclipse.jetty.http.HttpURI;
import org.eclipse.jetty.servlets.ProxyServlet;


/**
 * ProxyServlet to forward incoming requests.
 * 
 * http://localhost:2000/http://myEndpoint.com:80/sparql
 * 
 * In the above example the request is forwarded to 
 * http://myEndpoint.com:80/sparql. Note that in the general case
 * this URL is urlEncoded.
 * 
 * @author as
 *
 */
public class MyAsyncProxyServlet extends ProxyServlet {

	public static Logger log = Logger.getLogger(MyAsyncProxyServlet.class);
	
	protected RequestHandler reqHandler = null;
	
	public MyAsyncProxyServlet() {
		String clazz = StartJettyProxy.requestHandler;
		log.info("Servlet is being initialized. Registering " + clazz  + " as request handler.");
		try {
			reqHandler = (RequestHandler)Class.forName(clazz).newInstance();
		} catch (Exception e) {
			log.fatal("Servlet could not be initialized.", e);
			throw new RuntimeException(e);
		} 
	}
	
	@Override
	protected HttpURI proxyHttpURI(String scheme, String serverName, int serverPort, String uri) throws MalformedURLException {
			
		URL forwardURL;
		
		try {
			forwardURL = new URL( URLDecoder.decode(uri.substring(1), "UTF-8"));
		} catch (UnsupportedEncodingException e) {
			log.error("Unexpected encoding exception: " + e.getMessage());
			throw new RuntimeException(e);
		}
	
		
		int port = forwardURL.getPort() == -1 ? 80 : forwardURL.getPort();
		
		if (log.isTraceEnabled())
			log.trace("Incoming forward request for " + forwardURL);
		
		return new HttpURI(forwardURL.getProtocol()+"://"+forwardURL.getHost()+":"+port+ forwardURL.getFile());
	}
	
	
	@Override
	public void service(ServletRequest req, ServletResponse res) throws ServletException, IOException {
		boolean superProcess = reqHandler.handleRequest(req, res);
		if (superProcess)
			super.service(req, res);
	}

}
