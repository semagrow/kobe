package org.semanticweb.fbench.report;

import java.util.ArrayList;
import java.util.GregorianCalendar;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

import org.apache.log4j.Logger;
import org.semanticweb.fbench.query.Query;

public abstract class MemoryReportStream implements ReportStream {

	protected static class QueryStats {
		public Query query;
		public long duration;
		public int run;
		public int numberOfResults;
		public QueryStats(Query query, long duration, int run,
				int numberOfResults) {
			super();
			this.query = query;
			this.duration = duration;
			this.run = run;
			this.numberOfResults = numberOfResults;
		}		
	}
	
	protected static class DatasetStats {
		public String id;
		public String name; 
		public String location;
		public String type;
		public long loadTime;
		public DatasetStats(String id, String name, String location, String type,
				long loadTime) {
			super();
			this.id = id;
			this.name = name;
			this.location = location;
			this.type = type;
			this.loadTime = loadTime;
		}		
	}
	
	public static Logger log = Logger.getLogger(MemoryReportStream.class);
	
	protected String evaluationID;
	protected GregorianCalendar evaluationDate;
	protected String dataConfig;
	protected List<String> querySet;
	protected int numberOfQueries;
	protected int numberOfRuns;
	protected long initializationDuration;
	protected long overallDuration;
	protected List<Query> queries;
	protected Map<Query, List<QueryStats>> queryEvaluation;
	protected List<Long> runDurations;
	protected List<DatasetStats> datasetStats;
	
	public MemoryReportStream() {
		this.queryEvaluation = new HashMap<Query, List<QueryStats>>();
		this.queries = new ArrayList<Query>();
		this.runDurations = new ArrayList<Long>();
		this.datasetStats = new ArrayList<DatasetStats>();
	}
	
	@Override
	public void beginEvaluation(String dataConfig, List<String> querySet,
			int numberOfQueries, int numberOfRuns) {
		this.dataConfig = dataConfig;
		this.querySet = querySet;
		this.numberOfQueries = numberOfQueries;
		this.numberOfRuns = numberOfRuns;	
	}

	@Override
	public void beginQueryEvaluation(Query query, int run) {
		if (run==1)
			queries.add(query);
	}

	@Override
	public void beginRun(int run, int totalNumberOfRuns) {
		;		
	}

	@Override
	public void close() throws Exception {
		;		
	}

	@Override
	public void endEvaluation(long duration) {
		this.overallDuration = duration;
		try {
			writeData();
		} catch (Exception e) {
			log.error("Error writing message in " + this.getClass().getSimpleName() + "#writeData() (" + e.getClass().getSimpleName() + "): " + e.getMessage());
			log.debug("Exception details:", e);
			throw new RuntimeException("Error while writing report", e);
		}		
	}

	@Override
	public void endQueryEvaluation(Query query, int run, long duration,
			int numberOfResults) {
		
		List<QueryStats> list = queryEvaluation.get(query);
		if (list == null) {
			list = new ArrayList<QueryStats>();
			queryEvaluation.put(query, list);
		}
		// check for duplicate results
		for (QueryStats qs : list) {
			if (qs.run==run) {
				qs.duration = qs.duration==-1 || duration==-1 ? -1 : qs.duration;
				return;
			}
		}
		list.add( new QueryStats(query, duration, run, numberOfResults));		
	}

	@Override
	public void endRun(int run, int totalNumberOfRuns, long duration) {
		runDurations.add(duration);		
	}

	@Override
	public void initializationBegin() {
		;		
	}

	@Override
	public void initializationEnd(long duration) {
		this.initializationDuration = duration;		
	}

	@Override
	public void open() throws Exception {
		this.evaluationID = Long.toString( System.currentTimeMillis() );
		this.evaluationDate = new GregorianCalendar();
	}

	@Override
	public void reportDatasetLoadTime(String id, String name, String location,
			String type, long duration) {
		this.datasetStats.add( new DatasetStats(id, name, location, type, duration));
	}

	
	public long getAverageQueryDuration(Query q) {
		long sum = 0;
		// TODO what about "some" timeouts (timeout duration=-1)? 
		List<QueryStats> l = queryEvaluation.get(q);
		for (QueryStats qstat : l)
			sum += (qstat.duration + 1);		// add normalization of 1ms
		if (sum==0)
			return -1;	// all requests timeout
		return sum / l.size();
	}
	
	public abstract void writeData() throws Exception;
}
