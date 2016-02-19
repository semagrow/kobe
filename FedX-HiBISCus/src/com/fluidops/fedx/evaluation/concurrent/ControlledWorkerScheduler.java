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

import java.lang.reflect.UndeclaredThrowableException;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.HashSet;
import java.util.LinkedList;
import java.util.List;
import java.util.Map;
import java.util.Set;

import org.apache.log4j.Logger;
import org.openrdf.query.QueryEvaluationException;

import com.fluidops.fedx.Config;
import com.fluidops.fedx.evaluation.join.ControlledWorkerBoundJoin;
import com.fluidops.fedx.evaluation.join.ControlledWorkerJoin;
import com.fluidops.fedx.evaluation.union.ControlledWorkerUnion;



/**
 * ControlledWorkerScheduler is a task scheduler that uses a FIFO queue for managing
 * its process. Each instance has a pool with a fixed number of worker threads. Once
 * notified a worker picks the next task from the queue and executes it. The results
 * is then returned to the controlling instance retrieved from the task.
 * 
 * Note that access to shared instance variables is synchronized.
 * 
 * Usage:
 * 
 * <code>
 *  ControlledWorkerScheduler scheduler = ...;	// get your instance => singleton?
 * 
 *  while (available) {
 *  	scheduler.schedule(task);
 *  }
 * 
 *  scheduler.informFinish(this);
 * 
 *  // wait until all tasks are executed
 *  synchronized (this) {
 * 		try {
 * 			// check to avoid deadlock
 * 			if (scheduler.isRunning(this))
 * 				this.wait();
 * 			} catch (InterruptedException e) {
 * 				;	// no-op
 * 			}
 * 		}	
 *  }
 * </code>
 * 
 * @author Andreas Schwarte
 * 
 * @see ControlledWorkerUnion
 * @see ControlledWorkerJoin
 * @see ControlledWorkerBoundJoin
 */
public class ControlledWorkerScheduler<T> implements Scheduler<T> {

	protected static final Logger log = Logger.getLogger(ControlledWorkerScheduler.class);
	protected static int NEXT_ID = 1;
	
	protected LinkedList<ParallelTask<T>> taskQueue = new LinkedList<ParallelTask<T>>();
	protected List<WorkerThread> workers = new ArrayList<WorkerThread>();
	protected Map<ParallelExecutor<T>, ControlStatus> controlMap = new HashMap<ParallelExecutor<T>, ControlStatus>();
	protected int nWorkers;
	protected int idleWorkers;
	protected IdleWorkersMonitor idleWorkersMonitor = null;
	protected String name;
	
		
	/**
	 * Construct a new instance with 20 workers.
	 */
	public ControlledWorkerScheduler() {
		this(20);
	}
	
	/**
	 * Construct a new instance with the specified number of workers.
	 * 
	 * @param nWorkers
	 * 			the number of worker threads
	 */
	public ControlledWorkerScheduler(int nWorkers) {
		this(nWorkers, "CW #" + (NEXT_ID++));
		
	}
	
	/**
	 * Construct a new instance with the specified number of workers and the
	 * given name.
	 * 
	 * @param nWorkers
	 * @param name
	 */
	public ControlledWorkerScheduler(int nWorkers, String name) {
		this.nWorkers = nWorkers;
		initWorkerThreads();
		this.name = name;
	}
	
	
	/**
	 * Schedule the specified parallel task.
	 * 	
	 * @param task
	 * 			the task to schedule
	 */
	@Override
	public void schedule(ParallelTask<T> task) {
		
		synchronized (taskQueue) {
			ControlStatus status = controlMap.get(task.getControl());
			if (status==null) {
				status = new ControlStatus(0, false);
				controlMap.put(task.getControl(), status);
			}
			status.waiting++;
			taskQueue.addLast(task);
			taskQueue.notify();
		}
		
	}	
	
	
	/**
	 * Schedule the given tasks and inform about finish using the same lock, i.e.
	 * all tasks are scheduled one after the other.
	 * @param tasks
	 * @param control
	 */
	public void scheduleAll(List<ParallelTask<T>> tasks, ParallelExecutor<T> control) {
		
		synchronized (taskQueue) {
			ControlStatus status = controlMap.get(control);
			if (status==null) {
				status = new ControlStatus(0, false);
				controlMap.put(control, status);
			}
			for (ParallelTask<T> task : tasks) {
				status.waiting++;
				taskQueue.addLast(task);
				if (idleWorkers>0)
					taskQueue.notify();
			}
			status.done = true;
		}
		
	}
	
	
	protected void initWorkerThreads() {
		for (int i=0; i<nWorkers; i++) {
			WorkerThread t = new WorkerThread();
			workers.add(t); 
			t.start();
		}

		if (Config.getConfig().isDebugWorkerScheduler()) {
			log.info("Registering IdleStatusMonitor. Total number of workers: " + nWorkers);
			idleWorkersMonitor = new IdleWorkersMonitor();
			idleWorkersMonitor.start();
		}
	}
	
