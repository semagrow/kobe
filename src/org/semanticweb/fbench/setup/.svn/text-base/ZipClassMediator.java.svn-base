package org.semanticweb.fbench.setup;

import java.io.BufferedInputStream;
import java.io.BufferedOutputStream;
import java.io.File;
import java.io.FileInputStream;
import java.io.FileOutputStream;
import java.util.zip.ZipEntry;
import java.util.zip.ZipInputStream;

import org.apache.log4j.Logger;

/**
 * Unzip a file and copy to specified destination.
 * 
 * Assumes that dest is a directory!
 * 
 * @author as
 *
 */
public class ZipClassMediator implements ClassMediator {
	
	public static Logger log = Logger.getLogger(ZipClassMediator.class);
	
	@Override
	public void perform(File file, File dest) throws Exception {
		ZipInputStream zip = new ZipInputStream( new BufferedInputStream(new FileInputStream(file)));
		dest.getParentFile().mkdirs();

		int BUFFER = 32*1024;
		ZipEntry entry;
        while((entry = zip.getNextEntry()) != null) {
        	File out = new File(dest, entry.getName());
        	log.info("Extracting file '" + entry + "' to '" + out.getAbsolutePath() + "'");
        	out.getParentFile().mkdirs();
        	if (entry.isDirectory()) {
        		out.mkdir();
        		continue;
        	}       	
			
			int count;
			byte data[] = new byte[BUFFER];
			FileOutputStream fos = new FileOutputStream(out);
			BufferedOutputStream bout = new BufferedOutputStream(fos);
			while ((count = zip.read(data, 0, BUFFER)) != -1) {
				bout.write(data, 0, count);
			}
			bout.flush();
			bout.close();
        }
        zip.close();

	}
}
