package org.semanticweb.fbench.proxy;

import java.io.IOException;

import javax.servlet.ServletException;
import javax.servlet.ServletRequest;
import javax.servlet.ServletResponse;
import javax.servlet.http.HttpServletRequest;

import org.apache.log4j.Logger;

public class LogRequestHandler implements RequestHandler {

	public static Logger log = Logger.getLogger(LogRequestHandler.class);
	
	private final static int COUNTER_MAX = 10;
	private int count = 0;
	
	@Override
	public boolean handleRequest(ServletRequest req, ServletResponse res)
			throws ServletException, IOException {
		
		if (log.isDebugEnabled())
			log.debug("Incoming request for " + ((HttpServletRequest)req).getRequestURI() );
		
		count++; 
		if (count%COUNTER_MAX==0)
			log.info("handling request " + count);
				
		return true;
	}

}
