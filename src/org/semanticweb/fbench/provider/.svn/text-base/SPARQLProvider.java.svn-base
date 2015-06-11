package org.semanticweb.fbench.provider;

import java.net.URLEncoder;
import java.util.Iterator;

import org.apache.log4j.Logger;
import org.openrdf.model.Graph;
import org.openrdf.model.Resource;
import org.openrdf.model.Statement;
import org.openrdf.model.impl.URIImpl;
import org.openrdf.repository.Repository;
import org.openrdf.repository.sparql.SPARQLRepository;
import org.semanticweb.fbench.Config;


/**
 * Provider to integrate a SPARQL endpoint into a Sail.<p>
 * 
 * Sample dataConfig:<p>
 * 
 * <code>
 * <http://DBpedia> fluid:store "SPARQLEndpoint";
 * fluid:SPARQLEndpoint "http://dbpedia.org/sparql".
 * 
 * <http://NYtimes> fluid:store "SPARQLEndpoint";
 * fluid:SPARQLEndpoint "http://api.talis.com/stores/nytimes/services/sparql".
 * </code>
 * 
 * If a http://fluidops.org/config#proxyUrl is specified, this overwrites the 
 * global proxyUrl setting in config.prop specific for this endpoint:
 * 
 * <code>
 * <http://NYtimes> fluid:store "SPARQLEndpoint";
 * fluid:SPARQLEndpoint "http://api.talis.com/stores/nytimes/services/sparql";
 * fluid:proxyUrl "http://localhost:2000/".
 * </code>
 * 
 * It is also possible to explicitely enable/disable a proxy using 
 * http://fluidops.org/config#useProxy and a boolean value (true/false). If
 * no local proxyUrl is specified, the global proxyUrl is used. Note that
 * this settings can switch off proxy usage for a given endpoint.
 * 
 * <code>
 * local proxy URL:
 * 
 * <http://NYtimes> fluid:store "SPARQLEndpoint";
 * fluid:SPARQLEndpoint "http://api.talis.com/stores/nytimes/services/sparql";
 * fluid:proxyUrl "http://localhost:2000/";
 * fluid:useProxy "true" .
 * 
 * global proxy URL:
 * 
 * <http://NYtimes> fluid:store "SPARQLEndpoint";
 * fluid:SPARQLEndpoint "http://api.talis.com/stores/nytimes/services/sparql";
 * fluid:useProxy "true" .
 * </code>
 * 
 * @author (mz), as
 *
 */
public class SPARQLProvider implements RepositoryProvider {

	public static Logger log = Logger.getLogger(SPARQLProvider.class);
	
	@Override
	public Repository load(Graph graph, Resource repNode) throws Exception {
		Iterator<Statement> iter = graph.match(repNode, new URIImpl("http://fluidops.org/config#SPARQLEndpoint"), null);
		Statement s = iter.next();
		String sparqlEndpoint = s.getObject().stringValue();
		
		log.info("Registering SPARQL endpoint " + sparqlEndpoint);
		
		boolean useProxy = false;
		String proxyUrl = null;
				
		// check if optional proxy value is specified
		iter = graph.match(repNode, new URIImpl("http://fluidops.org/config#proxyUrl"), null);
		if (iter.hasNext()) {
			proxyUrl = iter.next().getObject().stringValue();
		} 
		
		iter = graph.match(repNode, new URIImpl("http://fluidops.org/config#useProxy"), null);
		if (iter.hasNext()) {
			useProxy = Boolean.parseBoolean( iter.next().getObject().stringValue());
		} else if (Config.getConfig().useGlobalProxy()) {
			useProxy = true;		// use global proxy if specified
		} else if (proxyUrl!=null) {
			useProxy = true;		// implicetly use specified proxy
		}
		
		if (useProxy && proxyUrl==null) {
			proxyUrl = Config.getConfig().getProxyUrl();
		}
		
		if (useProxy) {
			if (proxyUrl==null) {
				log.error("No proxyURL found though proxy is enabled.");
				throw new RuntimeException("No proxyUrl specified though proxies are enabled. Check your configuration.");
			} 
			
			log.info("Proxy enabled for " + sparqlEndpoint + ": " + proxyUrl);
			
			if (!proxyUrl.endsWith("/"))
				log.warn("ProxyUrl does not end with a '/'.");
			
			sparqlEndpoint = proxyUrl + URLEncoder.encode(sparqlEndpoint, "UTF-8");
		} 
				
		SPARQLRepository rep = new SPARQLRepository(sparqlEndpoint);
		rep.initialize();
	
		return rep;
	}

	@Override
	public String getLocation(Graph graph, Resource repNode) {
		Iterator<Statement> iter = graph.match(repNode, new URIImpl("http://fluidops.org/config#SPARQLEndpoint"), null);
		Statement s = iter.next();
		String sparqlEndpoint = s.getObject().stringValue();
		return sparqlEndpoint;
	}

	@Override
	public String getId(Graph graph, Resource repNode) {
		String id = repNode.stringValue().replace("http://", "");
		id = id.replace("/", "_");
		return "sparql_" + id;
	}

}
