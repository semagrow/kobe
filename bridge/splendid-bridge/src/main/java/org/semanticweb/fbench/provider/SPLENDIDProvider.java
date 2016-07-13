package org.semanticweb.fbench.provider;

import java.io.File;
import java.io.IOException;
import java.io.InputStream;
import java.net.URL;
import java.util.Iterator;

import org.openrdf.model.Graph;
import org.openrdf.model.Resource;
import org.openrdf.model.Statement;
import org.openrdf.model.impl.GraphImpl;
import org.openrdf.model.impl.URIImpl;
import org.openrdf.repository.Repository;
import org.openrdf.repository.config.RepositoryConfig;
import org.openrdf.repository.config.RepositoryFactory;
import org.openrdf.repository.config.RepositoryImplConfig;
import org.openrdf.repository.config.RepositoryRegistry;
import org.openrdf.rio.RDFFormat;
import org.openrdf.rio.RDFParser;
import org.openrdf.rio.Rio;
import org.openrdf.rio.helpers.StatementCollector;
import org.openrdf.sail.config.SailConfigException;

public class SPLENDIDProvider implements RepositoryProvider {
	
	@Override
	public Repository load(Graph graph, Resource repNode) throws Exception {
		Iterator<Statement> iter = graph.match(repNode, new URIImpl("http://fluidops.org/config#rdfFile"), null);
		String fileName = iter.next().getObject().stringValue();
		return load(new File(fileName).toURI().toURL());
	}
	
	@Override
	public String getId(Graph graph, Resource repNode) {
		String id = repNode.stringValue().replace("http://", "");
		return id.replace("/", "_");
	}

	@Override
	public String getLocation(Graph graph, Resource repNode) {
		return "SPLENDID: unknown location";
	}

	
	protected Repository load(URL url) {
		
		try{
			Graph model = parse(url);
			RepositoryConfig repConf = RepositoryConfig.create(model, null);
			repConf.validate();
			RepositoryImplConfig implConf = repConf.getRepositoryImplConfig();
			RepositoryRegistry registry = RepositoryRegistry.getInstance();
			RepositoryFactory factory = (RepositoryFactory)registry.get(implConf.getType());
			if (factory == null)
				throw new SailConfigException("Unsupported repository type: " + implConf.getType());
			
			Repository repository = factory.getRepository(implConf);
			repository.initialize();
			return repository;
			
		} catch (Exception e) {
			e.printStackTrace();
		}
		
		return null;
	}
	
	protected Graph parse(URL url) throws SailConfigException, IOException {
		
		RDFFormat format = Rio.getParserFormatForFileName(url.getFile());
		if (format==null)
			throw new SailConfigException("Unsupported file format: " + url.getFile().toString());
		RDFParser parser = Rio.createParser(format);
		Graph model = new GraphImpl();
		parser.setRDFHandler(new StatementCollector(model));
		InputStream stream = url.openStream();
		
		try {
			parser.parse(stream, url.toString());
		} catch (Exception e) {
			throw new SailConfigException("Error parsing file!");
		}
		
		stream.close();
		return model;
	}	
}