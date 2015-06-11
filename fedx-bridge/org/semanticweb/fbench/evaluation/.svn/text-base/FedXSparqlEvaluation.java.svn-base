package org.semanticweb.fbench.evaluation;


import java.io.File;
import java.io.IOException;
import java.util.List;

import org.semanticweb.fbench.Config;
import org.semanticweb.fbench.misc.FileUtil;
import org.semanticweb.fbench.query.Query;
import org.semanticweb.fbench.report.SparqlQueryRequestReport;

import com.fluidops.fedx.FedXFactory;
import com.fluidops.fedx.FederationManager;
import com.fluidops.fedx.structures.Endpoint;


public class FedXSparqlEvaluation extends SesameSparqlEvaluation{

	private FedXMonitoringReport fedxReport = null;
	
	public FedXSparqlEvaluation() {
		super();
	}
	
	@Override
	public void initialize() throws Exception {
		log.info("Performing Sesame-Extension Initialization...");
		
		com.fluidops.fedx.Config.initialize(Config.getConfig().getProperty("fedxConfig"));
		
		File dataConfig = FileUtil.getFileLocation( Config.getConfig().getDataConfig() );
		List<Endpoint> endpoints = DataConfigReader.loadFederationMembers(dataConfig, report);
		
		if (Boolean.parseBoolean(Config.getConfig().getProperty("fedxRequestReport", "false"))) {
			com.fluidops.fedx.Config.getConfig().set("enableMonitoring", "true");
			fedxReport = new FedXMonitoringReport(endpoints);
		}
		
		sailRepo = FedXFactory.initializeFederation(endpoints);
		if (!Config.getConfig().isFill())
			conn = sailRepo.getConnection();
		
		new File("_shutdown").delete();
		for (File f : new File(".").listFiles())
			if (f.getName().endsWith(".pid"))
				f.delete();
		repoInformation = getRepoInformation();
		
		if (Config.getConfig().isSparqlRequestReport()) {
			sparqlReport = new SparqlQueryRequestReport();
			sparqlReport.init(repoInformation);
		}
				
		reinitializeSystem();
				
		log.info("Sesame Repository successfully initialized.");
	}


	protected void reinitializeSystem() throws Exception {
		log.info("Reinitializing system...");
		if (FederationManager.isInitialized())
			FederationManager.getInstance().reset();
		
		if (runningServers!=0) {
			log.info("Wating for graceful shutdown of servers, give them " + extraWait + "ms time.");
			File cFile = new File("_shutdown");
			cFile.createNewFile();
			System.gc();
			Thread.sleep(extraWait);
			cFile.delete();
			log.info("Checking for processes that did not shutdown gracefully and killing them by force..");
			checkAndKillServers();			
			runningServers=0;
		}
		
		
		log.info("Deleting possible locks in the file system.");
		for (RepoInformation r : repoInformation) {
			if (r.repoLoc == null)
				continue;
			for (File f : r.repoLoc.listFiles())
				if (f.isDirectory() && f.getName().equals("lock")) {
					if (!deleteFolder(f))
						log.fatal("Lock folder " + f.getAbsolutePath() + " could not be deleted!");
				}
		}
		
		log.info("Starting server processes ... ");
		int port = 10000;
		for (RepoInformation r : repoInformation) {
			if (r.repoLoc == null)
				continue;
			startSparqlServer(r.repoLoc, port++);
			runningServers++;
		}
				
		log.info("Waiting for " + extraWait + " ms to give server time for initialization");
		System.gc();
		Thread.sleep(extraWait);
		
		log.debug("Loading repositories from scratch.");
		sailRepo.initialize();
		conn = sailRepo.getConnection();
		log.debug("Reinitialize done.");
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
	public void finish() throws Exception {
		FederationManager.getInstance().getCache().persist();
		if (fedxReport!=null)
			fedxReport.finish();
		super.finish();
	}

	@Override
	protected void queryRunEnd(Query query, boolean error)
	{
		if (fedxReport!=null)
			try	{
				fedxReport.handleQuery(query);
			} catch (IOException e)	{
				throw new RuntimeException(e);
			}
		super.queryRunEnd(query, error);
	}
	
	
}
