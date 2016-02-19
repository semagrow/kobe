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

import com.fluidops.fedx.evaluation.concurrent.ControlledWorkerScheduler;
import com.fluidops.fedx.structures.QueryInfo;

/**
 * Execution of union tasks with {@link ControlledWorkerScheduler}. Tasks can be added
 * using the provided functions. Note that the union operation is to be executed
 * with the {@link #run()} method (also threaded execution is possible). Results are
 * then contained in this iteration.
 *
 * @author Andreas Schwarte
 *
 */
public class ControlledWorkerUnion<T> extends WorkerUnionBase<T> {

	public static int waitingCount = 0;
	public static int awakeCount = 0;
	
	protected final ControlledWorkerScheduler<T> scheduler;
	
	public ControlledWorkerUnion(ControlledWorkerScheduler<T> scheduler, QueryInfo queryInfo) {
		super(queryInfo);
		this.scheduler = scheduler;
	}
			
	
	@Override
	protected void union() throws Exception {
		
		// schedule all tasks and inform about finish
		scheduler.scheduleAll(tasks, this);
		
		// wait until all tasks are executed
		synchronized (this) {
			try {
				// check to avoid deadlock
				// FIXME there is a possible deadlock situation!!! tiny, but it is there
				// maybe use a monitor that checks for sleeping tasks?
//				boolean checked = false;
				while (scheduler.isRunning(this)) {
//					if (!checked)
//						System.out.println("Waiting Threads Count " + (++waitingCount));
					this.wait();
//					if (!checked)
//						System.out.println("Awake Threads Count " + (++awakeCount));
//					checked = true;
				}
			} catch (InterruptedException e) {
				//log.warn("Union thread was interrupted.");
			}
		}
	}
}
