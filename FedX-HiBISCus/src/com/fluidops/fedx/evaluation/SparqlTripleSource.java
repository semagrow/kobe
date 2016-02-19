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

package com.fluidops.fedx.evaluation;

import info.aduna.iteration.CloseableIteration;
import info.aduna.iteration.EmptyIteration;
import info.aduna.iteration.ExceptionConvertingIteration;

import org.openrdf.model.Resource;
import org.openrdf.model.Statement;
import org.openrdf.model.URI;
import org.openrdf.model.Value;
import org.openrdf.query.BindingSet;
import org.openrdf.query.BooleanQuery;
import org.openrdf.query.MalformedQueryException;
import org.openrdf.query.QueryEvaluationException;
import org.openrdf.query.QueryLanguage;
import org.openrdf.query.TupleQuery;
import org.openrdf.query.algebra.StatementPattern;
import org.openrdf.query.algebra.TupleExpr;
import org.openrdf.repository.RepositoryConnection;
import org.openrdf.repository.RepositoryException;
import org.openrdf.repository.RepositoryResult;

import com.fluidops.fedx.FederationManager;
import com.fluidops.fedx.algebra.FilterValueExpr;
import com.fluidops.fedx.evaluation.iterator.FilteringInsertBindingsIteration;
import com.fluidops.fedx.evaluation.iterator.FilteringIteration;
import com.fluidops.fedx.evaluation.iterator.InsertBindingsIteration;
import com.fluidops.fedx.exception.ExceptionUtil;
import com.fluidops.fedx.structures.Endpoint;
import com.fluidops.fedx.util.QueryStringUtil;



/**
 * A triple source to be used for (remote) SPARQL endpoints. 
 * 
 * @author Andreas Schwarte
 *
 */
public class SparqlTripleSource extends TripleSourceBase implements TripleSource {

	
	SparqlTripleSource(Endpoint endpoint) {
		super(FederationManager.getMonitoringService(), endpoint);
	}
	
	@Override
	public CloseableIteration<BindingSet, QueryEvaluationException> getStatements(
			String preparedQuery, RepositoryConnection conn, BindingSet bindings, FilterValueExpr filterExpr)
			throws RepositoryException, MalformedQueryException,
			QueryEvaluationException {
		
		TupleQuery query = conn.prepareTupleQuery(QueryLanguage.SPARQL, preparedQuery, null);
		
		try {			
			
			// evaluate the query
			monitorRemoteRequest();
			CloseableIteration<BindingSet, QueryEvaluationException> res = query.evaluate();
			
			// apply filter and/or insert original bindings
			if (filterExpr!=null) {
				if (bindings.size()>0) 
					res = new FilteringInsertBindingsIteration(filterExpr, bindings, res);
				else
					res = new FilteringIteration(filterExpr, res);
				if (!res.hasNext())
					return new EmptyIteration<BindingSet, QueryEvaluationException>();
			} else if (bindings.size()>0) {
				res = new InsertBindingsIteration(res, bindings);
			}
	
			return res;
			
		} catch (QueryEvaluationException ex) {
			throw ExceptionUtil.traceExceptionSourceAndRepair(conn, ex, "Subquery: " + preparedQuery);			
		}
	}

	@Override
	public CloseableIteration<BindingSet, QueryEvaluationException> getStatements(
			StatementPattern stmt, RepositoryConnection conn,
			BindingSet bindings, FilterValueExpr filterExpr)
			throws RepositoryException, MalformedQueryException,
			QueryEvaluationException  {
		
		throw new RuntimeException("NOT YET IMPLEMENTED.");
	}

	@Override
	public boolean hasStatements(StatementPattern stmt, RepositoryConnection conn,
			BindingSet bindings) throws RepositoryException,
			MalformedQueryException, QueryEvaluationException {

		/* remote boolean query */
		String queryString = QueryStringUtil.askQueryString(stmt, bindings);
		BooleanQuery query = conn.prepareBooleanQuery(QueryLanguage.SPARQL, queryString, null);
		
		try {
			monitorRemoteRequest();
			boolean hasStatements = query.evaluate();
			return hasStatements;
		} catch (QueryEvaluationException ex) {
			throw ExceptionUtil.traceExceptionSourceAndRepair(conn, ex, "Subquery: " + queryString);			
		}

		
	}

	@Override
	public boolean usePreparedQuery() {
		return true;
	}

	@Override
	public CloseableIteration<BindingSet, QueryEvaluationException> getStatements(
			TupleExpr preparedQuery, RepositoryConnection conn,
			BindingSet bindings, FilterValueExpr filterExpr)
			throws RepositoryException, MalformedQueryException,
			QueryEvaluationException {
		
		throw new RuntimeException("NOT YET IMPLEMENTED.");
	}

	@Override
	public CloseableIteration<Statement, QueryEvaluationException> getStatements(
			RepositoryConnection conn, Resource subj, URI pred, Value obj,
			Resource... contexts) throws RepositoryException,
			MalformedQueryException, QueryEvaluationException
	{
		
		// TODO add handling for contexts
		monitorRemoteRequest();
		RepositoryResult<Statement> repoResult = conn.getStatements(subj, pred, obj, true);
		
		return new ExceptionConvertingIteration<Statement, QueryEvaluationException>(repoResult) {
			@Override
			protected QueryEvaluationException convert(Exception arg0) {
				return new QueryEvaluationException(arg0);
			}
		};		
	}

}
