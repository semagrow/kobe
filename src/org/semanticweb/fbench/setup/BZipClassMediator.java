package org.semanticweb.fbench.setup;

import java.io.BufferedInputStream;
import java.io.BufferedOutputStream;
import java.io.File;
import java.io.FileInputStream;
import java.io.FileOutputStream;

import org.apache.commons.compress.compressors.bzip2.BZip2CompressorInputStream;
import org.apache.log4j.Logger;

/**
 * Unzip (BZIP) a file and copy to specified destination.
 * 
 * @author as
 *
 */
public class BZipClassMediator implements ClassMediator {

	public static Logger log = Logger.getLogger(BZipClassMediator.class);
			
	@Override
	public void perform(File file, File dest) throws Exception {
	
		BufferedInputStream bin = new BufferedInputStream(new FileInputStream(file));
		
		BZip2CompressorInputStream zip = new BZip2CompressorInputStream(bin);
		dest.getParentFile().mkdirs();
		
		log.info("Extracting file '" + file + "' to '" + dest.getAbsolutePath() + "'");
		
		BufferedOutputStream bout = new BufferedOutputStream(new FileOutputStream(dest));
		
		byte[] buffer = new byte[32 * 1024];
		
		int n = 0;
		while (-1 != (n = zip.read(buffer))) {
		    bout.write(buffer, 0, n);
		}

		bout.flush();
		bout.close();		
	
		zip.close();
 	}
}
