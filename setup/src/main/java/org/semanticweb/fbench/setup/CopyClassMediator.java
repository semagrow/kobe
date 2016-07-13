package org.semanticweb.fbench.setup;

import java.io.File;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.semanticweb.fbench.misc.FileUtil;

/**
 * Copies a file to specified destination.
 * 
 * @author as
 *
 */
public class CopyClassMediator implements ClassMediator {

	public static Logger log = LoggerFactory.getLogger(CopyClassMediator.class);
	

	public void perform(File file, File dest) throws Exception {
		log.info("Copying file '" + file + "' to '" + dest.getAbsolutePath() + "'");
		FileUtil.copyFile(file, dest);		
	}
	
	

}
