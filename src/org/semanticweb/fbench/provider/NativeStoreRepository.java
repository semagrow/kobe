package org.semanticweb.fbench.provider;

import java.io.File;
import java.util.Iterator;

import org.openrdf.model.Graph;
import org.openrdf.model.Resource;
import org.openrdf.model.Statement;
import org.openrdf.model.impl.URIImpl;
import org.openrdf.repository.Repository;
import org.openrdf.repository.sail.SailRepository;
import org.openrdf.sail.nativerdf.NativeStore;
import org.semanticweb.fbench.misc.FileUtil;



/**
 * Provider for a native sesame repository.<p>
 * 
 * Sample dataConfig:<p>
 * 
 * <code>
 * 
 * relative Path for repo location (relative to baseDir)
 * 
 * <http://DBpedia> fluid:store "NativeStore";
 * fluid:RepositoryLocation "data\\repositories\\native-storage.dbpedia".
 * 
 * 
 * absolute Path
 * <http://DBpedia> fluid:store "NativeStore";
 * fluid:RepositoryLocation "D:\\data\\repositories\\native-storage.dbpedia".
 * </code>
 * 
 * @author (mz), as
 *
 */
public class NativeStoreRepository implements RepositoryProvider {

	@Override
	public Repository load(Graph graph, Resource repNode) throws Exception {
		Iterator<Statement> iter = graph.match(repNode, new URIImpl("http://fluidops.org/config#RepositoryLocation"), null);
		Statement s = iter.next();
		
		// retrieve the location from dataConfig, first check absolute location
		String repoLocation = s.getObject().stringValue();
		
		File store = FileUtil.getFileLocation(repoLocation);
		if (!store.exists()){
			throw new RuntimeException("Store does not exist at '" + repoLocation + "'.");
		}

		NativeStore ns = new NativeStore(store);
		SailRepository rep = new SailRepository(ns);
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
