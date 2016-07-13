package org.semanticweb.fbench.owlim;

import java.io.File;
import java.util.Iterator;
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
import org.openrdf.sail.Sail;
import org.semanticweb.fbench.misc.FileUtil;
import org.semanticweb.fbench.provider.RepositoryProvider;


/**
 * Provider to fill a native OWLim store.<p>
 * 
 * Sample dataConfig:<p>
 * 
 * <code>
 * 
 * relative Path for storeFile (relative to baseDir)
 * 
 * <http://NYTimes.Locations> fluid:store "org.semanticweb.fbench.owlim.OwlimRepositoryFiller";
 * fluid:rdfFile "D:\\datasets\\nytimes\\locations.rdf";
 * fluid:RepositoryLocation "data\\OwlimManager";
 * fluid:context <http://nytimes.org>.
 * 
 * 
 * absolute Path for repo location
 * 
 * <http://NYTimes.Locations> fluid:store "org.semanticweb.fbench.owlim.OwlimRepositoryFiller";
 * fluid:rdfFile "D:\\datasets\\nytimes\\locations.rdf";
 * fluid:RepositoryLocation "D:\\data\\OwlimManager";
 * fluid:context <http://nytimes.org>.
 * </code>
 * 
 * @author (mz), as
 *
 */
public class OwlimRepositoryFiller implements RepositoryProvider {
	
	@Override
	public Repository load(Graph graph, Resource repNode) throws Exception{
		SailRepository rep;
		
		Iterator<Statement> iter = graph.match(repNode, new URIImpl("http://fluidops.org/config#rdfFile"), null);
		Statement s = iter.next();
		String fileName = s.getObject().stringValue();
		iter = graph.match(repNode, new URIImpl("http://fluidops.org/config#RepositoryLocation"), null);
		s = iter.next();
		String repoLocationString = s.getObject().stringValue();
		iter = graph.match(repNode, new URIImpl("http://fluidops.org/config#context"), null);
		s = iter.next();
		String context = s.getObject().toString();
		
		File rdfFile = FileUtil.getFileLocation(fileName);
		if (!rdfFile.exists())
			throw new RuntimeException("RDF file does not exist at '" + fileName + "'.");
		
		File store = FileUtil.getFileLocation(repoLocationString);
		
		System.setProperty("ruleset", "empty"); // for performance reasons
        System.setProperty("repository-type", "weighted-file-repository");
		System.setProperty("console-thread", "false");
		System.setProperty("storage-folder", store.getAbsolutePath());
	
		// we initialize with empty ruleset
		rep = new SailRepository((Sail) Class
				.forName("com.ontotext.trree.OwlimSchemaRepository")
				.getConstructor().newInstance());
		rep.initialize();
    		
		RDFFormat rdfFormat = RDFFormat.forFileName(rdfFile.getName());
		URI u = ValueFactoryImpl.getInstance().createURI(context);
		System.out.println("Adding dataset under context " + u.toString());
		if (rdfFormat != null){
			RepositoryConnection conn = rep.getConnection();
		    try {
		    	conn.add(rdfFile, null, rdfFormat, u);
		    }
		    finally {
		    	conn.close();
		    }
		}
		rep.shutDown();
		
		
		return null;
	}

	@Override
	public String getLocation(Graph graph, Resource repNode) {
		Iterator<Statement> iter = graph.match(repNode, new URIImpl("http://fluidops.org/config#rdfFile"), null);
		Statement s = iter.next();
		String fileName = s.getObject().stringValue();
		return fileName;
	}

	@Override
	public String getId(Graph graph, Resource repNode) {
		Iterator<Statement> iter = graph.match(repNode, new URIImpl("http://fluidops.org/config#RepositoryLocation"), null);
		Statement s = iter.next();
		String id = new File(s.getObject().stringValue()).getName();
		return id;
	}

}
