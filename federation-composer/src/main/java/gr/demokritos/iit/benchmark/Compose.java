/**
 * 
 */
package gr.demokritos.iit.benchmark;

import java.io.File;
import java.io.FileWriter;
import java.io.IOException;
import java.util.ArrayList;
import java.util.LinkedList;
import java.util.List;

import com.fasterxml.jackson.core.JsonParseException;
import com.fasterxml.jackson.core.type.TypeReference;
import com.fasterxml.jackson.databind.JsonMappingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.dataformat.yaml.YAMLFactory;

/**
 * @author Yiannis Mouchakis
 *
 */
public class Compose {
	
	static final String compose_location = System.getProperty("user.dir") + "/docker-compose.yml";
	static final String config_location = System.getProperty("user.dir") + "/config.yml";
	
	List<ComposeBlock> compose = new LinkedList<>();

	public Compose() {
		super();
	}
	
	public boolean addBlock(ComposeBlock block) {
		return compose.add(block);
	}
	
	public boolean removeBlock(String name) {
		for (ComposeBlock block : compose) {
			if (block.getName().equals(name)) {
				return compose.remove(block);
			}
		}
		return false;
	}
	
	public List<String> list() {
		List<String> names = new ArrayList<>(compose.size());
		for (ComposeBlock block : compose) {
			names.add(block.getName());
		}
		return names;
	}
	
	public boolean contains(String name) {
		for (ComposeBlock block : compose) {
			if (name.equals(block.getName())) {
				return true;
			}
		}
		return false;
	}
	
	public void save() throws IOException {
				
		ObjectMapper mapper = new ObjectMapper(new YAMLFactory());
		mapper.writeValue(new File(config_location), compose);
		
		FileWriter writer = new FileWriter(compose_location);
		writer.write("version: '2'\n" + 
				"\n" + 
				"services:\n" + 
				"\n");//clears file
		
		for (ComposeBlock block : compose) {
			writer.append(block.generateBlock());
		}
		
		writer.close();
	}
	
	public void load() throws JsonParseException, JsonMappingException, IOException {
		
		load(config_location);
		
	}
	
	public void load(String file) throws JsonParseException, JsonMappingException, IOException {
		
		//check if file exists and if not create it and initialize it
		File config = new File(file);
		if ( ! config.exists() ) {
			FileWriter writer = new FileWriter(file);
			writer.write("[]");//initialize
			writer.close();
		}
		
		ObjectMapper mapper = new ObjectMapper(new YAMLFactory());
		compose = mapper.readValue(new File(file), new TypeReference<List<ComposeBlock>>(){});
		
	}
	
	public void clear() {
		this.compose.clear();
	}
	

}
