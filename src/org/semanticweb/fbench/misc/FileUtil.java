package org.semanticweb.fbench.misc;

import java.io.BufferedInputStream;
import java.io.BufferedOutputStream;
import java.io.File;
import java.io.FileInputStream;
import java.io.FileNotFoundException;
import java.io.FileOutputStream;
import java.io.IOException;
import java.net.HttpURLConnection;
import java.net.URL;
import java.net.URLConnection;

import org.apache.log4j.Logger;
import org.semanticweb.fbench.Config;



/**
 * Utility class for file operations.
 * 
 * @author as
 *
 */
public class FileUtil {

	public static Logger log = Logger.getLogger(FileUtil.class);
	
	/**
	 * location utility.<p>
	 * 
	 *  if the specified path is absolute, it is returned as is, 
	 *  otherwise a location relative to {@link Config#getBaseDir()} is returned<p>
	 *  
	 *  examples:<p>
	 *  
	 *  <code>
	 *  /home/data/myPath -> absolute linux path
	 *  c:\\data -> absolute windows path
	 *  \\\\myserver\data -> absolute windows network path (see {@link File#isAbsolute()})
	 *  data/myPath -> relative path (relative location to baseDir is returned)
	 *  </code>
	 *  
	 * @param path
	 * @return
	 * 			the file corresponding to the abstract path
	 */
	public static File getFileLocation(String path) {
		
		// check if path is an absolute path that already exists
		File f = new File(path);
		
		if (f.isAbsolute())
			return f;
		
		f = new File( Config.getConfig().getBaseDir() + path);
		return f;
	}
	
	
	
    public static void copyDir(File source, File target) throws FileNotFoundException, IOException {
        
        File[] files = source.listFiles();
        target.mkdirs();
        for (File file : files) {
            if (file.isDirectory()) {
                copyDir(file, new File(target.getAbsolutePath() + System.getProperty("file.separator") + file.getName()));
            }
            else {
                copyFile(file, new File(target.getAbsolutePath() + System.getProperty("file.separator") + file.getName()));
            }
        }
    }
    
    public static void copyFile(File file, File target) throws FileNotFoundException, IOException {
        
        BufferedInputStream in = new BufferedInputStream(new FileInputStream(file));
        
        if (target.exists())
        	log.warn("WARNING: file " + target.getAbsolutePath() + " will be overridden.");
      
        if (!target.getParentFile().exists()) {
        	target.getParentFile().mkdirs();
        }
        
        BufferedOutputStream out = new BufferedOutputStream(new FileOutputStream(target));
        int bytes = 0;
        while ((bytes = in.read()) != -1) {
            out.write(bytes);
        }
        in.close();
        out.close();
    } 
    
    public static void download(URL url, File dest) throws Exception {
		URLConnection con = url.openConnection();
		if (con instanceof HttpURLConnection && ((HttpURLConnection)con).getResponseCode()!=200)
			throw new RuntimeException("Download of " + url + " failed with error '" + ((HttpURLConnection)con).getResponseMessage() + "'");
		
		log.debug("Download of '" + url + "' to '" + dest + "'");
		BufferedInputStream in = new BufferedInputStream(con.getInputStream());
		FileOutputStream out = new FileOutputStream(dest);
		int i = 0;
		byte[] bytesIn = new byte[1024];
		while ((i = in.read(bytesIn)) >= 0) {
			out.write(bytesIn, 0, i);
		}
		out.close();
		in.close();
	}
}
