/**
 * 
 */
package gr.demokritos.iit.benchmark;

import com.fasterxml.jackson.databind.annotation.JsonDeserialize;

/**
 * @author Yiannis Mouchakis
 *
 */
@JsonDeserialize(as = DataBlock.class)
public class DataBlock extends ComposeBlock {
	
	private String download_url;
	
	/**
	 * @param name
	 * @param image
	 * @param mount_dir
	 * @param node
	 * @param download_url
	 */
	public DataBlock(String name, String image, String mount_dir, String node, String download_url) {
		super(name, image, mount_dir, node);
		this.download_url = download_url;
	}
	
	public DataBlock() {
		super();
	}
	
	
	/**
	 * @return the download_url
	 */
	public String getDownload_url() {
		return download_url;
	}


	/**
	 * @param download_url the download_url to set
	 */
	public void setDownload_url(String download_url) {
		this.download_url = download_url;
	}


	@Override
	public String generateBlock() {
		
		StringBuilder builder = new StringBuilder();
		
		builder.append("    " + getName() + ":" + System.lineSeparator());
		builder.append("        image: " + getImage() + System.lineSeparator());
		builder.append("        container_name: " + getName() + System.lineSeparator());
		builder.append("        volumes_from:" + System.lineSeparator());
		builder.append("            - " + getName() + "-data" + System.lineSeparator());
		if (getDownload_url() != null && ! getDownload_url().equals("") ) {
			builder.append("        environment:" + System.lineSeparator());
			builder.append("            - DOWNLOAD_URL=" + download_url + System.lineSeparator());
		}
		builder.append("    " + getName() + "-data:" + System.lineSeparator());
		builder.append("        image: busybox" + System.lineSeparator());
		if (getNode() != null && ! getNode().equals("")) {
			builder.append("        environment:" + System.lineSeparator());
			builder.append("            - \"constraint:node==" + getNode() + "\"" + System.lineSeparator());
		}
		builder.append("        volumes:" + System.lineSeparator());
		builder.append("            - /data" + System.lineSeparator());
		if (getMount_dir() != null && ! getMount_dir().equals("")) {
			builder.append("            - " + getMount_dir() + ":/data/toLoad");
		}
		builder.append(System.lineSeparator());
		
		return builder.toString();
	}
	
}