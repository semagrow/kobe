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

package com.fluidops.fedx.algebra;

import info.aduna.iteration.CloseableIteration;
import info.aduna.iteration.EmptyIteration;

import org.openrdf.query.BindingSet;
import org.openrdf.query.MalformedQueryException;
import org.openrdf.query.QueryEvaluationException;
import org.openrdf.query.algebra.StatementPattern;
import org.openrdf.repository.RepositoryConnection;
import org.openrdf.repository.RepositoryException;

import com.fluidops.fedx.EndpointManager;
import com.fluidops.fedx.evaluation.TripleSource;
import com.fluidops.fedx.evaluation.iterator.SingleBindingSetIteration;
import com.fluidops.fedx.exception.IllegalQueryException;
import com.fluidops.fedx.structures.Endpoint;
import com.fluidops.fedx.structures.QueryInfo;
import com.fluidops.fedx.util.QueryStringUtil;



/**
 * Represents a StatementPattern that can only produce results at a single endpoint, the owner.
 * 
 * @author Andreas Schwarte
 */
public class ExclusiveStatement extends FedXStatementPattern implements StatementTupleExpr, FilterTuple {

	protected FilterValueExpr filterExpr = null;
	
	public ExclusiveStatement(StatementPattern node, StatementSource owner, QueryInfo queryInfo) {
		super(node, queryInfo);
		statementSources.add(owner);
	}	

	public StatementSource getOwner() {
		return getStatementSources().get(0);
	}	

	@Override
	public CloseableIteration<BindingSet, QueryEvaluationException> evaluate(
			BindingSet bindings) throws QueryEvaluationException {
		
		try {
			
			Endpoint ownedEndpoint = EndpointManager.getEndpointManager().getEndpoint(getOwner().getEndpointID());
			RepositoryConnection ownedConnection = ownedEndpoint.getConn();
			TripleSource t = ownedEndpoint.getTripleSource();
			
			/*
			 * Implementation note: for some endpoint types it is much more efficient to use prepared queries
			 * as there might be some overhead (obsolete optimization) in the native implementation. This
			 * is for instance the case for SPARQL connections. In contrast for NativeRepositories it is
			 * much more efficient to use getStatements(subj, pred, obj) instead of evaluating a prepared query.
			 */			
		
			if (t.usePreparedQuery()) {
				
				Boolean isEvaluated = false;	// is filter evaluated
				String preparedQuery;
				try {
					preparedQuery = QueryStringUtil.selectQueryString(this, bindings, filterExpr, isEvaluated);
				} catch (IllegalQueryException e1) {
					// TODO there might be an issue with filters being evaluated => investigate
					/* all vars are bound, this must be handled as a check query, can occur in joins */
					if (t.hasStatements(this, ownedConnection, bindings))
						return new SingleBindingSetIteration(bindings);
					return new EmptyIteration<BindingSet, QueryEvaluationException>();
				}
								
				return t.getStatements(preparedQuery, ownedConnection, bindings, (isEvaluated ? null : filterExpr) );
				
			} else {
				return t.getStatements(this, ownedConnection, bindings, filterExpr);
			}
				
		} catch (RepositoryException e) {
			throw new QueryEvaluationException(e);
		} catch (MalformedQueryException e) {
			throw new QueryEvaluationException(e);
		}
	}
}
