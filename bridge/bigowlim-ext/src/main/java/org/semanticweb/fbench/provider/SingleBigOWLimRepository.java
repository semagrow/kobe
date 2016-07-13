package org.semanticweb.fbench.provider;

import java.io.File;
import java.util.Iterator;

import org.openrdf.model.Graph;
import org.openrdf.model.Resource;
import org.openrdf.model.Statement;
import org.openrdf.model.impl.URIImpl;
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
 * <http://CrossStore> fluid:store "BigOWLimRepo";
 * fluid:RepositoryLocation "data\\owlim-storage.SingleStore.Cross".
 * 
 * 
 * absolute Path for repo location
 * <http://CrossStore> fluid:store "BigOWLimRepo";
 * fluid:RepositoryLocation "D:\\data\\owlim-storage.SingleStore.Cross".
 * </code>
 * 
 * @author (mz), as
 *
 */
public class SingleBigOWLimRepository implements RepositoryProvider {
	
	@Override
	public SailRepository load(Graph graph, Resource repNode) throws Exception {
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
		return s.getObject().stringValue();		// the repo location
	}

	@Override
	public String getId(Graph graph, Resource repNode) {
		Iterator<Statement> iter = graph.match(repNode, new URIImpl("http://fluidops.org/config#RepositoryLocation"), null);
		Statement s = iter.next();
		String id = new File(s.getObject().stringValue()).getName();
		return id;
	}
}
