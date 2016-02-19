/*
 * Copyright (C) 2008-2012, fluid Operations AG
 *
 * FedX is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 * 
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 * 
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package com.fluidops.fedx.util;

import java.io.File;
import java.io.FileReader;
import java.io.IOException;
import java.net.MalformedURLException;
import java.net.URL;
import java.util.ArrayList;
import java.util.Iterator;
import java.util.List;

import org.openrdf.model.Graph;
import org.openrdf.model.Resource;
import org.openrdf.model.Statement;
import org.openrdf.model.Value;
import org.openrdf.model.impl.GraphImpl;
import org.openrdf.model.impl.LiteralImpl;
import org.openrdf.model.impl.URIImpl;
import org.openrdf.rio.RDFFormat;
import org.openrdf.rio.RDFHandler;
import org.openrdf.rio.RDFHandlerException;
import org.openrdf.rio.RDFParser;
import org.openrdf.rio.Rio;

import com.fluidops.fedx.Config;
import com.fluidops.fedx.exception.FedXException;
import com.fluidops.fedx.exception.FedXRuntimeException;
import com.fluidops.fedx.provider.EndpointProvider;
import com.fluidops.fedx.provider.NativeGraphRepositoryInformation;
import com.fluidops.fedx.provider.NativeStoreProvider;
import com.fluidops.fedx.provider.RemoteRepositoryGraphRepositoryInformation;
import com.fluidops.fedx.provider.RemoteRepositoryProvider;
import com.fluidops.fedx.provider.RepositoryInformation;
import com.fluidops.fedx.provider.SPARQLGraphRepositoryInformation;
import com.fluidops.fedx.provider.SPARQLProvider;
import com.fluidops.fedx.structures.Endpoint;
import com.fluidops.fedx.structures.Endpoint.EndpointType;



/**
 * Utility class providing various methods to create Endpoints to be used as federation members.
 * 
 * @author Andreas Schwarte
 *
 */
public class EndpointFactory {

	
	/**
	 * Construct a SPARQL endpoint using the the provided information.
	 * 
	 * @param name
	 * 			a descriptive name, e.g. http://dbpedia
	 * @param endpoint
	 * 			the URL of the SPARQL endpoint, e.g. http://dbpedia.org/sparql
	 * 
	 * @return
	 * 		an initialized {@link Endpoint} containing the repository
	 * 
	 * @throws Exception
	 */
	public static Endpoint loadSPARQLEndpoint(String name, String endpoint) throws FedXException {
		
		EndpointProvider repProvider = new SPARQLProvider();
		String id = "sparql_" + endpoint.replace("http://", "").replace("/", "_");
		return repProvider.loadEndpoint( new RepositoryInformation(id, name, endpoint, EndpointType.SparqlEndpoint));		
	}
	
	
	/**
	 * Construct a SPARQL endpoint using the the provided information and the host of the url as name.
	 * 
	 * @param endpoint
	 * 			the URL of the SPARQL endpoint, e.g. http://dbpedia.org/sparql
	 * 
	 * @return
	 * 		an initialized {@link Endpoint} containing the repository
	 * 
	 * @throws FedXException
	 */
	public static Endpoint loadSPARQLEndpoint(String endpoint) throws FedXException {
		try {
			String id = new URL(endpoint).getHost();
			if (id.equals("localhost"))
				id = id + "_" + new URL(endpoint).getPort();
			return loadSPARQLEndpoint(id, endpoint);
		} catch (MalformedURLException e) {
			throw new FedXException("Malformed URL: " + endpoint);
		}
	}
	
	
	public static Endpoint loadRemoteRepository(String repositoryServer, String repositoryName) throws FedXException {
		EndpointProvider repProvider = new RemoteRepositoryProvider();
		return repProvider.loadEndpoint( new RemoteRepositoryGraphRepositoryInformation(repositoryServer, repositoryName));		
	
	}
	
	/**
	 * Construct a NativeStore endpoint using the provided information.
	 * 
	 * @param name
	 * 			a descriptive name, e.g. http://dbpedia
	 * @param location
	 * 			the location of the data store, either absolute or relative to {@link Config#getBaseDir()}
	 * 
	 * @return
	 * 		an initialized endpoint containing the repository
	 * 
	 * @throws Exception
	 */
	public static Endpoint loadNativeEndpoint(String name, String location) throws FedXException {
		
		EndpointProvider repProvider = new NativeStoreProvider();
		String id = new File(location).getName();
		return repProvider.loadEndpoint( new RepositoryInformation(id, name, location, EndpointType.NativeStore) );
	}
	
