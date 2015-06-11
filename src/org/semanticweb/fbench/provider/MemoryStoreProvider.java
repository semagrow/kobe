package org.semanticweb.fbench.provider;

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
import org.openrdf.sail.memory.MemoryStore;

public class MemoryStoreProvider implements RepositoryProvider {

private Repository rep;
	
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

		// TODO
		// add absolute / relative path check for fileName and repoLocation
		
		File repositoryDir = new File(repoLocation);
    	boolean exists = repositoryDir.exists();
		MemoryStore ms = new MemoryStore();
    	ms.setDataDir(repositoryDir);
    	ms.setPersist(true);
		rep = new SailRepository(ms);
		rep.initialize();
		if (!exists){
		    RDFFormat rdfFormat = RDFFormat.forFileName(fileName);
		    URI u = ValueFactoryImpl.getInstance().createURI(context);
		    System.out.println("Adding dataset under context " + u.toString());
		    if (rdfFormat != null){
		    	RepositoryConnection conn = rep.getConnection();
		    	try {
		    		conn.add(new File(fileName), null, rdfFormat, u);
		    	}
		    	finally {
		    		conn.close();
		    	}
		    }
		}

		return rep;	
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
		Iterator<Statement> iter = graph.match(repNode, new URIImpl("http://fluidops.org/config#rdfFile"), null);
		Statement s = iter.next();
		String id = "mem_" + new File(s.getObject().stringValue()).getName();
		return id;
	}

}
