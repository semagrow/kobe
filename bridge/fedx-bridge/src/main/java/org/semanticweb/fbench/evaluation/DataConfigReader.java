package org.semanticweb.fbench.evaluation;

import java.io.File;
import java.io.FileReader;
import java.io.IOException;
import java.util.ArrayList;
import java.util.List;

import org.openrdf.model.Model;
import org.openrdf.model.Resource;
import org.openrdf.model.Statement;
import org.openrdf.model.Value;
import org.openrdf.model.impl.TreeModel;
import org.openrdf.model.impl.URIImpl;
import org.openrdf.rio.RDFFormat;
import org.openrdf.rio.RDFHandler;
import org.openrdf.rio.RDFHandlerException;
import org.openrdf.rio.RDFParser;
import org.openrdf.rio.Rio;
import org.semanticweb.fbench.report.ReportStream;

import com.fluidops.fedx.structures.Endpoint;
import com.fluidops.fedx.util.EndpointFactory;





public class DataConfigReader {
	
	
	public static List<Endpoint> loadFederationMembers(File dataConfig, ReportStream ...report) throws IOException, Exception {
			
		Model graph = new TreeModel();
		RDFParser parser = Rio.createParser(RDFFormat.N3);
		RDFHandler handler = new DefaultRDFHandler(graph);
		parser.setRDFHandler(handler);
		parser.parse(new FileReader(dataConfig), "http://fluidops.org/config#");
	
		List<Endpoint> res = new ArrayList<Endpoint>();
		for (Statement st : graph.filter(null, new URIImpl("http://fluidops.org/config#store"), null)) {
			Endpoint e = loadEndpoint(graph, st.getSubject(), st.getObject());
			res.add(e);
		}

		
		return res;
	}
	
	protected static Endpoint loadEndpoint(Model graph, Resource repNode, Value repType, ReportStream ...report) throws Exception {

		long datasetLoadStart = System.currentTimeMillis();
		
		Endpoint e = EndpointFactory.loadEndpoint(graph, repNode, repType);
		if (report.length!=0)
			report[0].reportDatasetLoadTime(e.getId(), e.getName(), e.getEndpoint(), e.getType().toString(), System.currentTimeMillis()-datasetLoadStart);
		
		return e;
	}
	
	protected static class DefaultRDFHandler implements RDFHandler {

		protected final Model graph;
				
		public DefaultRDFHandler(Model graph) {
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
