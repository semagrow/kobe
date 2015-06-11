package org.semanticweb.fbench.misc;

import java.util.Calendar;
import java.util.GregorianCalendar;

import javax.xml.bind.DatatypeConverter;

import org.apache.log4j.Logger;
import org.openrdf.repository.Repository;
import org.openrdf.repository.RepositoryConnection;
import org.openrdf.repository.RepositoryException;
import org.openrdf.repository.sail.SailRepository;

public class Utils {

	public static String dateToXsd(Calendar calendar) {
		return DatatypeConverter.printDateTime(calendar);
	}
	
	public static String nowToXsd() {
		return dateToXsd( new GregorianCalendar() );
	}
	
	
	public static long getPID() {
		String processName = java.lang.management.ManagementFactory.getRuntimeMXBean().getName();
		return Long.parseLong(processName.split("@")[0]);
	}


	public static boolean closeConnectionTimeout(final RepositoryConnection conn, long timeout) {
		TimedInterrupt t = new TimedInterrupt();
		
		return t.run( new Runnable() {
			@Override
			public void run() {
				try {
					conn.close();
					
				} catch (RepositoryException e) {
					Logger.getLogger(this.getClass()).error("Error closing conenction: " + e.getMessage());
				}						
			}
		}, timeout);
	}
	
//	public static boolean shutdownRepositoryTimeout(final SailRepository repo, long timeout) {
	public static boolean shutdownRepositoryTimeout(final Repository repo, long timeout) {
		TimedInterrupt t = new TimedInterrupt();
		
		return t.run( new Runnable() {
			@Override
			public void run() {
				try {
					repo.shutDown();
					
				} catch (RepositoryException e) {
					Logger.getLogger(this.getClass()).error("Error shutting down repository: " + e.getMessage());
				}						
			}
		}, timeout);
	}
	
	
	
	public static void main(String[] args) {
		System.out.println("now: " + nowToXsd());
	}
}
