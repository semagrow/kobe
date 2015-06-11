package org.semanticweb.fbench.evaluation;


import java.io.File;
import java.util.List;

import org.semanticweb.fbench.Config;
import org.semanticweb.fbench.misc.FileUtil;
import org.semanticweb.fbench.query.Query;

import com.fluidops.fedx.FedXFactory;
import com.fluidops.fedx.FederationManager;
import com.fluidops.fedx.structures.Endpoint;


public class FedXEvaluation extends SesameEvaluation{

	public FedXEvaluation() {
		super();
	}
	
	@Override
	public void initialize() throws Exception {
		log.info("Performing Sesame-Extension Initialization...");
		
		File dataConfig = FileUtil.getFileLocation( Config.getConfig().getDataConfig() );
		List<Endpoint> endpoints = DataConfigReader.loadFederationMembers(dataConfig, report);
		
		com.fluidops.fedx.Config.initialize(Config.getConfig().getProperty("fedxConfig"));
				
		sailRepo = FedXFactory.initializeFederation(endpoints);
		if (!Config.getConfig().isFill())
			conn = sailRepo.getConnection();
		
		log.info("Sesame Repository successfully initialized.");
	}

	@Override
	public int runQuery(Query q, int run) throws Exception {
		
		// use fedbench optional config param fedxClearCache to specify up to
		// which run the cache is cleared prior to executing the query 
		// 0 means disabled
		int clearCacheRuns = Integer.parseInt(Config.getConfig().getProperty("fedxClearCache", "0"));
		if (clearCacheRuns >= run)
			FederationManager.getInstance().getCache().clear();
		return super.runQuery(q, run);
	}
	
	@Override
	public void reInitialize() throws Exception {
		
		log.info("Reinitializing repository and connection due to error in past results.");
		FederationManager.getInstance().reset();	// reset worker thread pools
		
		sailRepo.initialize();
		conn = sailRepo.getConnection();
		log.debug("reinitialize done.");
	}	

	@Override
	public void finish() throws Exception {
		FederationManager.getInstance().getCache().persist();
		super.finish();
	}
}