	@Override
	public void abort() {
		log.info("Aborting workers of " + name + ".");
		if (idleWorkersMonitor!=null) {
			idleWorkersMonitor.close();
			idleWorkersMonitor.interrupt();
		}
		synchronized (taskQueue) {
			taskQueue.clear();
			taskQueue.notifyAll();
			for (WorkerThread t : workers) {
				t.abort();
				t.interrupt();
			}
		}						
	}
	
	/**
	 * Abort all task belonging to control
	 * 
	 * @param control
	 */
	public void abort(ParallelExecutor<T> control) {
		log.debug("Aborting tasks for executor " + control + " due to previous error.");
		// remove all tasks belonging to this control from the queue
		synchronized (taskQueue) {
			if (control.isFinished()) {
				log.debug("Join executor is already finished or aborted, no tasks are aborted.");
				return;
			}
			LinkedList<ParallelTask<T>> copy = (LinkedList<ParallelTask<T>>)taskQueue.clone();
			taskQueue.clear();
			for (ParallelTask<T> task : copy) {
				if (!task.getControl().equals(control))
					taskQueue.add(task);
			}
			
			ControlStatus st = controlMap.get(control);
			if (st!=null) {
				st.waiting=0;
				st.done=true;
				controlMap.remove(control);
			}
			
			// interrupt workers working on this task and start new threads
			int idlecount=0;
			ArrayList<WorkerThread> new_workers = new ArrayList<WorkerThread>(workers.size());
			for (WorkerThread t : workers) {
				if (t.isHandlingTask(control)) {
					t.abort();
					t.interrupt();
					WorkerThread tmp = new WorkerThread();
					new_workers.add(tmp); 
					tmp.start();
				} else {
					if (t.isIdle())
						idlecount++;
					new_workers.add(t);
				}				
			}
			workers = new_workers;
			idleWorkers = idlecount;
		}
		
	}
	
	/**
	 * Abort all task belonging to query with provided id (i.e. from QueryInfo)
	 * 
	 * @param queryId
	 * 			the valid queryId or -1 if not available
	 */
	public void abort(int queryId) {
		log.debug("Aborting tasks for query with id " + queryId + ".");
		if (queryId<0)
			return;
		
		// remove all tasks belonging to this queryId
		synchronized (taskQueue) {
			
			Set<ParallelExecutor<T>> queryControls = new HashSet<ParallelExecutor<T>>();
			LinkedList<ParallelTask<T>> copy = (LinkedList<ParallelTask<T>>)taskQueue.clone();
			taskQueue.clear();
			for (ParallelTask<T> task : copy) {
				if (task.getControl().getQueryId()!=queryId)
					taskQueue.add(task);
			}
			
			for (ParallelExecutor<T> control : queryControls) {
				ControlStatus st = controlMap.get(control);
				if (st!=null) {
					st.waiting=0;
					st.done=true;
					controlMap.remove(control);
				}
			}
			
			// interrupt workers working on this task and start new threads
			int idlecount=0;
			ArrayList<WorkerThread> new_workers = new ArrayList<WorkerThread>(workers.size());
			for (WorkerThread t : workers) {
				if (t.isHandlingTask(queryId)) {
					t.abort();
					t.interrupt();
					WorkerThread tmp = new WorkerThread();
					new_workers.add(tmp); 
					tmp.start();
				} else {
					if (t.isIdle())
						idlecount++;
					new_workers.add(t);
				}				
			}
			workers = new_workers;
			idleWorkers = idlecount;
		}
		
	}
	

	@Override
	public void done() {
		/* not needed here, implementations call informFinish(control) to notify done status */		
	}


	@Override
	public void handleResult(CloseableIteration<T, QueryEvaluationException> res) {
		/* not needed here since the result is passed directly to the control instance */		
		throw new RuntimeException("Unsupported Operation for this scheduler.");
	}

	@Override
	public void informFinish() {
		throw new RuntimeException("Unsupported Operation for this scheduler!");	
	}
	
	/**
	 * Inform this scheduler that the specified control instance will no longer
	 * submit tasks.
	 * 
	 * @param control
	 */
	public void informFinish(ParallelExecutor<T> control) {
		
		synchronized (taskQueue) {
			ControlStatus st = controlMap.get(control);
			if (st!=null)
				st.done=true;
		}
	}
	

	@Override
	public boolean isRunning() {
		/* Note: this scheduler can only determine runtime for a given control instance! */
		throw new RuntimeException("Unsupported Operation for this scheduler.");
	}

	
	/**
	 * Determine if there are still task running or queued for the specified control.
	 * 
	 * @param control
	 * 
	 * @return
	 * 		true, if there are unfinished tasks, false otherwise
	 */
	public boolean isRunning(ParallelExecutor<T> control) {
		synchronized (taskQueue) {
			ControlStatus st = controlMap.get(control);
			if (st==null)
				return false;
			return st.waiting>0;
		}
	}
	
	
	@Override
	public void toss(Exception e) {
		/* not needed here: exceptions are directly tossed to the controlling instance */		
		throw new RuntimeException("Unsupported Operation for this scheduler.");
	}

	
	
