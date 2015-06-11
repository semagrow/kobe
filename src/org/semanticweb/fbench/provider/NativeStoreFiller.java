package org.semanticweb.fbench.provider;

import java.io.File;
import java.util.Iterator;

import org.apache.log4j.Logger;
import org.openrdf.model.Graph;
import org.openrdf.model.Resource;
import org.openrdf.model.Statement;
import org.openrdf.model.URI;
import org.openrdf.model.impl.URIImpl;
import org.openrdf.model.impl.ValueFactoryImpl;
import org.openrdf.repository.Repository;
import org.openrdf.repository.RepositoryConnection;
import org.openrdf.repository.sail.SailRepository;
import org.openrdf.rio.RDFFormat;
import org.openrdf.sail.nativerdf.NativeStore;
import org.semanticweb.fbench.misc.FileUtil;



/**
 * Provider to fill a native Sesame store.<p>
 * 
 * Sample dataConfig:<p>
 * 
 * <code>
 * 
 * relative Path for storeFile
 * 
 * <http://NYTimes.Locations> fluid:store "Native";
 * fluid:RepositoryLocation "data\\native-storage.SingleStore.Cross";
 * fluid:rdfFile "D:\\datasets\\nytimes\\locations.rdf";
 * fluid:context <http://nytimes.org>.
 * 
 * 
 * absolute Path for storeFile:
 * 
 * <http://NYTimes.Organizations> fluid:store "Native";
 * fluid:RepositoryLocation "D:\\data\\native-storage.SingleStore.Cross";
 * fluid:rdfFile "D:\\datasets\\nytimes\\organizations.rdf";
 * fluid:context <http://nytimes.org>.
 * </code>
 * 
 * Note: if rdfFile is a directory, the tool adds all rdf files that
 * are contained within the directory.
 * 
 * <code>
 * <http://DBpedia.Properties> fluid:store "Native";
 * fluid:RepositoryLocation "data\\repositories\\native-storage.dbpedia351";
 * fluid:rdfFile "data\\rdf\\dbpedia351\\mapingbased_properties_en_chunked_cleaned\\";
 * fluid:context <http://DBpedia.org>.
 * </code>
 * 
 * By default the system tries to guess the RDFFormat from the file's extension.
 * However, in some cases this heuristic might fail. In such a case the RDFFormat
 * can be specified manually, as illustrated in the following example. The value
 * can be any of the {@link RDFFormat} enum type, e.g RDF/XML or N-TRIPLES
 * 
 * <code>
 * <http://Geonames> fluid:store "Native";
 * fluid:RepositoryLocation "data\\repositories\\native-storage.geonames";
 * fluid:rdfFile "data\\rdf\\geonames\\all-geonames-rdf.txt";
 * fluid:rdfFormat "RDF/XML";
 * fluid:context <http://Geonames.org>.
 * </code>
 * 
 * Specialized readers might be necessary for certain data formats. These can be
 * specified using the dataReader property. An example is given below. The property
 * must be set to a fully qualified class name of a DataReader implementation. Note
 * that this class must be on the classpath at runtime.
 * 
 * <code>
 * <http://Geonames> fluid:store "Native";
 * fluid:RepositoryLocation "data\\repositories\\native-storage.geonames";
 * fluid:rdfFile "data\\rdf\\geonames\\all-geonames-rdf.txt";
 * fluid:dataReader "RDF/XML";
 * fluid:context <http://Geonames.org>.
 * </code>
 * 
 * @author (mz), as
 *
 */
public class NativeStoreFiller implements RepositoryProvider {

	public static Logger log = Logger.getLogger(NativeStoreFiller.class);
	
