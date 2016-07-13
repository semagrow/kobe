package org.semanticweb.fbench.owlim;

import java.io.File;
import java.util.Iterator;

import org.openrdf.model.Graph;
import org.openrdf.model.Resource;
import org.openrdf.model.Statement;
import org.openrdf.model.impl.URIImpl;
import org.openrdf.repository.Repository;
import org.openrdf.repository.sail.SailRepository;
import org.openrdf.sail.Sail;

import org.semanticweb.fbench.misc.FileUtil;
import org.semanticweb.fbench.provider.RepositoryProvider;




/**
 * Provider for a native OWLIM repository.<p>
 * 
 * Sample dataConfig:<p>
 * 
 * <code>
 * 
 * relative Path for repo location (relative to baseDir)
 * 
 * <http://DBpedia> fluid:store "org.semanticweb.fbench.owlim.OwlimRepository";
 * fluid:RepositoryLocation "data\\repositories\\owlim-storage.dbpedia".
 * 
 * 
 * absolute Path
 * <http://DBpedia> fluid:store "org.semanticweb.fbench.owlim.OwlimRepository";
 * fluid:RepositoryLocation "C:\\data\\repositories\\owlim-storage.dbpedia".
 * </code>
 * 
 * @author (mz), as
 *
 */
public class OwlimRepository implements RepositoryProvider {
	
	@Override
	public Repository load(Graph graph, Resource repNode) throws Exception {
		Iterator<Statement> iter = graph.match(repNode, new URIImpl("http://fluidops.org/config#RepositoryLocation"), null);
		Statement s = iter.next();
		String repoLocation = s.getObject().stringValue();
		
		File store = FileUtil.getFileLocation(repoLocation);
		if (!store.exists()){
			throw new RuntimeException("Store does not exist at '" + repoLocation + "'.");
		}

		System.setProperty("ruleset", "empty"); // for performance reasons
        System.setProperty("repository-type", "weighted-file-repository");
		System.setProperty("console-thread", "false");
		System.setProperty("storage-folder", store.getAbsolutePath());
	
		// we initialize with empty ruleset
		SailRepository rep = new SailRepository((Sail) Class
				.forName("com.ontotext.trree.OwlimSchemaRepository")
				.getConstructor().newInstance());
		rep.initialize();
		
		return rep;
	}

	@Override
	public String getLocation(Graph graph, Resource repNode) {
		Iterator<Statement> iter = graph.match(repNode, new URIImpl("http://fluidops.org/config#RepositoryLocation"), null);
		Statement s = iter.next();
		String repoLocation = s.getObject().stringValue();
		return repoLocation;
	}

	@Override
	public String getId(Graph graph, Resource repNode) {
		Iterator<Statement> iter = graph.match(repNode, new URIImpl("http://fluidops.org/config#RepositoryLocation"), null);
		Statement s = iter.next();
		String id = new File(s.getObject().stringValue()).getName();
		return id;
	}
}
