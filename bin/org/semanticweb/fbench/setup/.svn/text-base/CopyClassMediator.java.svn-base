package org.semanticweb.fbench.setup;

import java.io.File;

import org.apache.log4j.Logger;
import org.semanticweb.fbench.misc.FileUtil;

/**
 * Copies a file to specified destination.
 * 
 * @author as
 *
 */
public class CopyClassMediator implements ClassMediator {

	public static Logger log = Logger.getLogger(CopyClassMediator.class);
	
	@Override
	public void perform(File file, File dest) throws Exception {
		log.info("Copying file '" + file + "' to '" + dest.getAbsolutePath() + "'");
		FileUtil.copyFile(file, dest);		
	}
	
	

}