	/**
	 * Load NativeStore from location relative to baseDir
	 * 
	 * @param name
	 * @param location
	 * @param baseDir
	 * @return
	 * @throws Exception
	 */
	public static Endpoint loadNativeEndpoint(String name, String location, File baseDir) throws FedXException {
		return loadNativeEndpoint(name, baseDir.getAbsolutePath() + "/" + location);
	}
	
	
	/**
	 * Construct a NativeStore endpoint using the provided information and the file location as name.
	 * 
	 * @param location
	 * 			the location of the data store
	 * 
	 * @return
	 * 		an initialized endpoint containing the repository
	 * 
	 * @throws Exception
	 */
	public static Endpoint loadNativeEndpoint(String location) throws FedXException {
		return loadNativeEndpoint(new File(location).getName(), location);
	}
	
	
	
	/**
	 * Utility function to load federation members from a data configuration file. A data configuration 
	 * file provides information about federation members in form of ntriples. Currently the types
	 * NativeStore and SPARQLEndpoint are supported. For details please refer to the documentation
	 * in {@link NativeGraphRepositoryInformation} and {@link SPARQLGraphRepositoryInformation}.
	 * 
	 * @param dataConfig
	 * 
	 * @return
	 * 			a list of initialized endpoints, i.e. the federation members
	 * 
	 * @throws IOException
	 * @throws Exception
	 */
	public static List<Endpoint> loadFederationMembers(File dataConfig) throws FedXException {
		
		if (!dataConfig.exists())
			throw new FedXRuntimeException("File does not exist: " + dataConfig.getAbsolutePath());
		
		Graph graph = new GraphImpl();
		RDFParser parser = Rio.createParser(RDFFormat.N3);
		RDFHandler handler = new DefaultRDFHandler(graph);
		parser.setRDFHandler(handler);
		try {
			parser.parse(new FileReader(dataConfig), "http://fluidops.org/config#");
		} catch (Exception e) {
			throw new FedXException("Unable to parse dataconfig " + dataConfig + ":" + e.getMessage());
		} 
		
		List<Endpoint> res = new ArrayList<Endpoint>();
		Iterator<Statement> iter = graph.match(null, new URIImpl("http://fluidops.org/config#store"), null);
				
		while (iter.hasNext()){
			Statement s = iter.next();
			Endpoint e = loadEndpoint(graph, s.getSubject(), s.getObject());
			res.add(e);
		}
		
		return res;
	}
	
	
	public static Endpoint loadEndpoint(Graph graph, Resource repNode, Value repType) throws FedXException {
		
		EndpointProvider repProvider;
		
		// NativeStore => Sesame native store implementation
		if (repType.equals(new LiteralImpl("NativeStore"))){
			repProvider = new NativeStoreProvider();
			return repProvider.loadEndpoint( new NativeGraphRepositoryInformation(graph, repNode) );
		} 
		
		// SPARQL Repository => SPARQLRepository 
		else if (repType.equals(new LiteralImpl("SPARQLEndpoint"))){
			repProvider =  new SPARQLProvider();	 
			return repProvider.loadEndpoint( new SPARQLGraphRepositoryInformation(graph, repNode) );
		} 
		
		// Remote Repository
		else if (repType.equals(new LiteralImpl("RemoteRepository"))){
			repProvider =  new RemoteRepositoryProvider();	 
			return repProvider.loadEndpoint( new RemoteRepositoryGraphRepositoryInformation(graph, repNode) );
		} 
		
		// other generic type
		else if (repType.equals(new LiteralImpl("Other"))) {
			
			// TODO add reflection techniques to allow for flexibility
			throw new UnsupportedOperationException("Operation not yet supported for generic type.");
			
		}
		
		else {
			throw new FedXRuntimeException("Repository type not supported: " + repType.stringValue());
		}
		
		
	}
	
	
	
	/**
	 * Construct a unique id for the provided SPARQL Endpoint, e.g
	 * 
	 * http://dbpedia.org/ => %type%_dbpedia.org
	 * 
	 * @param endpoint
	 * @param type
	 * 			the repository type, e.g. native, sparql, etc
	 * 
	 * @return
	 */
	public static String getId(String endpointID, String type) {
		String id = endpointID.replace("http://", "");
		id = id.replace("/", "_");
		return type + "_" + id;
	}
	
	
	
	protected static class DefaultRDFHandler implements RDFHandler {

		protected final Graph graph;
				
		public DefaultRDFHandler(Graph graph) {
			super();
			this.graph = graph;
		}

		@Override
		public void endRDF() throws RDFHandlerException {
			; // no-op
		}

		@Override
		public void handleComment(String comment) throws RDFHandlerException {
			; // no-op			
		}

		@Override
		public void handleNamespace(String prefix, String uri)
				throws RDFHandlerException {
			; // no-op			
		}

		@Override
		public void handleStatement(Statement st) throws RDFHandlerException {
			graph.add(st);			
		}

		@Override
		public void startRDF() throws RDFHandlerException {
			; // no-op			
		}
	}
}
