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

package com.fluidops.fedx;

import java.io.FileInputStream;
import java.io.FileNotFoundException;
import java.io.IOException;
import java.util.Properties;

import org.apache.log4j.Logger;

import com.fluidops.fedx.cache.MemoryCache;
import com.fluidops.fedx.evaluation.concurrent.ControlledWorkerScheduler;
import com.fluidops.fedx.exception.FedXException;
import com.fluidops.fedx.exception.FedXRuntimeException;
import com.fluidops.fedx.monitoring.QueryLog;
import com.fluidops.fedx.monitoring.QueryPlanLog;


/**
 * Configuration properties for FedX based on a properties file. Prior to using this configuration
 * {@link #initialize(String)} must be invoked with the location of the properties file.
 * 
 * @author Andreas Schwarte
 *
 */
public class Config {

	protected static Logger log = Logger.getLogger(Config.class);
	
	private static Config instance = null;
	
	public static Config getConfig() {
		if (instance==null)
			throw new FedXRuntimeException("Config not initialized. Call Config.initialize() first.");
		return instance;
	}
	
	protected static void reset() {
		instance=null;
	}
	
	/**
	 * Initialize the configuration with the specified properties file.
	 * 
	 * @param fedxConfig
	 * 			the optional location of the properties file. If not specified the default configuration is used.
	 * 
	 * @throws FileNotFoundException
	 * @throws IOException
	 * @throws IllegalArgumentException
	 */
	public static void initialize(String ...fedxConfig) throws FedXException {
		if (instance!=null)
			throw new FedXRuntimeException("Config is already initialized.");
		instance = new Config();
		String cfg = fedxConfig!=null && fedxConfig.length==1 ? fedxConfig[0] : null;
		instance.init(cfg);
	}
	

	
	private Properties props;
	private Config() {
		props = new Properties();
	}
	
	private void init(String configFile) throws FedXException {
		if (configFile==null) {
			log.warn("No configuration file specified. Using default config initialization.");
			return;
		}
		log.info("FedX Configuration initialized from file '" + configFile + "'.");
		try {
			FileInputStream in = new FileInputStream(configFile);
			props.load( in );
			in.close();
		} catch (IOException e) {
			throw new FedXException("Failed to initialize FedX configuration with " + configFile + ": " + e.getMessage());
		}
	}
	
	public String getProperty(String propertyName) {
		return props.getProperty(propertyName);
	}
	
	public String getProperty(String propertyName, String def) {
		return props.getProperty(propertyName, def);
	}
	
	/**
	 * the base directory for any location used in fedx, e.g. for repositories
	 * 
	 * @return
	 */
	public String getBaseDir() {
		return props.getProperty("baseDir", "");
	}
	
	/**
	 * The location of the dataConfig.
	 * 
	 * @return
	 */
	public String getDataConfig() {
		return props.getProperty("dataConfig");
	}
	
	
	/**
	 * The location of the cache, i.e. currently used in {@link MemoryCache}
	 * 
	 * @return
	 */
	public String getCacheLocation() {
		return props.getProperty("cacheLocation", "cache.db");
	}
	
	/**
	 * The number of join worker threads used in the {@link ControlledWorkerScheduler}
	 * for join operations. Default is 20.
	 * 
	 * @return
	 */
	public int getJoinWorkerThreads() {
		return Integer.parseInt( props.getProperty("joinWorkerThreads", "20"));
	}
	
	/**
	 * The number of join worker threads used in the {@link ControlledWorkerScheduler}
	 * for join operations. Default is 20
	 * 
	 * @return
	 */
	public int getUnionWorkerThreads() {
		return Integer.parseInt( props.getProperty("unionWorkerThreads", "20"));
	}
	
	/**
	 * The block size for a bound join, i.e. the number of bindings that are integrated
	 * in a single subquery. Default is 15.
	 * 
	 * @return
	 */
	public int getBoundJoinBlockSize() {
		return Integer.parseInt( props.getProperty("boundJoinBlockSize", "15"));
	}
	
	/**
	 * Get the maximum query time in seconds used for query evaluation. Applied in CLI
	 * or in general if {@link QueryManager} is used to create queries.<p>
	 * 
	 * Set to 0 to disable query timeouts.
	 * 
	 * @return
	 */
	public int getEnforceMaxQueryTime() {
		return Integer.parseInt( props.getProperty("enforceMaxQueryTime", "30"));
	}
	
	/**
	 * Flag to enable/disable monitoring features. Default=false.
	 * 
	 * @return
	 */
	public boolean isEnableMonitoring() {
		return Boolean.parseBoolean( props.getProperty("enableMonitoring", "false"));	
	}
	
	/**
	 * Flag to enable/disable query plan logging via {@link QueryPlanLog}. Default=false
	 * The {@link QueryPlanLog} facility allows to retrieve the query execution plan
	 * from a variable local to the executing thread.
	 * 
	 * @return
	 */
	public boolean isLogQueryPlan() {
		return Boolean.parseBoolean( props.getProperty("monitoring.logQueryPlan", "false"));	
	}
	
	/**
	 * Flag to enable/disable query logging via {@link QueryLog}. Default=false
	 * The {@link QueryLog} facility allows to log all queries to a file. See 
	 * {@link QueryLog} for details. 
	 * 
	 * @return
	 */
	public boolean isLogQueries() {
		return Boolean.parseBoolean( props.getProperty("monitoring.logQueries", "false"));	
	}
	
	/**
	 * Returns the path to a property file containing prefix declarations as 
	 * "namespace=prefix" pairs (one per line).<p> Default: no prefixes are 
	 * replaced. Note that prefixes are only replaced when using the CLI
	 * or the {@link QueryManager} to create/evaluate queries.
	 * 
	 * Example:
	 * 
	 * <code>
	 * foaf=http://xmlns.com/foaf/0.1/
	 * rdf=http://www.w3.org/1999/02/22-rdf-syntax-ns#
	 * =http://mydefaultns.org/
	 * </code>
	 * 			
	 * @return
	 */
	public String getPrefixDeclarations() {
		return props.getProperty("prefixDeclarations");
	}
	
	/**
	 * The debug mode for worker scheduler, the scheduler prints usage stats regularly
	 * if enabled
	 * 
	 * @return
	 * 		false
	 */
	public boolean isDebugWorkerScheduler() {
		return Boolean.parseBoolean( props.getProperty("debugWorkerScheduler", "false"));
	}
	
	/**
	 * The debug mode for query plan. If enabled, the query execution plan is
	 * printed to stdout
	 * 
	 * @return
	 * 		false
	 */
	public boolean isDebugQueryPlan() {
		return Boolean.parseBoolean( props.getProperty("debugQueryPlan", "false"));
	}
	
	/**
	 * Set some property at runtime
	 * @param key
	 * @param value
	 */
	public void set(String key, String value) {
		props.setProperty(key, value);
	}
}
