package org.semanticweb.fbench.proxy;

import java.io.IOException;

import javax.servlet.ServletException;
import javax.servlet.ServletRequest;
import javax.servlet.ServletResponse;

public interface RequestHandler {

	/**
	 * Handle an incoming request to this system. Returns true if the original
	 * processor shall take care for processing the actual request, which 
	 * is the case in most situations.
	 * 
	 * This method is a callback to intercept requests, and for instance add
	 * a delay.
	 * 
	 * @param req
	 * @param res
	 * @return
	 * 		true, if the original processor shall handle this request
	 * 
	 * @throws ServletException
	 * @throws IOException
	 */
	public boolean handleRequest(ServletRequest req, ServletResponse res) throws ServletException, IOException;
}
