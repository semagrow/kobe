package org.semanticweb.fbench.report;


import java.io.BufferedOutputStream;
import java.io.File;
import java.io.FileInputStream;
import java.io.FileOutputStream;
import java.io.OutputStream;
import java.util.List;
import java.util.Properties;

import javax.xml.datatype.DatatypeFactory;

import org.openrdf.model.Resource;
import org.openrdf.model.Statement;
import org.openrdf.model.URI;
import org.openrdf.model.Value;
import org.openrdf.model.ValueFactory;
import org.openrdf.model.impl.ValueFactoryImpl;
import org.openrdf.repository.RepositoryConnection;
import org.openrdf.repository.sail.SailRepository;
import org.openrdf.rio.RDFFormat;
import org.openrdf.rio.RDFWriter;
import org.openrdf.rio.Rio;
import org.openrdf.sail.memory.MemoryStore;
import org.semanticweb.fbench.Config;
import org.semanticweb.fbench.query.Query;


/**
 * Report configuration that writes RDF triples into result/result.rdf. RDF triples
 * can be visualized directly in the information workbench (tables, graphs, etc)
 * 
 * Note that query triples are written to result/queries_%querySet%.nt
 * 
 * @author as
 *
 */
public class RdfReportStream extends MemoryReportStream {

	protected SailRepository repo;
	protected RepositoryConnection conn;
	protected ValueFactory vf;
	
