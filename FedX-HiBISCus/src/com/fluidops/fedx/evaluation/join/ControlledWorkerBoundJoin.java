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

import java.util.ArrayList;
import java.util.List;

import org.apache.log4j.Logger;
import org.openrdf.query.BindingSet;
import org.openrdf.query.QueryEvaluationException;
import org.openrdf.query.algebra.StatementPattern;
import org.openrdf.query.algebra.TupleExpr;

import com.fluidops.fedx.Config;
import com.fluidops.fedx.algebra.CheckStatementPattern;
import com.fluidops.fedx.algebra.IndependentJoinGroup;
import com.fluidops.fedx.algebra.StatementTupleExpr;
import com.fluidops.fedx.evaluation.FederationEvalStrategy;
import com.fluidops.fedx.evaluation.concurrent.ControlledWorkerScheduler;
import com.fluidops.fedx.evaluation.concurrent.ParallelTask;
import com.fluidops.fedx.structures.QueryInfo;



/**
 * Execute the nested loop join in an asynchronous fashion, using grouped requests,
 * i.e. group bindings into one SPARQL request using the UNION operator.
 * 
 * The number of concurrent threads is controlled by a {@link ControlledWorkerScheduler} which
 * works according to the FIFO principle and uses worker threads.
 * 
 * This join cursor blocks until all scheduled tasks are finished, however the result iteration
 * can be accessed from different threads to allow for pipelining.
 * 
 * @author Andreas Schwarte
 * 
 */
public class ControlledWorkerBoundJoin extends ControlledWorkerJoin {

	public static Logger log = Logger.getLogger(ControlledWorkerBoundJoin.class);
	
	public ControlledWorkerBoundJoin(ControlledWorkerScheduler<BindingSet> scheduler, FederationEvalStrategy strategy,
			CloseableIteration<BindingSet, QueryEvaluationException> leftIter,
			TupleExpr rightArg, BindingSet bindings, QueryInfo queryInfo)
			throws QueryEvaluationException {
		super(scheduler, strategy, leftIter, rightArg, bindings, queryInfo);
	}

	
	@Override
	protected void handleBindings() throws Exception {
		// XXX use something else as second check, e.g. an empty interface
		if (! ((rightArg instanceof StatementPattern) || (rightArg instanceof IndependentJoinGroup) ) ) {
			log.warn("Right argument is not a StatementPattern nor a IndependentJoinGroup. Fallback on ControlledWorkerJoin implementation: " + rightArg.getClass().getCanonicalName());
			super.handleBindings();	// fallback
			return;
		}
		
		int nBindingsCfg = Config.getConfig().getBoundJoinBlockSize();	
		int totalBindings = 0;		// the total number of bindings
		TupleExpr expr = rightArg;
		
		TaskCreator taskCreator = null;
				
		// first item is always sent in a non-bound way
		if (!closed && leftIter.hasNext()) {
			BindingSet b = leftIter.next();
			totalBindings++;
			if (expr instanceof StatementTupleExpr) {
				StatementTupleExpr stmt = (StatementTupleExpr)expr;
				if (stmt.hasFreeVarsFor(b)) {
					taskCreator = new BoundJoinTaskCreator(this, strategy, stmt);
				} else {
					expr = new CheckStatementPattern(stmt);
					taskCreator = new CheckJoinTaskCreator(this, strategy, (CheckStatementPattern)expr);
				}
			} else if (expr instanceof IndependentJoinGroup) {
				taskCreator = new IndependentJoinGroupTaskCreator(this, strategy, (IndependentJoinGroup)expr);
			} else {
				throw new RuntimeException("Expr is of unexpected type: " + expr.getClass().getCanonicalName());
			}
			scheduler.schedule( new ParallelJoinTask(this, strategy, expr, b) );
		}
		
		int nBindings;	
		List<BindingSet> bindings = null;
		while (!closed && leftIter.hasNext()) {
			
			
			/*
			 * XXX idea:
			 * 
			 * make nBindings dependent on the number of intermediate results of the left argument.
			 * 
			 * If many intermediate results, increase the number of bindings. This will result in less
			 * remote SPARQL requests.
			 * 
			 */
			
			if (totalBindings>10)
				nBindings = nBindingsCfg;
			else
				nBindings = 3;

			bindings = new ArrayList<BindingSet>(nBindings);
			
			int count=0;
			while (count < nBindings && leftIter.hasNext()) {
				bindings.add(leftIter.next());
				count++;
			}
			
			totalBindings += count;		
		
			scheduler.schedule( taskCreator.getTask(bindings) );
		}
		
		scheduler.informFinish(this);
		
		log.debug("JoinStats: left iter of join #" + this.joinId + " had " + totalBindings + " results.");
				
		// wait until all tasks are executed
		synchronized (this) {
			try {
				// check to avoid deadlock
				if (scheduler.isRunning(this))
					this.wait();
			} catch (InterruptedException e) {
				;	// no-op
			}
		}	
	}

	
	
	protected interface TaskCreator {
		public ParallelTask<BindingSet> getTask(List<BindingSet> bindings);
	}

	protected class BoundJoinTaskCreator implements TaskCreator {
		protected final ControlledWorkerBoundJoin _control;
		protected final FederationEvalStrategy _strategy;
		protected final StatementTupleExpr _expr;
		public BoundJoinTaskCreator(ControlledWorkerBoundJoin control,
				FederationEvalStrategy strategy, StatementTupleExpr expr) {
			super();
			_control = control;
			_strategy = strategy;
			_expr = expr;
		}
		@Override
		public ParallelTask<BindingSet> getTask(List<BindingSet> bindings) {
			return new ParallelBoundJoinTask(_control, _strategy, _expr, bindings);
		}		
	}
	
	protected class CheckJoinTaskCreator implements TaskCreator {
		protected final ControlledWorkerBoundJoin _control;
		protected final FederationEvalStrategy _strategy;
		protected final CheckStatementPattern _expr;
		public CheckJoinTaskCreator(ControlledWorkerBoundJoin control,
				FederationEvalStrategy strategy, CheckStatementPattern expr) {
			super();
			_control = control;
			_strategy = strategy;
			_expr = expr;
		}
		@Override
		public ParallelTask<BindingSet> getTask(List<BindingSet> bindings) {
			return new ParallelCheckJoinTask(_control, _strategy, _expr, bindings);
		}		
	}
	
	protected class IndependentJoinGroupTaskCreator implements TaskCreator {
		protected final ControlledWorkerBoundJoin _control;
		protected final FederationEvalStrategy _strategy;
		protected final IndependentJoinGroup _expr;
		public IndependentJoinGroupTaskCreator(ControlledWorkerBoundJoin control,
				FederationEvalStrategy strategy, IndependentJoinGroup expr) {
			super();
			_control = control;
			_strategy = strategy;
			_expr = expr;
		}
		@Override
		public ParallelTask<BindingSet> getTask(List<BindingSet> bindings) {
			return new ParallelIndependentGroupJoinTask(_control, _strategy, _expr, bindings);
		}		
	}
	
}
