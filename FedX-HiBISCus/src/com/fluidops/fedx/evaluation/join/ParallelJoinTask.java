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

package com.fluidops.fedx.evaluation.join;

import info.aduna.iteration.CloseableIteration;

import org.openrdf.query.BindingSet;
import org.openrdf.query.QueryEvaluationException;
import org.openrdf.query.algebra.TupleExpr;

import com.fluidops.fedx.evaluation.FederationEvalStrategy;
import com.fluidops.fedx.evaluation.concurrent.ParallelExecutor;
import com.fluidops.fedx.evaluation.concurrent.ParallelTask;

/**
 * A task implementation representing a join, i.e. the provided expression is evaluated 
 * with the given bindings.
 * 
 * @author Andreas Schwarte
 */
public class ParallelJoinTask implements ParallelTask<BindingSet> {
	
	protected final FederationEvalStrategy strategy;
	protected final TupleExpr expr;
	protected final BindingSet bindings;
	protected final ParallelExecutor<BindingSet> joinControl;
	
	public ParallelJoinTask(ParallelExecutor<BindingSet> joinControl, FederationEvalStrategy strategy, TupleExpr expr, BindingSet bindings) {
		this.strategy = strategy;
		this.expr = expr;
		this.bindings = bindings;
		this.joinControl = joinControl;
	}

	@Override
	public CloseableIteration<BindingSet, QueryEvaluationException> performTask() throws Exception {
		return strategy.evaluate(expr, bindings);		
	}

	@Override
	public ParallelExecutor<BindingSet> getControl() {
		return joinControl;
	}
}
