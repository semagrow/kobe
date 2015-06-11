package org.semanticweb.fbench.provider;

import java.io.File;
import java.util.Iterator;

import org.openrdf.model.Graph;
import org.openrdf.model.Resource;
import org.openrdf.model.Statement;
import org.openrdf.model.impl.URIImpl;
import org.openrdf.repository.Repository;
import org.openrdf.repository.sail.SailRepository;
import org.semanticweb.fbench.misc.BigOWLimFileHandler;
import org.semanticweb.fbench.misc.FileUtil;

import com.ontotext.trree.OwlimSchemaRepository;



/**
 * Provider for a native BigOWLim repository.<p>
 * 
 * Sample dataConfig:<p>
 * 
 * <code>
 * 
 * relative Path for repo location (relative to baseDir)
 * 
 * <http://DBpedia> fluid:store "BigOWLimRepo";
 * fluid:RepositoryLocation "data\\repositories\\owlim-storage.dbpedia".
 * 
 * 
 * absolute Path
 * <http://DBpedia> fluid:store "BigOWLimRepo";
 * fluid:RepositoryLocation "C:\\data\\repositories\\owlim-storage.dbpedia".
 * </code>
 * 
 * @author (mz), as
 *
 */
public class BigOwlimRepository implements RepositoryProvider {
	
	@Override
	public Repository load(Graph graph, Resource repNode) throws Exception {
		Iterator<Statement> iter = graph.match(repNode, new URIImpl("http://fluidops.org/config#RepositoryLocation"), null);
		Statement s = iter.next();
		String repoLocationString = s.getObject().stringValue();
		
		BigOWLimFileHandler repoLocation = new BigOWLimFileHandler(FileUtil.getFileLocation(repoLocationString));
		
		OwlimSchemaRepository osr = new OwlimSchemaRepository();
		osr.setParameter("ruleset", "empty");
		osr.setParameter("console-thread", "false");
		osr.setParameter("storage-folder", repoLocation.getRepoLocationName());
		osr.setDataDir(repoLocation.getBaseDir());
		SailRepository rep = new SailRepository(osr);
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
