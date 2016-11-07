/**
 * 
 */
package gr.demokritos.iit.benchmark;

import com.fasterxml.jackson.databind.annotation.JsonDeserialize;

/**
 * @author Yiannis Mouchakis
 *
 */
@JsonDeserialize(using = ComposeBlockDeserializer.class)
public abstract class ComposeBlock {
	
	private String name;
	private String image;
	private String mount_dir;
	private String node;
		
	/**
	 * @param name
	 * @param image
	 * @param mount_dir
	 * @param node
	 */
	public ComposeBlock(String name, String image, String mount_dir, String node) {
		super();
		this.name = name;
		this.image = image;
		this.mount_dir = mount_dir;
		this.node = node;
	}

	public ComposeBlock() {
		super();
	}
	
	/**
	 * @return the name
	 */
	public String getName() {
		return name;
	}

	/**
	 * @param name the name to set
	 */
	public void setName(String name) {
		this.name = name;
	}

	/**
	 * @return the image
	 */
	public String getImage() {
		return image;
	}

	/**
	 * @param image the image to set
	 */
	public void setImage(String image) {
		this.image = image;
	}

	/**
	 * @return the mount_dir
	 */
	public String getMount_dir() {
		return mount_dir;
	}

	/**
	 * @param mount_dir the mount_dir to set
	 */
	public void setMount_dir(String mount_dir) {
		this.mount_dir = mount_dir;
	}

	/**
	 * @return the node
	 */
	public String getNode() {
		return node;
	}

	/**
	 * @param node the node to set
	 */
	public void setNode(String node) {
		this.node = node;
	}

	public abstract String generateBlock();
	

}
