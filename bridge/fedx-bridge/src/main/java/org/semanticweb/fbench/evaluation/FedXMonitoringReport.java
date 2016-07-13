package org.semanticweb.fbench.evaluation;

import java.io.BufferedWriter;
import java.io.FileWriter;
import java.io.IOException;
import java.util.List;

import org.apache.log4j.Logger;
import org.semanticweb.fbench.query.Query;
import org.semanticweb.fbench.report.SparqlQueryRequestReport;

import com.fluidops.fedx.monitoring.MonitoringImpl.MonitoringInformation;
import com.fluidops.fedx.monitoring.MonitoringService;
import com.fluidops.fedx.monitoring.MonitoringUtil;
import com.fluidops.fedx.structures.Endpoint;

public class FedXMonitoringReport
{

	
	public static Logger log = Logger.getLogger(SparqlQueryRequestReport.class);

	protected List<Endpoint> endpoints;
	protected BufferedWriter bout;
	
	
	public FedXMonitoringReport(List<Endpoint> endpoints) throws IOException {
		this.endpoints = endpoints;
		bout = new BufferedWriter( new FileWriter("result/fedx_stats.csv"));
		bout.write("query;");
		for (Endpoint e : endpoints)
			bout.write(e.getId() + ";");
		bout.write("total;");
		bout.write("\r\n");
	}
	

	public void handleQuery(Query query) throws IOException {
		log.info("Reporting query statatistics");
		bout.write(query.getIdentifier() + ";");
		
		MonitoringService m = MonitoringUtil.getMonitoringService();
		try {
			int total = 0;
			for (Endpoint e : endpoints) {			
				try {					
					MonitoringInformation mInfo = m.getMonitoringInformation(e);
					int requestCount = mInfo!=null ? mInfo.getNumberOfRequests() : 0;
					
					total += requestCount;
					
					log.info("## " + e.getId() + " => " + requestCount + " requests");
					bout.write(requestCount + ";");
				} catch (Exception ex) {
					log.warn("Exception (" + ex.getClass() + "):" + ex.getMessage());
					bout.write("-1;");
				}
			}
			bout.write(total + ";");
			bout.write("\r\n");
			bout.flush();
		} finally {	
			m.resetMonitoringInformation();
		}
		
	}

	public void finish() throws IOException {
		bout.flush();
		bout.close();
	}
}
