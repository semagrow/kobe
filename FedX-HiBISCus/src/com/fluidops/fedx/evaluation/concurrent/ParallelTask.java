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
 * Interface for any parallel task that can be performed in Scheduler implementations.
 * 
 * @author Andreas Schwarte
 *
 */
public interface ParallelTask<T> {

	public CloseableIteration<T, QueryEvaluationException> performTask() throws Exception;
	
	/**
	 * return the controlling instance, e.g. in most cases the instance of a thread. Shared variables
	 * are used to inform the thread about new events.
	 * 
	 * @return
	 */
	public ParallelExecutor<T> getControl();
}
