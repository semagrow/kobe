package org.semanticweb.fbench;

import org.apache.log4j.Logger;
import org.semanticweb.fbench.evaluation.Evaluation;
import org.semanticweb.fbench.query.QueryManager;
import org.semanticweb.fbench.setup.Setup;


/**
 * Main class to execute a benchmark evaluation.<p>
 * 
 * <ul>
 *  <li>noArgs -> default config path (config\config.prop)</li>
 *  <li>"config\myConfig.prop" -> use the specified config</li>
 *  <li>"-fill" -> activated fill mode</li>
 * </ul>
 * 
 * @author (mz), as
 */
public class FederationEval {

	public static Logger log;
	
			
	public static void main(String[] args) {
		
		// initialize config: load properties
		// if no arg specified, config's location is config\config.prop
		try {
			Config.initialize(args);
		} catch (Exception e) {
			System.out.println("Could not initialize Config: " + e.getMessage());
			System.exit(1);
		}
		
		log = Logger.getLogger(FederationEval.class);
		
		if (Config.getConfig().isSetup()) {
			log.info("Setup mode enabled. Beginning preparation of data sources.");
			try {
				Setup.prepareDataSources();
				System.exit(0);
			} catch (Exception e) {
				log.fatal("Error during setup (" + e.getClass().getSimpleName() + "): " + e.getMessage());
				log.debug("Exception details", e);
				System.exit(1);
			}
		}
		
		// initialize the query manager, i.e. load all queries
		try {
			QueryManager.initialize();
		} catch (Exception e) {
			log.fatal("Could not initialize query manager: " + e.getMessage());
			log.debug("Exception details", e);
			System.exit(1);
		}
		
		
		// Determine the Evaluation class to be used and run the evaluation
		try {
			Evaluation eval = (Evaluation)Class.forName(Config.getConfig().getEvaluationClass()).newInstance();
			eval.run();
		} catch (Exception e) {
			log.fatal("Error while performing evaluation (" + e.getClass().getSimpleName() + "): " + e.getMessage());
			log.debug("Exception details", e);
			System.exit(1);
		}
		
		System.exit(0);
	}

}
