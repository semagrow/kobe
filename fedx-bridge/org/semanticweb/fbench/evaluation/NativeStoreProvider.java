package org.semanticweb.fbench.evaluation;

import java.io.File;
import java.util.Iterator;

import org.apache.log4j.Logger;
import org.openrdf.model.Graph;
import org.openrdf.model.Resource;
import org.openrdf.model.Statement;
import org.openrdf.model.impl.URIImpl;
import org.openrdf.repository.Repository;
import org.openrdf.repository.sail.SailRepository;
import org.openrdf.sail.nativerdf.NativeStore;
import org.openrdf.sail.nativerdf.NativeStoreExt;
import org.semanticweb.fbench.misc.FileUtil;
import org.semanticweb.fbench.provider.RepositoryProvider;



/**
 * Provider for a native sesame repository.<p>
 * 
 * This is a extension which tries to load the NativeStoreExt repository from FedX. The 
 * extension allows for evaluation of precompiled queries without prior optimization
 * overhead.
 * 
 * Sample dataConfig:<p>
 * 
 * <code>
 * 
 * relative Path for repo location (relative to baseDir)
 * 
 * <http://DBpedia> fluid:store "NativeRepo";
 * fluid:RepositoryLocation "data\\repositories\\native-storage.dbpedia".
 * 
 * 
 * absolute Path
 * <http://DBpedia> fluid:store "NativeRepo";
 * fluid:RepositoryLocation "D:\\data\\repositories\\native-storage.dbpedia".
 * </code>
 * 
 * @author (mz), as
 *
 */
public class NativeStoreProvider implements RepositoryProvider {

	public static Logger log = Logger.getLogger(NativeStoreProvider.class);
	
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
		
		NativeStore ns = new NativeStoreExt(store);
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
