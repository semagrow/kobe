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
import info.aduna.iteration.EmptyIteration;
import info.aduna.iteration.LookAheadIteration;

import org.apache.log4j.Logger;
import org.openrdf.query.QueryEvaluationException;

import com.fluidops.fedx.evaluation.concurrent.ParallelExecutor;
import com.fluidops.fedx.evaluation.iterator.QueueCursor;


/**
 * Base class for any parallel union executor.
 * 
 * Note that this class extends {@link LookAheadIteration} and thus any implementation of this 
 * class is applicable for pipelining when used in a different thread (access to shared
 * variables is synchronized).
 * 
 * @author Andreas Schwarte
 *
 */
public abstract class UnionExecutorBase<T> extends LookAheadIteration<T, QueryEvaluationException> implements ParallelExecutor<T> {

	public static Logger log = Logger.getLogger(UnionExecutorBase.class);
	protected static int NEXT_UNION_ID = 1;
	
	/* Constants */
	protected final int unionId;							// the union id
	
	/* Variables */
	protected volatile boolean closed;
	protected boolean finished = true;
	
	protected QueueCursor<CloseableIteration<T, QueryEvaluationException>> result = new QueueCursor<CloseableIteration<T, QueryEvaluationException>>(1024);
	protected CloseableIteration<T, QueryEvaluationException> rightIter;
	
	
	public UnionExecutorBase() {
		this.unionId = NEXT_UNION_ID++;
	}
	

	@Override
	public final void run() {

		try {
			union();
		} catch (Exception e) {
			toss(e);
		} finally {
			finished=true;
			result.done();
		}
		
	}
	

	/**
	 * 
	 * Note: this method must block until the union is executed completely. Otherwise
	 * the result queue is marked as committed while this isn't the case. The blocking
	 * behavior in general is no problem: If you need concurrent access to the result
	 * (i.e. pipelining) just run the union in a separate thread. Access to the result
	 * iteration is synchronized.
	 * 
	 * @throws Exception
	 */
	protected abstract void union() throws Exception;
	
	
	@Override
	public void addResult(CloseableIteration<T, QueryEvaluationException> res)  {
		/* optimization: avoid adding empty results */
		if (res instanceof EmptyIteration<?,?>)
			return;
		
		try {
			result.put(res);
		} catch (InterruptedException e) {
			throw new RuntimeException("Error adding element to result queue", e);
		}

	}
		
	@Override
	public void done() {
		;	// no-op
	}
	
	@Override
	public void toss(Exception e) {
		log.warn("Error executing union operator: " + e.getMessage());
		result.toss(e);
	}
	
	
	@Override
	public T getNextElement() throws QueryEvaluationException	{
		// TODO check if we need to protect rightQueue from synchronized access
		// wasn't done in the original implementation either
		// if we see any weird behavior check here !!

		while (rightIter != null || result.hasNext()) {
			if (rightIter == null) {
				rightIter = result.next();
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
		
		if (rightIter != null) {
			rightIter.close();
			rightIter = null;
		}

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
}
