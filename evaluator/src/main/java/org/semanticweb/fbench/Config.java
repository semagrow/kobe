package org.semanticweb.fbench;

import org.semanticweb.fbench.evaluation.SesameSparqlEvaluationReactive;
import org.semanticweb.fbench.misc.ArgumentParser;
import org.semanticweb.fbench.provider.SPARQLProvider;
import org.semanticweb.fbench.report.*;

import java.io.FileInputStream;
import java.io.FileNotFoundException;
import java.io.IOException;
import java.util.ArrayList;
import java.util.List;
import java.util.Properties;




/**
 * Configuration for the benchmark evaluation. Is initialized with the 
 * specified properties file (i.e. config.prop) at runtime.
 * 
 * @author as
 *
 */
public class Config {
	
	public static class Property {
		public String key;
		public String value;
		public Property(String key, String value) {
			super();
			this.key = key;
			this.value = value;
		}		
	}
	

	private static Config instance = null;
	
	public static Config getConfig() {
		if (instance==null)
			throw new RuntimeException("Config not initialized. Call Config.load() first.");
		return instance;
	}
	
	public static void initialize(String[] args) throws FileNotFoundException, IOException, IllegalArgumentException {
		/*
		if (System.getProperty("log4j.configuration")==null) {
			//System.setProperty("log4j.configuration", "file:config/log4j.properties");
			System.setProperty("log4j.configuration", "log4j.properties");
		}
		*/
		LogUtils.setMDC();
		// necessary for RDFXML format, e.g. for Jamendo dataset, to not abort with RDFParseException
		if (System.getProperty("entityExpansionLimit")==null)
			System.setProperty("entityExpansionLimit", "10000000");
		
		instance = new Config();
		
		for (Property p : ArgumentParser.parseArguments(args))
			instance.setProperty(p.key, p.value);
		
		instance.init(instance.getProperty("configFile"));
	}
	
	
	
	private Properties props;
	
	private Config() {
		props = new Properties();
	}
	
	public void setProperty(String prop, String value) {
		props.setProperty(prop, value);
	}
	
	private void init(String configFile) throws IOException {
        if (configFile != null) {
            FileInputStream in = new FileInputStream(configFile);
            props.load(in);
            in.close();
        }
	}
	
	
	public String getProperty(String propertyName) {
		return props.getProperty(propertyName);
	}
	
	public String getProperty(String propertyName, String def) {
		return props.getProperty(propertyName, def);
	}

	
	/**
	 * @return
	 * 		a directory that contains a list of text files, and each text file is a query.
     * 	    if none specified, then /etc/querySet/ directory is used
	 */
	public String getQuerySetPath() throws IllegalArgumentException {
        return props.getProperty("querySet", "/queries");
	}

    public List<String> getQuerySet() {
        List<String> l = new ArrayList<>();
        l.add(Config.getConfig().getQuerySetPath());
        return l;
    }

	/**
	 * @return
	 * 		the base directory for the evaluation benchmark or the empty string if none specified
	 */
	public String getBaseDir() {
		return props.getProperty("baseDir", "");
	}
	
	
	/**
	 * @return
	 * 		true, if fill mode is enabled (commandline arg "-fill", i.e. no queries are executed)
	 */
	public boolean isFill() {
		return Boolean.parseBoolean( props.getProperty("fill", "false"));
	}

	/**
	 * 
	 * @return
	 * 		true, if setup mode is enabled (commandline arg "-setup", i.e. no queries are executed)
	 */
	public boolean isSetup() {
		return Boolean.parseBoolean( props.getProperty("setup", "false"));
	}

	/**
	 * @return
	 * 		true, if property showResults is set
	 */
	public boolean isShowResults() {
		return Boolean.parseBoolean( props.getProperty("showResults", "false"));
	}

	/**
	 * @return
	 * 		true, if property debugMode is set
	 */
	public boolean isDebugMode() {
		return Boolean.parseBoolean( props.getProperty("debugMode", "false"));
	}

