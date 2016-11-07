/**
 * 
 */
package gr.demokritos.iit.benchmark;

import com.fasterxml.jackson.databind.annotation.JsonDeserialize;

/**
 * @author Yiannis Mouchakis
 *
 */
@JsonDeserialize(as = FederationBlock.class)
public class FederationBlock extends ComposeBlock {
	
	private String federation_type;
	
	/**
	 * @param name
	 * @param image
	 * @param mount_dir
	 * @param node
	 * @param federation_type
	 */
	public FederationBlock(String name, String image, String mount_dir, String node, String federation_type) {
		super(name, image, mount_dir, node);
		this.federation_type = federation_type;
	}
	
	public FederationBlock() {
		super();
	}

	/**
	 * @return the federation_type
	 */
	public String getFederation_type() {
		return federation_type;
	}

	/**
	 * @param federation_type the federation_type to set
	 */
	public void setFederation_type(String federation_type) {
		this.federation_type = federation_type;
	}

	@Override
	public String generateBlock() {
		
		StringBuilder builder = new StringBuilder();
		
		builder.append("    " + getName() + ":" + System.lineSeparator());
		builder.append("        image: " + getImage() + System.lineSeparator());
		builder.append("        container_name: " + getName() + System.lineSeparator());
		if (getNode() != null && ! getNode().equals("") ) {
			builder.append("        environment:" + System.lineSeparator());
			builder.append("            - \"constraint:node==" + getNode() + "\"" + System.lineSeparator());
		}
		builder.append("        volumes:" + System.lineSeparator());
		if (getFederation_type().equals("semagrow")) {
			builder.append("            - " + getMount_dir() + ":/etc/default/semagrow");
		} else if (getFederation_type().equals("fedx")) {
			builder.append("            - " + getMount_dir() + ":/etc/fedx");
		}
		
		builder.append(System.lineSeparator());
		
		return builder.toString();
	}

	

}
