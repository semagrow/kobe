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

package com.fluidops.fedx.provider;

import java.util.Iterator;

import org.openrdf.model.Graph;
import org.openrdf.model.Resource;
import org.openrdf.model.Statement;
import org.openrdf.model.impl.URIImpl;

import com.fluidops.fedx.structures.Endpoint.EndpointType;


/**
 * Graph information for Sesame SPARQLRepository initialization.
 * 
 * Format:
 * 
 * <code>
 * <%name%> fluid:store "SPARQLEndpoint";
 * fluid:SPARQLEndpoint "%location%"
 * 
 * <http://DBpedia> fluid:store "SPARQLEndpoint";
 * fluid:SPARQLEndpoint "http://dbpedia.org/sparql".
 * 
 * <http://NYtimes> fluid:store "SPARQLEndpoint";
 * fluid:SPARQLEndpoint "http://api.talis.com/stores/nytimes/services/sparql".
 * </code>
 * 
 * Note: the id is constructed from the name: http://dbpedia.org/ => sparql_dbpedia.org
 * 
 * 
 * @author Andreas Schwarte
 *
 */
public class SPARQLGraphRepositoryInformation extends RepositoryInformation {

	public SPARQLGraphRepositoryInformation(Graph graph, Resource repNode) {
		super(EndpointType.SparqlEndpoint);
		initialize(graph, repNode);
	}

	protected void initialize(Graph graph, Resource repNode) {
		
		// name: the node's value
		setProperty("name", repNode.stringValue());
				
		// location
		Iterator<Statement> iter = graph.match(repNode, new URIImpl("http://fluidops.org/config#SPARQLEndpoint"), null);
		String repoLocation = iter.next().getObject().stringValue();
		setProperty("location", repoLocation);
		
		// id: the name of the location
		String id = repNode.stringValue().replace("http://", "");
		id = "sparql_" + id.replace("/", "_");
		setProperty("id", id);
	}
}
