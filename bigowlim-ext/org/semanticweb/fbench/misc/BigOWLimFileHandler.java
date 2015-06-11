package org.semanticweb.fbench.misc;

import java.io.File;


/**
 * Utility tool for BigOWLim repository location.<p>
 * 
 * There is a strange bug in BigOWLim's OwlimSchemaRepository implementation
 * which prevent the usage of absolute paths for repository locations, i.e. 
 * the parameter "storage-folder" cannot be used with absolute paths.<p>
 * 
 * A work-around is to adjust the base directory and use only the last folder
 * name as a relative path.<p>
 * 
 * This class implements a way to obtain the base directory and the repoLocation's
 * name from a given repo location.
 * 
 * @author as
 *
 */
public class BigOWLimFileHandler {

	
	private File baseDir;
	private String repoLocationName;
	
	public BigOWLimFileHandler(String repoLocation) {
		File f = new File(repoLocation);
		baseDir = f.getParentFile();
		repoLocationName = f.getName();
	}
	
	public BigOWLimFileHandler(File repoLocation) {
		baseDir = repoLocation.getParentFile();
		repoLocationName = repoLocation.getName();
	}
	
	public File getBaseDir() {
		return baseDir;
	}

	public String getRepoLocationName() {
		return repoLocationName;
	}

	
	
	public static void main(String[] args) {
		
		test("C:\\\\data\\owlimStorage", new File("C:\\\\data"), "owlimStorage");
		test("C:\\\\data\\owlimStorage\\", new File("C:\\\\data"), "owlimStorage");
		
	}
	
	protected static void test(String loc, File expectedDir, String expectedName) {
		BigOWLimFileHandler test = new BigOWLimFileHandler(loc);
		System.out.println("*** Testing " + loc + " ***");
		System.out.println("BaseDir: " + test.getBaseDir() + " (" + expectedDir + ")");
		System.out.println("Name: " + test.getRepoLocationName() + "(" +expectedName + ")");
		if (!test.getBaseDir().equals(expectedDir) || !test.getRepoLocationName().equals(expectedName) )
			throw new RuntimeException("Test failed!");
	}
}
