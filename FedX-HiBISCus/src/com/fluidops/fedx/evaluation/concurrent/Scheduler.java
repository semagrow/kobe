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

package com.fluidops.fedx.evaluation.concurrent;

import info.aduna.iteration.CloseableIteration;

import org.openrdf.query.QueryEvaluationException;

/**
 * Interface for any scheduler. 
 * 
 * @author Andreas Schwarte
 *
 * @see ControlledWorkerScheduler
 */
public interface Scheduler<T> {

	/**
	 * Schedule the provided task.
	 * 
	 * @param task
	 */
	public void schedule(ParallelTask<T> task);
	
	/**
	 * Callback to handle the result.
	 * 
	 * @param res
	 */
	public void handleResult(CloseableIteration<T, QueryEvaluationException> res);
	
	/**
	 * Inform the scheduler that a certain task is done.
	 * 
	 */
	public void done();
	
	/**
	 * Toss an exception to the scheduler.
	 * 
	 * @param e
	 */
	public void toss(Exception e);
	
	/**
	 * Abort the execution of running and queued tasks.
	 * 
	 */
	public void abort();
	
	/**
	 * Inform the scheduler that no more tasks will be scheduled.
	 */
	public void informFinish();
	
	/**
	 * Determine if the scheduler has unfisnihed tasks.
	 * 
	 * @return
	 */
	public boolean isRunning();
	
}
