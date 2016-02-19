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

package com.fluidops.fedx.structures;

import org.openrdf.model.Resource;
import org.openrdf.model.URI;
import org.openrdf.model.Value;

import com.fluidops.fedx.util.QueryStringUtil;



/**
 * Structure to maintain query information during evaluation, is attached to algebra nodes. 
 * Each instance is uniquely attached to the query.
 * 
 * The queryId can be used to abort tasks belonging to a particular evaluation.
 * 
 * @author Andreas Schwarte
 *
 */
public class QueryInfo {

	protected static int NEXT_QUERY_ID = 1;		// static id count
	
	private final int queryID;
	private final String query;
	private final QueryType queryType;
	
	public QueryInfo(String query, QueryType queryType) {
		super();
		synchronized (QueryInfo.class) {
			this.queryID = NEXT_QUERY_ID++;
		}
		this.query = query;
		this.queryType = queryType;
	}
	public QueryInfo(Resource subj, URI pred, Value obj) {
		this(QueryStringUtil.toString(subj, pred, obj), QueryType.GET_STATEMENTS);
	}

	public int getQueryID() {
		return queryID;
	}

	public String getQuery() {
		return query;
	}	
	
	public QueryType getQueryType() {
		return queryType;
	}

	@Override
	public int hashCode()
	{
		final int prime = 31;
		int result = 1;
		result = prime * result + queryID;
		return result;
	}

	@Override
	public boolean equals(Object obj)
	{
		if (this == obj)
			return true;
		if (obj == null)
			return false;
		if (getClass() != obj.getClass())
			return false;
		QueryInfo other = (QueryInfo) obj;
		if (queryID != other.queryID)
			return false;
		return true;
	}	
	
	
}
