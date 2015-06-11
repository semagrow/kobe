package org.semanticweb.fbench.provider;

import java.io.File;
import java.io.IOException;

import org.openrdf.model.Resource;
import org.openrdf.repository.RepositoryConnection;
import org.openrdf.repository.RepositoryException;

/**
 * DataReader implementations can be used to process special data formats, e.g. the
 * Geonames dump. See {@link NativeStoreFiller} for usage information.
 * 
 * @author as
 *
 */
public interface DataReader {

	/**
	 * load the data provided by file into the repository
	 * 
	 * @param conn
	 * 			the repository connections, remains open after processing
	 * @param file
	 * @param context
	 * @throws RepositoryException
	 * @throws IOException
	 */
	public void loadData(RepositoryConnection conn, File file, Resource context) throws RepositoryException, IOException;
}
