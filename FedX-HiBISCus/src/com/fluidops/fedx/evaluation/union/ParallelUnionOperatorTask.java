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

package com.fluidops.fedx.evaluation.union;

import info.aduna.iteration.CloseableIteration;

import org.openrdf.query.BindingSet;
import org.openrdf.query.QueryEvaluationException;
import org.openrdf.query.algebra.TupleExpr;

import com.fluidops.fedx.evaluation.FederationEvalStrategy;
import com.fluidops.fedx.evaluation.concurrent.ParallelExecutor;
import com.fluidops.fedx.evaluation.concurrent.ParallelTask;

/**
 * A task implementation representing a UNION operator expression to be evaluated.
 * 
 * @author Andreas Schwarte
 */
public class ParallelUnionOperatorTask implements ParallelTask<BindingSet> {

	protected final ParallelExecutor<BindingSet> unionControl;
	protected final FederationEvalStrategy strategy;
	protected final TupleExpr expr;
	protected final BindingSet bindings;
	
	public ParallelUnionOperatorTask(ParallelExecutor<BindingSet> unionControl,
			FederationEvalStrategy strategy, TupleExpr expr, BindingSet bindings) {
		super();
		this.unionControl = unionControl;
		this.strategy = strategy;
		this.expr = expr;
		this.bindings = bindings;
	}
	
	@Override
	public ParallelExecutor<BindingSet> getControl() {
		return unionControl;
	}

	@Override
	public CloseableIteration<BindingSet, QueryEvaluationException> performTask()
			throws Exception {
		return strategy.evaluate(expr, bindings);
	}
}
