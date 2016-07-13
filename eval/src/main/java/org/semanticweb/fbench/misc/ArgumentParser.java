package org.semanticweb.fbench.misc;

import java.util.ArrayList;
import java.util.List;

import org.semanticweb.fbench.Config.Property;



/**
 * Utility class for command line argument parsing.
 * 
 * @author as
 *
 */
public class ArgumentParser {

	
	/**
	 * allowed arguments (in any order)
	 *  "configFile" -> location of the config file
	 *  "-fill" -> enable filling, i.e. do not execute queries
	 *  "-setup" -> enable setup mode
	 *  
	 * @param args
	 * @return
	 */
	public static List<Property> parseArguments(String args[]) {
		if (args==null || args.length==0)
			return new ArrayList<Property>();
		
		ArrayList<Property> res = new ArrayList<Property>();
		for (int i=0; i<args.length; i++) {
			if (args[i].equals("-fill"))
				res.add(new Property("fill", "true"));
			else if (args[i].equals("-setup"))
				res.add(new Property("setup", "true"));
			else
				res.add(new Property("configFile", args[i]));
		}
		
		return res;
	}
}
