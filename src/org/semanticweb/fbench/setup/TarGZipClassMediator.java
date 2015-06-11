package org.semanticweb.fbench.setup;

import java.io.BufferedOutputStream;
import java.io.File;
import java.io.FileInputStream;
import java.io.FileOutputStream;
import java.io.IOException;
import java.util.zip.GZIPInputStream;

import org.apache.commons.compress.archivers.tar.TarArchiveEntry;
import org.apache.commons.compress.archivers.tar.TarArchiveInputStream;
import org.apache.log4j.Logger;

/**
 * Unzip (.tar.gz) a file and copy to specified destination.
 * 
 * @author as
 *
 */
public class TarGZipClassMediator implements ClassMediator {

	public static Logger log = Logger.getLogger(TarGZipClassMediator.class);
			
	@Override
	public void perform(File file, File dest) throws Exception {
	
		FileInputStream fin = new FileInputStream(file);
		
		GZIPInputStream gzip = new GZIPInputStream(fin);
		TarArchiveInputStream tis = new TarArchiveInputStream(gzip);
		dest.getParentFile().mkdirs();
		
		for (TarArchiveEntry entry = tis.getNextTarEntry(); entry != null;) {
			unpackEntries(tis, entry, dest.getParentFile());
			entry = tis.getNextTarEntry();
		}
		
		gzip.close();
		tis.close();
 	}
	
	private static void unpackEntries(TarArchiveInputStream tis, TarArchiveEntry entry, File outputDir) throws IOException {
		if (entry.isDirectory()) {
			File subDir = new File(outputDir, entry.getName());
			subDir.mkdirs();
			for (TarArchiveEntry e : entry.getDirectoryEntries()) {
				unpackEntries(tis, e, subDir);
			}
			return;
		}
		File outputFile = new File(outputDir, entry.getName());
		if (!outputFile.getParentFile().exists()) {
			outputFile.getParentFile().mkdirs();
		}
		
		log.info("Extracting file '" + entry.getName() + "' to '" + outputFile.getAbsolutePath() + "'");
		
		BufferedOutputStream bout = new BufferedOutputStream(new FileOutputStream(outputFile));
		
		int BUFFER = 32*1024;
		byte[] content = new byte[BUFFER];
	
		int read=0;
		while (tis.available() > 0) {
			read = tis.read(content, 0, Math.min(BUFFER, tis.available()));
			bout.write(content, 0, read);
		}
			
		
		bout.flush();
		bout.close();
	}
}
