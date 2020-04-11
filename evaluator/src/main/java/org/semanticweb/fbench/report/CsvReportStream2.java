package org.semanticweb.fbench.report;

import java.io.BufferedWriter;
import java.io.File;
import java.io.FileWriter;
import java.util.List;

import org.semanticweb.fbench.Config;
import org.semanticweb.fbench.query.Query;


/**
 * Report query evaluation to %baseDir%\result\result.csv
 * Report loadTimes to file at %baseDir%\result\loadtimes.csv
 * 
 * Format:
 * <code>
 * eval: Query;run1;run2;...;runN;avg;numResults;minRes;maxRes
 * loadTimes: id;name;location;type;duration;
 * <code>
 * 
 * @author as
 *
 */
public class CsvReportStream2 extends MemoryReportStream {

	@Override
	public void writeData() throws Exception {
		writeLoadTimes();
		writeResult();
	}
	
	
	private void writeLoadTimes() throws Exception {
		
		String file2 = Config.getConfig().getBaseDir() + "result/loadTimes.csv"; 
		BufferedWriter loadOut = new BufferedWriter( new FileWriter(file2));
		loadOut.append("id;name;location;type;duration;\r\n");
		
		for (DatasetStats d : datasetStats) {
			loadOut.append(d.id+";"+d.name+";"+d.location+";"+d.type+";"+d.loadTime+";\r\n");
		}
		
		loadOut.flush();
		loadOut.close();
	}
	
	
	private void writeResult() throws Exception {
		
		String file = Config.getConfig().getBaseDir() + "result/result.csv"; 
		File outFile = new File(file);
		BufferedWriter bw = new BufferedWriter( new FileWriter( outFile ));
		
		// Query;run1;run2;...;runN;avg;numResults;minRes;maxRes
		bw.append("Query;");
		for (int i=1; i<=this.numberOfRuns; i++)
			bw.append("run" + i + ";");
		bw.append("avg;numResults;minRes;maxRes;\r\n");
		
		// write query results
		for (Query q : this.queries) {
			
			bw.append(q.getIdentifier()).append(";");		// Query;
			
			int minResults = Integer.MAX_VALUE;
			int maxResults = Integer.MIN_VALUE;
			int sumResults = 0;
			int resultsCount = 0;
			
			// run results
			List<QueryStats> qStats = queryEvaluation.get(q);
			for (QueryStats qStat : qStats) {
				bw.append(Long.toString(qStat.duration) + ";");
				minResults = Math.min(minResults, qStat.numberOfResults);
				maxResults = Math.max(maxResults, qStat.numberOfResults);
				if (qStat.duration>=0) {
					sumResults += qStat.numberOfResults;
					resultsCount++;
				}
			}
			
			// avgDuration
			// XXX this can be improved, and integrated into the loop above, but I am lazy ;-)
			bw.append( Long.toString(getAverageQueryDuration(q)) + ";");
			
			// number of results
			bw.append( Integer.toString( sumResults==0 ? maxResults : sumResults / resultsCount) + ";");
			bw.append( Integer.toString(minResults) + ";");
			bw.append( Integer.toString(maxResults) + ";\r\n");
			
		}
		
		bw.flush();
		bw.close();
	}

}
