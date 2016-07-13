package org.semanticweb.fbench.proxy;

import java.io.IOException;

import javax.servlet.ServletException;
import javax.servlet.ServletRequest;
import javax.servlet.ServletResponse;


/**
 * Adds a specified delay to any request to this proxy. See {@link StartJettyProxy}
 * on how to specify this delay.
 * 
 * @author as
 */
public class DelayRequestHandler implements RequestHandler {

	@Override
	public boolean handleRequest(ServletRequest req, ServletResponse res)
			throws ServletException, IOException {

		long delay = StartJettyProxy.requestDelay;
		
		try {
			Thread.sleep(delay);
		} catch (InterruptedException e) {
			; // no-op
		}
		
		return true;
	}

}