	/**
	 * Return the report stream implementation to be used
	 * 
	 * default: 
	 * 	 a) com.fluidops.iwb.benchmark.report.SimpleReportStream (if debug mode is on)
	 *   b) com.fluidops.iwb.benchmark.report.CsvReportStream (otherwise)
	 * 
	 * @return
	 * 		the reportStream setting, i.e. the fully qualified class that shall be used for reporting
	 */
	public String getReportStream() {
		String def = isDebugMode() ? SimpleReportStream.class.getCanonicalName() : CsvReportStream2.class.getCanonicalName();
		return props.getProperty("reportStream", def);
	}

	/**
	 * @return
	 * 		the timeout setting in ms, default is 0ms (=OFF)
	 */
	public long getTimeout() {
		return Long.parseLong(props.getProperty("timeout", "0"));
	}
		
	
	/**
	 * @return
	 * 		the number of evaluation runs, default is 1
	 */
	public int getEvalRuns() {
		return Integer.parseInt(props.getProperty("evalRuns", "1"));
	}
	
	/**
	 * @return
	 * 	 	the data configuration or baseDir/dataconfig.ttl if not specified
	 */
	public String getDataConfig() {
		return props.getProperty("dataConfig", getBaseDir() + "dataConfig.ttl");
	}
	
	/**
	 * 
	 * @return
	 * 		the evaluationClass setting, i.e. the fully qualified class that shall be used for evaluation
	 * 		default: org.semanticweb.fbench.evaluation.SesameSparqlEvaluationReactive
	 */
	public String getEvaluationClass() {
		return props.getProperty("evaluationClass", SesameSparqlEvaluationReactive.class.getCanonicalName());
	}

	/**
	 * 
	 * @return
	 * 		the class that shall be used for early results monitoring.
	 * 		default: org.semanticweb.fbench.report.NoOpEarlyResultsMonitor
	 */
	public String getEarlyResultsMonitorClass() {
		return props.getProperty("earlyResultsMonitorClass", NoOpEarlyResultsMonitor.class.getCanonicalName());
	}
	
	/**
	 * 
	 * @return
	 * 		the envConfig setting, i.e. the location of the environment properties that are used in {@link RdfReportStream}
	 * 		default: null
	 */
	public String getEnvConfig()  {
		return props.getProperty("envConfig");
	}
	
	/**
	 * 
	 * @return
	 * 		the description setting, i.e. a meaningful short description. Ex: "Endpoint Federation (Sesame) Cross Domain Queries"
	 * 		default null
	 */
	public String getDescription() {
		return props.getProperty("description");
	}
	
	/**
	 * 
	 * @return
	 * 		true iff a proxyUrl is specified
	 */
	public boolean useGlobalProxy() {
		return getProxyUrl()!=null;
	}
	
	/**
	 * 
	 * @return
	 * 		the proxyUrl setting, i.e. (if specified) the url which shall be used
	 * 		as proxy for HTTP SPARQL requests. Refer to {@link SPARQLProvider} for
	 * 		some further information
	 */
	public String getProxyUrl() {
		return props.getProperty("proxyUrl");
	}

	/**
	 * 
	 * @return
	 * 		the delay in ms, which is added to each sparql request. Note that this configuration (if enabled) 
	 * 		works only with local SPARQL server, since the property is added in SPARQLServlet2, i.e. it is 
	 * 		given as argument to the automatically started servers. Default is -1 (disabled)
	 */
	public int getSparqlRequestDelay() {
		return Integer.parseInt(props.getProperty("sparqlRequestDelay", "-1"));
	}

    /**
     *
     * @return
     * 		time break after query in ms. Default is 1000ms
     */
	public long getBreakAfterQuery() {
		return Long.parseLong(props.getProperty("breakAfterQuery", "1000"));
	}

    /**
	 * 
	 * @return
	 * 		the sparqlRequestReport setting, i.e. if the number of requests sent to the endpoint
	 *  	shall be reported per query (see {@link SparqlQueryRequestReport}) 
	 *  	
	 *  	true: report stats to result\sparql_stats.csv
	 *  	false: disabled
	 *  
	 *  	default true
	 */
	public boolean isSparqlRequestReport() {
		return Boolean.parseBoolean(props.getProperty("sparqlRequestReport", "false"));
	}

	/**
	 *
	 * @return
	 * 		Name of the experiment. Used for logging purposes.
	 *
	 *  	default kobeExperiment
	 */
	public String getExperimentName() {
		return props.getProperty("experimentName", "kobeExperiment");
	}
}