	@Override
	public Repository load(Graph graph, Resource repNode) throws Exception {
		Iterator<Statement> iter = graph.match(repNode, new URIImpl("http://fluidops.org/config#rdfFile"), null);
		Statement s = iter.next();
		String fileName = s.getObject().stringValue();
		iter = graph.match(repNode, new URIImpl("http://fluidops.org/config#RepositoryLocation"), null);
		s = iter.next();
		String repoLocation = s.getObject().stringValue();
		iter = graph.match(repNode, new URIImpl("http://fluidops.org/config#context"), null);
		s = iter.next();
		String context = s.getObject().toString();
		
		File rdfFile = FileUtil.getFileLocation(fileName);
		if (!rdfFile.exists())
			throw new RuntimeException("RDF file does not exist at '" + fileName + "'.");
		
		File store = FileUtil.getFileLocation(repoLocation);
		NativeStore ns = new NativeStore(store);		
		Repository rep = new SailRepository(ns);
		rep.initialize();
		
		RepositoryConnection conn = rep.getConnection();
		conn.setAutoCommit(false);	
		
		RDFFormat rdfFormat = getSpecifiedRdfFormat(graph, repNode);	// can still be null if not specified
		
		URI contextURI = ValueFactoryImpl.getInstance().createURI(context);
    	
		try {
			DataReader dataReader = hasCustomReader(graph, repNode);		// can be null
			
			long beforeTriples = conn.size();
			
			if (rdfFile.isDirectory()) {
				log.info("Adding contents of provided rdf directory " + rdfFile.getAbsolutePath());
				for (File f : rdfFile.listFiles()) {
					if (!f.isDirectory())
						addData(conn, f, rdfFormat, contextURI, dataReader);
				}
			} else {
				addData(conn, rdfFile, rdfFormat, contextURI, dataReader);
			}
			
			conn.commit();
			long afterTriples = conn.size();
			log.info("Loaded " + (afterTriples-beforeTriples) + " triples into repository for rdfFile " + rdfFile.getName() + ".");
		} finally {
			conn.close();
			rep.shutDown();
			rep=null;
			System.gc();
		}

		return null;
	}

	@Override
	public String getLocation(Graph graph, Resource repNode) {
		Iterator<Statement> iter = graph.match(repNode, new URIImpl("http://fluidops.org/config#rdfFile"), null);
		Statement s = iter.next();
		String fileName = s.getObject().stringValue();
		return fileName;
	}


	protected RDFFormat getSpecifiedRdfFormat(Graph graph, Resource repNode) {
		Iterator<Statement> iter = graph.match(repNode, new URIImpl("http://fluidops.org/config#rdfFormat"), null);
		if (!iter.hasNext())
			return null;
		Statement s = iter.next();
		RDFFormat res = RDFFormat.valueOf(s.getObject().stringValue());
		if (res==null)
			throw new RuntimeException("Specified rdfFormat is not applicable as RDFFormat: " + s.getObject().stringValue());
		return res;
	}
	
	protected DataReader hasCustomReader(Graph graph, Resource repNode) throws ClassNotFoundException, InstantiationException, IllegalAccessException {
		Iterator<Statement> iter = graph.match(repNode, new URIImpl("http://fluidops.org/config#dataReader"), null);
		if (!iter.hasNext())
			return null;
		Statement s = iter.next();
		return (DataReader)Class.forName(s.getObject().stringValue()).newInstance();
	}
	
	protected void addData(RepositoryConnection conn, File rdfFile, RDFFormat rdfFormat, URI context, DataReader dataReader) throws Exception {
		if (dataReader!=null) {
			log.info("Adding contents of " + rdfFile.getName() + " using provided dataReader: " + dataReader.getClass().getCanonicalName());
			dataReader.loadData(conn, rdfFile, context);
			return;
		}
		
    	rdfFormat = rdfFormat == null ? RDFFormat.forFileName(rdfFile.getName()) : rdfFormat;
    	if (rdfFormat != null){
    		log.info("Adding dataset " + rdfFile.getName() + " under context " + context.toString());
    		conn.add(rdfFile, null, rdfFormat, context);
    	} else {
    		log.warn("RDF format could not be determined from fileName. Could not add data.");
    		throw new RuntimeException("RDF format could not be determined for " + rdfFile.getName() + ". Specification in ttl data configuration necessary.");
    	}
	}

	@Override
	public String getId(Graph graph, Resource repNode) {
		Iterator<Statement> iter = graph.match(repNode, new URIImpl("http://fluidops.org/config#RepositoryLocation"), null);
		Statement s = iter.next();
		String id = new File(s.getObject().stringValue()).getName();
		return id;
	}
}