	public RdfReportStream() {
		super();
	}
	
	
	@Override
	public void writeData() throws Exception {
		
		init();
		
		writeEnvironment(); 
		
		writeLoadTimes();
		
		if (!Config.getConfig().isFill()) {
			writeQueryTuples();
		
			writeQueryEval();
		}
		
		finish();
		
	}

	
	protected void init() throws Exception {
		MemoryStore ms = new MemoryStore();
		repo = new SailRepository(ms);
		repo.initialize();
		conn = repo.getConnection();
		vf = ValueFactoryImpl.getInstance();
	}
	
	
	protected void writeEnvironment() throws Exception {
		
		URI eURI = createEvalURI();
		URI iURI = createInformationURI();
		
		Properties env = new Properties();
		FileInputStream in = new FileInputStream(new File( Config.getConfig().getEnvConfig()));	
		env.load( in );
		in.close();
		
		// type information + link
		conn.add( createStatement( eURI, RdfVocabulary.TYPE, RdfVocabulary.EVALUATION_TYPE));
		conn.add( createStatement( iURI, RdfVocabulary.TYPE, RdfVocabulary.INFORMATION_TYPE));
		conn.add( createStatement( eURI, RdfVocabulary.INFORMATION, iURI));
		
		// date, mode, description
		conn.add( createStatement(iURI, RdfVocabulary.DATE, vf.createLiteral( DatatypeFactory.newInstance().newXMLGregorianCalendar(evaluationDate))));
		conn.add( createStatement(iURI, RdfVocabulary.MODE, vf.createLiteral( (Config.getConfig().isFill() ? "fill" : "benchmark"))));
		conn.add( createStatement(iURI, RdfVocabulary.DATACONFIG, vf.createLiteral( (Config.getConfig().getDataConfig()))));

		
		// description:
		// fill mode: "Fill %ds.id" (=location name)"
		// benchmark: %description% (from config) (otherwise "Benchmark %configName% - Queryset") 
		String description;
		if (Config.getConfig().isFill()) {
			DatasetStats ds = datasetStats.get(0);
			description = "Fill " + ds.id;
		} else {
			String tmpDesc = Config.getConfig().getDescription();
			if (tmpDesc==null) {
				String dsName = new File(Config.getConfig().getDataConfig()).getName();
				description = "Benchmark " + dsName + " on " + Config.getConfig().getProperty("querySet").toLowerCase() + " queries"; 
			} else
				description = tmpDesc;
		}
		conn.add( createStatement(iURI, RdfVocabulary.DESCRIPTION, vf.createLiteral( description )));
		
		// other: memory, operating system 
		conn.add( createStatement(iURI, RdfVocabulary.MEMORY, vf.createLiteral( env.getProperty("memory", "n/a"))));
		conn.add( createStatement(iURI, RdfVocabulary.OPERATINGSYSTEM, vf.createLiteral( env.getProperty("operatingSystem", "n/a"))));
		conn.add( createStatement(iURI, RdfVocabulary.CPU, vf.createLiteral( env.getProperty("cpu", "n/a"))));
		conn.add( createStatement(iURI, RdfVocabulary.SOFTWARE, vf.createLiteral( env.getProperty("software", "n/a"))));
		conn.add( createStatement(iURI, RdfVocabulary.HARDDISK, vf.createLiteral( env.getProperty("harddisk", "n/a"))));
		conn.add( createStatement(iURI, RdfVocabulary.NOTE, vf.createLiteral( env.getProperty("note", "n/a"))));
		conn.add( createStatement(iURI, RdfVocabulary.ORGANIZATION, vf.createLiteral( env.getProperty("organization", "n/a"))));
		String timeout = Config.getConfig().isFill() ? "n/a" : Config.getConfig().getTimeout() + "";
		conn.add( createStatement(iURI, RdfVocabulary.TIMEOUT, vf.createLiteral( timeout )));
	}
	
	
	protected void writeLoadTimes() throws Exception {
		
		URI eURI = createEvalURI();
				
		for (DatasetStats d : datasetStats) {
			
			URI dstatsURI = createDatasetStatsURI(d);
			URI dURI = createDatasetURI(d);
			
			// link: evaluation -> dataset stats
			conn.add(createStatement(eURI, RdfVocabulary.DATASETSTATS, dstatsURI));
			
			// type information			
			conn.add( createStatement( dstatsURI, RdfVocabulary.TYPE, RdfVocabulary.DATASETSTATS_TYPE));
			conn.add( createStatement( dURI, RdfVocabulary.TYPE, RdfVocabulary.DATASET_TYPE));
			
			// link datasetset -> dataset
			conn.add( createStatement( dstatsURI, RdfVocabulary.DATASET, dURI));
			
			// stats information, i.e. loadtime
			conn.add( createStatement( dstatsURI, RdfVocabulary.NAME, vf.createLiteral(d.name)));
			conn.add( createStatement( dstatsURI, RdfVocabulary.LOADTIME, vf.createLiteral(d.loadTime)));
			conn.add( createStatement( dstatsURI, RdfVocabulary.DTYPE, vf.createLiteral(d.type)));
			
			// if fill mode create dataset entry
			if (Config.getConfig().isFill()) {
				URI dentryURI = createDatasetEntryURI(d);
				conn.add( createStatement( dentryURI, RdfVocabulary.TYPE, RdfVocabulary.DATASETENTRY_TYPE));
				conn.add( createStatement( dURI, RdfVocabulary.DATASETENTRY, dentryURI));
				conn.add( createStatement( dentryURI, RdfVocabulary.NAME, vf.createLiteral(d.name)));
				conn.add( createStatement( dentryURI, RdfVocabulary.LOCATION, vf.createLiteral(d.location)));
			}
			
			// otherwise add appropriate information to dataset itself
			else {
				conn.add( createStatement( dURI, RdfVocabulary.NAME, vf.createLiteral(d.name)));
				conn.add( createStatement( dURI, RdfVocabulary.DTYPE, vf.createLiteral(d.type)));
				conn.add( createStatement( dURI, RdfVocabulary.LOCATION, vf.createLiteral(d.location)));
			}
			
		}
	}
	
