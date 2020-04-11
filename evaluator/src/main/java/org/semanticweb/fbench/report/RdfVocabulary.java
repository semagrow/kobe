package org.semanticweb.fbench.report;

import org.openrdf.model.Resource;
import org.openrdf.model.Statement;
import org.openrdf.model.URI;
import org.openrdf.model.Value;
import org.openrdf.model.impl.ValueFactoryImpl;

public class RdfVocabulary {

	
	/**
	 * namespace for fbench
	 */
	public static String FBENCH = "http://semanticweb.org/fbench/";
	
	public static String RDF = "http://www.w3.org/1999/02/22-rdf-syntax-ns#";
	
	
	/* Evaluation (object) properties, e.g env */
	public static URI EVALUATION_TYPE = createFURI("Evaluation");
	public static URI DATASETSTATS = createFURI("datasetstats");
	public static URI TESTRUN = createFURI("testrun");
	public static URI QUERYSTATS = createFURI("querystats");
	public static URI INFORMATION = createFURI("information");
	
	/* Information (data) properties */
	public static URI INFORMATION_TYPE = createFURI("Information");
	public static URI DATACONFIG = createFURI("dataConfig");
	public static URI MODE = createFURI("mode");
	public static URI MEMORY =  createFURI("memory");
	public static URI OPERATINGSYSTEM =  createFURI("os");
	public static URI CPU =  createFURI("cpu");
	public static URI SOFTWARE =  createFURI("software");
	public static URI NOTE =  createFURI("note");
	public static URI HARDDISK =  createFURI("harddisk");
	public static URI ORGANIZATION =  createFURI("organization");
	public static URI TIMEOUT =  createFURI("timeout");
	
	/* Datastats (object/data) properties */
	public static URI DATASETSTATS_TYPE = createFURI("Datasetstats");
	public static URI DATASET = createFURI("dataset");
	public static URI LOADTIME = createFURI("loadtime");
	
	
	public static URI DATASET_TYPE = createFURI("Dataset");
	public static URI DATASETENTRY = createFURI("datasetentry");
	public static URI LOCATION = createFURI("location");
	
	public static URI DATASETENTRY_TYPE = createFURI("Datasetentry");
	public static URI NAME = createFURI("name");
	public static URI DTYPE = createFURI("dtype");	
	
	
	public static URI QUERY_TYPE = createFURI("Query");
	public static URI ID = createFURI("id");
	public static URI SPARQL = createFURI("sparql");
	
	public static URI QUERYSTATS_TYPE = createFURI("Querystats");
	public static URI AVGQUERYDURATION = createFURI("avgQueryDuration");
	
	
	public static URI TESTRUN_TYPE = createFURI("Testrun");
	public static URI QUERY = createFURI("query");
	public static URI RUN = createFURI("run");
	public static URI RUNDURATION = createFURI("duration");
	public static URI NUMBEROFRESULTS = createFURI("numberOfResults");
	
	
	
	public static URI TYPE = createURI(RDF, "type");
	
	
    public static final String dc = "http://purl.org/dc/elements/1.1/";
    public static final URI TITLE = (ValueFactoryImpl.getInstance()).createURI(dc, "title");
    public static final URI DATE = (ValueFactoryImpl.getInstance()).createURI(dc, "date");
    public static final URI CREATOR = (ValueFactoryImpl.getInstance()).createURI(dc, "creator");
    public static final URI FORMAT = (ValueFactoryImpl.getInstance()).createURI(dc, "format");
    public static final URI DESCRIPTION = (ValueFactoryImpl.getInstance()).createURI(dc, "description");
    public static final URI LICENSE = (ValueFactoryImpl.getInstance()).createURI(dc, "license");
    public static final URI HAS_VERSION = (ValueFactoryImpl.getInstance()).createURI(dc, "hasVersion");
    public static final URI dcModified = (ValueFactoryImpl.getInstance()).createURI(dc, "modified");
	
	
	public static URI createURI(String s) {
		return (ValueFactoryImpl.getInstance()).createURI(s);
	}
	
	public static URI createURI(String ns, String s) {
		return (ValueFactoryImpl.getInstance()).createURI(ns, s);
	}
		
	public static URI createFURI(String s) {
		return (ValueFactoryImpl.getInstance()).createURI(FBENCH, s);
	}

	public static Statement createStatement(Resource subject, URI predicate, Value object) {
		return (ValueFactoryImpl.getInstance()).createStatement(subject, predicate, object);
	}
}
