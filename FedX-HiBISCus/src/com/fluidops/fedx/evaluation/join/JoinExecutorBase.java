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
import info.aduna.iteration.EmptyIteration;
import info.aduna.iteration.LookAheadIteration;

import org.apache.log4j.Logger;
import org.openrdf.query.BindingSet;
import org.openrdf.query.QueryEvaluationException;
import org.openrdf.query.algebra.TupleExpr;

import com.fluidops.fedx.evaluation.FederationEvalStrategy;
import com.fluidops.fedx.evaluation.concurrent.ParallelExecutor;
import com.fluidops.fedx.evaluation.iterator.QueueCursor;
import com.fluidops.fedx.structures.QueryInfo;


/**
 * Base class for any join parallel join executor. 
 * 
 * Note that this class extends {@link LookAheadIteration} and thus any implementation of this 
 * class is applicable for pipelining when used in a different thread (access to shared
 * variables is synchronized).
 * 
 * @author Andreas Schwarte
 */
public abstract class JoinExecutorBase<T> extends LookAheadIteration<T, QueryEvaluationException> implements ParallelExecutor<T> {

	public static Logger log = Logger.getLogger(JoinExecutorBase.class);
	
	protected static int NEXT_JOIN_ID = 1;
	
	/* Constants */
	protected final FederationEvalStrategy strategy;		// the evaluation strategy
	protected final TupleExpr rightArg;						// the right argument for the join
	protected final BindingSet bindings;					// the bindings
	protected final int joinId;								// the join id
	protected final QueryInfo queryInfo;
	
	/* Variables */
	protected volatile Thread evaluationThread;
	protected CloseableIteration<T, QueryEvaluationException> leftIter;
	protected CloseableIteration<T, QueryEvaluationException> rightIter;
	protected volatile boolean closed;
	protected boolean finished = false;
	
	protected QueueCursor<CloseableIteration<T, QueryEvaluationException>> rightQueue = new QueueCursor<CloseableIteration<T, QueryEvaluationException>>(1024);

	
	public JoinExecutorBase(FederationEvalStrategy strategy, CloseableIteration<T, QueryEvaluationException> leftIter, TupleExpr rightArg,
			BindingSet bindings, QueryInfo queryInfo) throws QueryEvaluationException	{
		this.strategy = strategy;
		this.leftIter = leftIter;
		this.rightArg = rightArg;
		this.bindings = bindings;
		this.joinId = NEXT_JOIN_ID++;
		this.queryInfo = queryInfo;
	}
	

	@Override
	public final void run() {
		evaluationThread = Thread.currentThread();
		

		if (log.isTraceEnabled())
			log.trace("Performing join #" + joinId);
		
		try {
			handleBindings();
		} catch (Exception e) {
			toss(e);
		} finally {
			finished=true;
			evaluationThread = null;
			rightQueue.done();
		}
				
		if (log.isTraceEnabled())
			log.trace("Join #" + joinId + " is finished.");
	}
	
	/**
	 * Implementations must implement this method to handle bindings.
	 * 
	 * Use the following as a template
	 * <code>
	 * while (!closed && leftIter.hasNext()) {
	 * 		// your code
	 * }
	 * </code>
	 * 
	 * and add results to rightQueue. Note that addResult() is implemented synchronized
	 * and thus thread safe. In case you can guarantee sequential access, it is also
	 * possible to directly access rightQueue
	 * 
	 */
	protected abstract void handleBindings() throws Exception;
	
	
	@Override
	public void addResult(CloseableIteration<T, QueryEvaluationException> res)  {
		/* optimization: avoid adding empty results */
		if (res instanceof EmptyIteration<?,?>)
			return;
		
		try {
			rightQueue.put(res);
		} catch (InterruptedException e) {
			throw new RuntimeException("Error adding element to right queue", e);
		}
	}
		
	@Override
	public void done() {
		;	// no-op
	}
	
	@Override
	public void toss(Exception e) {
		rightQueue.toss(e);
	}
	
	
	@Override
	public T getNextElement() throws QueryEvaluationException	{
		// TODO check if we need to protect rightQueue from synchronized access
		// wasn't done in the original implementation either
		// if we see any weird behavior check here !!

		while (rightIter != null || rightQueue.hasNext()) {
			if (rightIter == null) {
				rightIter = rightQueue.next();
			}
			if (rightIter.hasNext()) {
				return rightIter.next();
			}
			else {
				rightIter.close();
				rightIter = null;
			}
		}

		return null;
	}

	
	@Override
	public void handleClose() throws QueryEvaluationException {
		closed = true;
		if (evaluationThread != null) {
			evaluationThread.interrupt();
		}
		
		if (rightIter != null) {
			rightIter.close();
			rightIter = null;
		}

		leftIter.close();
	}
	
	/**
	 * Return true if this executor is finished or aborted
	 * 
	 * @return
	 */
	public boolean isFinished() {
		synchronized (this) {
			return finished;
		}
	}
	
	/**
	 * Retrieve information about this join, joinId and queryId
	 * 
	 * @return
	 */
	public String getId() {
		return "ID=(id:" + joinId + "; query:" + getQueryId() + ")";
	}
	
	@Override
	public int getQueryId() {
		if (queryInfo!=null)
			return queryInfo.getQueryID();
		return -1;
	}
}
