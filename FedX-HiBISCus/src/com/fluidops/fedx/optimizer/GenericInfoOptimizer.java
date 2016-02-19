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

package com.fluidops.fedx.optimizer;

import java.util.ArrayList;
import java.util.List;

import org.openrdf.query.algebra.Filter;
import org.openrdf.query.algebra.Join;
import org.openrdf.query.algebra.Service;
import org.openrdf.query.algebra.StatementPattern;
import org.openrdf.query.algebra.TupleExpr;
import org.openrdf.query.algebra.Union;
import org.openrdf.query.algebra.helpers.QueryModelVisitorBase;

import com.fluidops.fedx.algebra.NJoin;
import com.fluidops.fedx.exception.OptimizationException;
import com.fluidops.fedx.structures.QueryInfo;


/**
 * Generic optimizer
 * 
 * Tasks:
 * - Collect information (hasUnion, hasFilter, hasService)
 * - Collect all statements in a list (for source selection), do not collect SERVICE expressions
 * - Collect all Join arguments and group them in the NJoin structure for easier optimization (flatten)
 * 
 * @author Andreas Schwarte
 */
public class GenericInfoOptimizer extends QueryModelVisitorBase<OptimizationException> implements FedXOptimizer {

	protected boolean hasFilter = false;
	protected boolean hasUnion = false;
	protected boolean hasService = false;
	protected List<StatementPattern> stmts = new ArrayList<StatementPattern>();
	
	protected final QueryInfo queryInfo;
		
	public GenericInfoOptimizer(QueryInfo queryInfo) {
		super();
		this.queryInfo = queryInfo;
	}

	public boolean hasFilter() {
		return hasFilter;
	}
	
	public boolean hasUnion() {
		return hasUnion;
	}
	
	public List<StatementPattern> getStatements() {
		return stmts;
	}
	
	@Override
	public void optimize(TupleExpr tupleExpr) {
		
		try { 
			tupleExpr.visit(this);
		} catch (RuntimeException e) {
			throw e;
		} catch (Exception e) {
			throw new RuntimeException(e);
		}	
		
	}
	
	
	@Override
	public void meet(Union union) {
		hasUnion=true;
		super.meet(union);
	}
	
	@Override
	public void meet(Filter filter)  {
		hasFilter=true;
		super.meet(filter);
	}
	
	@Override
	public void meet(Service service) {
		hasService=true;
	}
	
	@Override
	public void meet(Join node) {
		
		/*
		 * Optimization task:
		 * 
		 * Collect all join arguments recursively and create the
		 * NJoin structure for easier join order optimization
		 */
				
		NJoin newJoin = OptimizerUtil.flattenJoin(node, queryInfo);
		newJoin.visitChildren(this);
		
		node.replaceWith(newJoin);
	}
	
	@Override
	public void meet(StatementPattern node) {
		stmts.add(node);
	}

	public boolean hasService()	{
		return hasService;
	}
}