	protected void writeQueryTuples() throws Exception {
		
		MemoryStore ms = new MemoryStore();
		SailRepository queryRepo = new SailRepository(ms);
		queryRepo.initialize();
		RepositoryConnection queryConn = queryRepo.getConnection();
		
		for (Query q : queries) {
			URI qURI = createQueryURI(q);
			// {fbench:q.id; rdf:type; fbench:Query}
			// {fbench:q.id; fbench:query-sparql; ".."}
			queryConn.add( createStatement( qURI, RdfVocabulary.TYPE, RdfVocabulary.QUERY_TYPE));
			queryConn.add( createStatement( qURI, RdfVocabulary.ID, vf.createLiteral(q.getIdentifier())));
			queryConn.add( createStatement( qURI, RdfVocabulary.SPARQL, vf.createLiteral(q.getQuery())));
		}
		
		String file = Config.getConfig().getBaseDir() + "result\\queries_" + Config.getConfig().getProperty("querySet", "custom") + ".nt"; 
		File outFile = new File(file);
		OutputStream out = new BufferedOutputStream(new FileOutputStream(outFile));
		RDFWriter wr = Rio.createWriter(RDFFormat.NTRIPLES, out);
		wr.handleNamespace("rdf", "http://www.w3.org/1999/02/22-rdf-syntax-ns#");
		queryConn.export(wr);
		
		out.flush();
		out.close();
		queryConn.close();
		
	}
	
	
	protected void writeQueryEval() throws Exception {
		URI eURI = createEvalURI();
		
		for (Query q : queries) {
			List<QueryStats> qStats = queryEvaluation.get(q);
			URI qURI = createQueryURI(q);
			for (QueryStats qStat : qStats) {
				URI rURI = createRunURI(qStat);
				// {rURI; rdf:type; TestRun}
				// {rURI; fbench:testrun-query; ".."}
				// {rURI; fbench:testrun-duration; ".."}
				// {rURI; fbench:testrun-numberOfResults; ".."}
				// {rURI; fbench:testrun-run; ".."}
				conn.add( createStatement(rURI, RdfVocabulary.TYPE, RdfVocabulary.TESTRUN_TYPE));
				conn.add( createStatement(rURI, RdfVocabulary.QUERY, qURI));
				conn.add( createStatement(rURI, RdfVocabulary.RUNDURATION, vf.createLiteral(qStat.duration)));
				conn.add( createStatement(rURI, RdfVocabulary.NUMBEROFRESULTS, vf.createLiteral(qStat.numberOfResults)));
				conn.add( createStatement(rURI, RdfVocabulary.RUN, vf.createLiteral(qStat.run)));
			
				// {eURI; fbench:testrun; run}
				conn.add(createStatement(eURI, RdfVocabulary.TESTRUN, rURI));
			}
			
			// avg query stats
			URI qstatsURI = createQuerystatsURI(q);
			conn.add(createStatement(eURI, RdfVocabulary.QUERYSTATS, qstatsURI));
			conn.add(createStatement(qstatsURI, RdfVocabulary.TYPE, RdfVocabulary.QUERYSTATS_TYPE));
			conn.add(createStatement(qstatsURI, RdfVocabulary.AVGQUERYDURATION, vf.createLiteral( getAverageQueryDuration(q))));
			conn.add(createStatement(qstatsURI, RdfVocabulary.QUERY, qURI));
		}
	}
	
	
	protected void finish() throws Exception {
		String file = Config.getConfig().getBaseDir() + "result/result.nt"; 
		File outFile = new File(file);
		OutputStream out = new BufferedOutputStream(new FileOutputStream(outFile));
		RDFWriter wr = Rio.createWriter(RDFFormat.NTRIPLES, out);
		wr.handleNamespace("rdf", "http://www.w3.org/1999/02/22-rdf-syntax-ns#");
		conn.export(wr);
		
		out.flush();
		out.close();
		conn.close();
	}
	
	
	private Statement createStatement(Resource subject, URI predicate, Value object) {
		return RdfVocabulary.createStatement(subject, predicate, object);
	}
	
	/**
	 * @param s
	 * @return
	 * 			the URI relative to fbench namespace
	 */
	private URI createURI(String s) {
		return RdfVocabulary.createFURI(s);
	}
	
	
	private URI createDatasetStatsURI(DatasetStats d) {
		String name = d.name.replace("http://", "");
		name = name.replace("/", "_");
		return createURI("_" + evaluationID + "/d_" + name);
	}
	
	private URI createDatasetURI(DatasetStats d) {
		return createURI("d/" + d.id);
	}
	
	private URI createDatasetEntryURI(DatasetStats d) {
		String name = d.name.replace("http://", "");
		name = name.replace("/", "_");
		return createURI("e/" + name);
	}
	
	private URI createQueryURI(Query q) {
		return createURI("q/" + q.getIdentifier());
	}
	
	private URI createQuerystatsURI(Query q) {
		return createURI("_" + evaluationID + "/q_" + q.getIdentifier() );
	}
	
	private URI createRunURI(QueryStats qStats) {
		return createURI("_" + evaluationID + "/q_" + qStats.query.getIdentifier() + "/run_" + qStats.run);
	}
	
	private URI createEvalURI() {
		return createURI("_" + evaluationID );
	}
	
	private URI createInformationURI() {
		return createURI("_" + evaluationID + "/information" );
	}
	
	
}
