package org.semanticweb.fbench.misc;

import org.apache.log4j.Logger;


public class TimedInterrupt {
	
	public static Logger log = Logger.getLogger(TimedInterrupt.class);
	
	public boolean run(Runnable r, long timeout) {
	
		MyThread t = new MyThread(this, r);
		
		synchronized (this) {
			t.start();
			try {
//				log.info("waiting for " + timeout);
				this.wait(timeout);
			} catch (InterruptedException e) {
				// ignore
			}
		}
		
//		log.info("ready again");
		
		if (!t.isFinished()) {
//			log.info("Stopping task in thread " + t.getName() + ". Timeout reached.");
			t.interrupt();
			t.stop();
			return false;
		}
		
		return true;
			
	}
	
	
	
	
	
	protected class MyThread extends Thread {
		
		protected final TimedInterrupt parent;
		protected final Runnable r;
		
		public MyThread(TimedInterrupt parent, Runnable r) {
			super();
			this.parent = parent;
			this.r = r;
		}

		private boolean finished = false;
		
		@Override 
		public void run() {
//			log.info("Starting task in thread " + Thread.currentThread().getName());
			r.run();
			finished = true;
			synchronized (parent) {
				parent.notify();
			}
		}
		
		public boolean isFinished() {
			return finished;
		}
		
	}
}