	/**
	 * Worker implementation that performs the tasks available in the queue.
	 * 
	 * @author Andreas Schwarte
	 */
	protected class WorkerThread extends Thread {
		
		protected boolean aborted = false;
		protected boolean inTask = false;
		protected boolean idle = false;
		protected ParallelTask<T> task;
		
		public WorkerThread() {			
		}
				
		@Override
		public void run() {
			
			task = null;
			
			synchronized (taskQueue) {
				if (!taskQueue.isEmpty()) 
					task = taskQueue.removeFirst();
			}
			
			while (!isAborted()) {

				if (task!=null) {
					
					ParallelExecutor<T> taskControl = task.getControl();
					
					try {
						inTask = true;
						CloseableIteration<T, QueryEvaluationException> res = task.performTask();
						inTask = false;
						taskControl.addResult(res);						
						
					} catch (UndeclaredThrowableException e) {
						// happens if a thread gets interrupted
						if (e.getCause()!=null && (e.getCause() instanceof InterruptedException)) {
							if (isAborted())
								return;
						}
						throw e;
					} catch (Exception e) {
						if (isAborted()) {
							return;
						}
						log.warn("Exception encountered while evaluating task (" + e.getClass().getSimpleName() + "): " + e.getMessage());
						taskControl.toss(e);
					}
					
					if (isAborted()) {
						log.trace("Workerthread has been aborted.");
						return;
					}
					taskControl.done();		// in most cases this is a no-op
					
					synchronized (taskQueue) {
						
						ControlStatus st = controlMap.get(taskControl);
						if (st!=null) {
							st.waiting--;
							
							// check if finished: notify control as it might be sleeping
							if (st.done && st.waiting<=0) {
								synchronized (taskControl) {
									taskControl.notify();
								}
								controlMap.remove(taskControl);
							}
						}
						
						if (!taskQueue.isEmpty()) {
							task = taskQueue.removeFirst();
							continue;
						}
					}
				}
				
				// wait on signal on the queue => push
    			synchronized (taskQueue) {
    				try {
    					if (!taskQueue.isEmpty()) {
    						task = taskQueue.removeFirst();
       						continue;	// check to guarantee correctness
    					}
    					idleWorkers++;
    					idle=true;
						taskQueue.wait();
						idle=false;
						idleWorkers--;
						if (!taskQueue.isEmpty())
							task = taskQueue.removeFirst();
						else
							task = null;
    				} catch (InterruptedException e) {
						
    					if (aborted) {
    						log.trace("Workerthread has been aborted.");
    						return;
    					}
					}
     			}
    			
			}
		}
		
		/**
		 * Returns true, if this worker thread is currently handling a task
		 * of the provided control
		 * 
		 * @param control
		 * @return
		 */
		public boolean isHandlingTask(ParallelExecutor<T> control) {
			synchronized (this) {
				if (task!=null && task.getControl().equals(control))
					return true;
				return false;
			}
		}
		
		/**
		 * Returns true, if this worker thread is currently handling a task
		 * belonging to the provided query
		 * 
		 * @param control
		 * @return
		 */
		public boolean isHandlingTask(int queryId) {
			if (queryId<0)
				return false;
			synchronized (this) {
				if (task!=null && task.getControl().getQueryId()==queryId)
					return true;
				return false;
			}
		}
		
		public void abort() {
			synchronized (this) {
				aborted = true;
			}
		}
		
		public boolean isAborted() {
			synchronized (this) {
				return aborted;
			}			
		}
		
		
		public boolean isIdle() {
			synchronized (this) {
				return idle;
			}
		}
	}
	
	
	
	/**
	 * Structure to maintain the status for a given control instance.
	 * 
	 * @author Andreas Schwarte
	 */
	protected class ControlStatus {
		public int waiting;
		public boolean done;
		public ControlStatus(int waiting, boolean done) {
			this.waiting = waiting;
			this.done = done;
		}
	}
	
	
	/**
	 * Monitor for the workers.
	 * 
	 * @author Andreas Schwarte
	 *
	 */
	protected class IdleWorkersMonitor extends Thread {
    	
		protected boolean closed = false;
		
    	@Override
    	public void run() {
    		
    		while (!Thread.interrupted() && !closed) {
    			
    			int _idle;
    			int req;
    			synchronized (taskQueue) {
    				_idle = idleWorkers;
    				req = taskQueue.size();
    				
    				System.out.println("Worker Status (" + name + "): " + _idle + " idle, requests in queue: " + req);
        			        			
        			for (WorkerThread w : workers) {
        				if (!w.isIdle()) {
        					System.out.println("Worker " + w.getId() + ": inTask=" + Boolean.toString(w.inTask) + ", task: " + w.task);
        				}
        			}
    			}    		
    			
    			
    			try {
					Thread.sleep(5000);
				} catch (InterruptedException e) {
					// ignore
				}
    		}
    		log.debug("Idle Worker Monitor for scheduler " + name + " closed.");
    	}
    	
    	public void close() {
    		closed = true;
    	}
    }
}
