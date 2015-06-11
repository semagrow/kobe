package org.semanticweb.fbench.provider;

import java.io.BufferedReader;
import java.io.BufferedWriter;
import java.io.File;
import java.io.FileReader;
import java.io.FileWriter;
import java.io.IOException;
import java.io.StringReader;

import org.apache.log4j.Logger;
import org.openrdf.model.Resource;
import org.openrdf.repository.RepositoryConnection;
import org.openrdf.repository.RepositoryException;
import org.openrdf.rio.RDFFormat;
import org.openrdf.rio.RDFParseException;

/**
 * Specialized data reader for the Geonames Dump
 * 
 * Format:
 * 1: URI
 * 2: XML doc -> RDF/XML
 * 3: URI
 * 4: XML doc -> RDF/XML
 * ...
 * 
 * 
 * @author as
 *
 */
public class GeonamesDataReader implements DataReader {

	public static Logger log = Logger.getLogger(GeonamesDataReader.class);
	
	@Override
	public void loadData(RepositoryConnection conn, File file, Resource context) throws RepositoryException, IOException{
		
		BufferedReader bin = new BufferedReader(new FileReader(file));
		
		File errorFile = new File("geonames-error.rdf.txt");
		BufferedWriter errorDocs = new BufferedWriter(new FileWriter( errorFile ));
		int errorCount = 0;
		
		try {
			String line;
			while ( (line=bin.readLine())!=null) {
				
				String docId = line;
				
				// jump over it, the first line is just a URI
				line=bin.readLine();
				
				if (line==null || !line.startsWith("<?xml"))
					throw new IOException("Unexpected format. Line is expected to be a valid XML document. Line: \n" + line);
				
				StringReader nextDocument = new StringReader(line);
				try {
					conn.add(nextDocument, file.toURI().toString(), RDFFormat.RDFXML, context);
				} catch (RDFParseException e) {
					errorCount++;
					errorDocs.write(docId + "\r\n");
					errorDocs.write(line + "\r\n");
					errorDocs.flush();
				}
			}
		} finally {
			bin.close();
			errorDocs.flush();
			errorDocs.close();
		}
		
		if (errorCount>0)
			log.error("Not all documents could be added. Check " + errorFile.getName() + " for details.");
		else
			errorFile.delete();
	
	}

}
