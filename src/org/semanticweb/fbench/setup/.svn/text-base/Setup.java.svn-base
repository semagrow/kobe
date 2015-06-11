package org.semanticweb.fbench.setup;

import java.io.File;
import java.io.FileReader;
import java.net.URL;
import java.util.Iterator;

import org.apache.log4j.Logger;
import org.openrdf.model.Graph;
import org.openrdf.model.Resource;
import org.openrdf.model.Statement;
import org.openrdf.model.Value;
import org.openrdf.model.impl.GraphImpl;
import org.openrdf.model.impl.URIImpl;
import org.openrdf.rio.RDFFormat;
import org.openrdf.rio.RDFHandler;
import org.openrdf.rio.RDFHandlerException;
import org.openrdf.rio.RDFParser;
import org.openrdf.rio.Rio;
import org.semanticweb.fbench.Config;
import org.semanticweb.fbench.misc.FileUtil;


/**
 * Generic setup class to prepare datasources according to specified configuration file.
 * This tool is able to download files from the web, e.g. via HTTP and to use a {@link ClassMediator}
 * for copying/unzipping or customized subtasks.
 * 
 * Sample config:
 * 
 * <code>
 * <http://NYtimes> fluid:localSource "D:\\Test\\people.rdf";
 * fluid:webSource "http://data.nytimes.com/people.rdf";
 * fluid:destPath "data\\rdf\\nytimes\\people.rdf";
 * fluid:classMediator "NULL" .
 * 
 * <http://DBpedia.Instance-Types> fluid:localSource "D:\\Test\\instance_types_en.nt.bz2";
 * fluid:webSource "http://downloads.dbpedia.org/3.5.1/en/instance_types_en.nt.bz2";
 * fluid:destPath "data\\rdf\\dbpedia351\\instance_types_en.nt";
 * fluid:classMediator "org.semanticweb.fbench.setup.BZipClassMediator" .
 * </code>
 * 
 * 
 * @author Andreas
 *
 */
public class Setup {

	public static Logger log = Logger.getLogger(Setup.class);
	
	public static void prepareDataSources() throws Exception {
		final Graph graph = new GraphImpl();
		final String dataConfig = Config.getConfig().getDataConfig();
		RDFParser parser = Rio.createParser(RDFFormat.N3);
		RDFHandler handler = new RDFHandler(){

			@Override
			public void endRDF() throws RDFHandlerException {						
			}

			@Override
			public void handleComment(String arg0){
			}

			@Override
			public void handleNamespace(String arg0, String arg1){
			}

			@Override
			public void handleStatement(Statement arg0)
					throws RDFHandlerException {
				graph.add(arg0);
			}

			@Override
			public void startRDF() throws RDFHandlerException {
			}
			
		};
		
		parser.setRDFHandler(handler);
		
		parser.parse(new FileReader(dataConfig), "http://fluidops.org/config#");
		Iterator<Statement> iter = graph.match(null, new URIImpl("http://fluidops.org/config#localSource"), null);
		Statement s;
		Resource repNode;
		Value repType;
		boolean error = false;
		while (iter.hasNext()){
			s = iter.next();
			repNode = s.getSubject();
			repType = s.getObject();
			
			String localSource = repType.stringValue();
			String webSource = getObjectValue(graph, repNode, "http://fluidops.org/config#webSource");
			String destPath = getObjectValue(graph, repNode, "http://fluidops.org/config#destPath");
			String classMediator = getObjectValue(graph, repNode, "http://fluidops.org/config#classMediator");
			
			try {
				downloadAndMove(repNode.stringValue(), localSource, webSource, destPath, classMediator);
			} catch (Exception e) {
				log.error("Error preparing source for " + repType.stringValue() + ". Please prepare source manually. ");
				log.error("Details (" + e.getClass().getSimpleName() + "): " + e.getMessage());
				log.debug("exception", e);
				error = true;
			}
		}
		
		if (error)
			log.info("At least one data source could not be prepared. Consult log statements and prepare the respective data source manually.");
		else
			log.info("Setup successfully finished.");
	}
	
	
	
	public static void downloadAndMove(String id, String localSource, String webSource, String destPath, String classMediator) throws Exception {
		
		File dest = FileUtil.getFileLocation(destPath);
		ClassMediator mediator = (classMediator.toLowerCase().equals("null")) ? new CopyClassMediator() : (ClassMediator)Class.forName(classMediator).newInstance();
		
		// check if the file exists at the local source
		// if so, copy it to destPath
		File local = FileUtil.getFileLocation(localSource);
		if (local.exists()) {
			log.info("Using local source for " + id + ". Copy task starting...");
			log.debug("Using local source for " + id + ": source=" + local.getAbsolutePath() + ", dest=" + dest.getAbsolutePath() + ", mediator=" + mediator.getClass().getCanonicalName());
			mediator.perform(local, dest);
			return;
		}
		
		// otherwise: try to download from web source. use tmp folder
		log.info("Using web source for " + id + ". Download task starting...");
		URL url = new URL(webSource);
		
		File tmp = new File("tmp", url.getPath());
		tmp.getParentFile().mkdirs();		
		
		FileUtil.download(url, tmp);
		
		log.info("Download done. Running copy task...");
		log.debug("Using web source for " + id + ": source=" + webSource + ", dest=" + dest.getAbsolutePath() + ", mediator=" + mediator.getClass().getCanonicalName());
		
		log.info("Copying downloaded file to local destination ...");
		FileUtil.copyFile(tmp, local);
		
		mediator.perform(tmp, dest);
	}
	
	
	/**
	 * convenience method to retrieve the object value for the given subject and predicate. 
	 * If not present a RuntimeException is thrown.
	 * 
	 * @param graph
	 * @param subject
	 * @param predicateURI
	 * @return
	 */
	protected static String getObjectValue(Graph graph, Resource subject, String predicateURI) {
		Iterator<Statement> iter = graph.match(subject, new URIImpl(predicateURI), null);
		if (!iter.hasNext())
			throw new RuntimeException("Expected pattern match for <" + subject.stringValue() + "; " + predicateURI + "; null>. However, no result found.");
		return iter.next().getObject().stringValue();
	}
}
