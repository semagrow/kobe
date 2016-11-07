/**
 * 
 */
package gr.demokritos.iit.benchmark;

import java.io.Console;
import java.io.FileReader;
import java.io.IOException;
import java.net.URISyntaxException;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.ArrayList;
import java.util.List;
import java.util.Properties;

import org.apache.commons.cli.CommandLine;
import org.apache.commons.cli.CommandLineParser;
import org.apache.commons.cli.DefaultParser;
import org.apache.commons.cli.HelpFormatter;
import org.apache.commons.cli.Option;
import org.apache.commons.cli.OptionGroup;
import org.apache.commons.cli.Options;
import org.apache.commons.cli.ParseException;

import com.fasterxml.jackson.core.JsonParseException;
import com.fasterxml.jackson.databind.JsonMappingException;

/**
 * @author Yiannis Mouchakis
 *
 */
public class Main {

	private static Compose compose = new Compose();
	
	/**
	 * @param args
	 * @throws JsonMappingException 
	 * @throws JsonParseException 
	 * @throws IOException 
	 * @throws URISyntaxException 
	 */
	public static void main(String[] args) throws JsonParseException, JsonMappingException, IOException, URISyntaxException  {
		
		OptionGroup group = new OptionGroup();
		group.setRequired(true);
		
		Option list = Option.builder("l")
				.desc("list current containers")
				.longOpt("list")
				.build();
		group.addOption(list);
		
		Option clear = Option.builder("c")
				.desc("clear current containers")
				.longOpt("clear")
				.build();
		group.addOption(clear);
		
		Option remove = Option.builder("r")
				.desc("remove a container")
				.longOpt("remove")
				.hasArgs()
				.build();
		group.addOption(remove);
		
		Option file = Option.builder("f")
				.desc("load containers from file")
				.longOpt("file")
				.hasArg()
				.build();
		group.addOption(file);
		
		Option add = Option.builder("a")
				.desc("add containers")
				.longOpt("add")
				.build();
		group.addOption(add);
		
		Option help = Option.builder("h")
				.desc("prints this help message and exits")
				.longOpt("help")
				.build();
		group.addOption(help);
		
		Options options = new Options();
		options.addOptionGroup(group);
		
		CommandLineParser parser = new DefaultParser();
				
		try {
			
			CommandLine cmd = parser.parse( options, args);
			
			compose.load();
			
			if (cmd.hasOption(help.getOpt()) || cmd.getOptions().length == 0) {
				
				printHelp(options);
				
			} 
			else if (cmd.hasOption(clear.getOpt())) {
				
				compose.clear();
				
			}
			else if (cmd.hasOption(remove.getOpt())) {
				
				for (String name : cmd.getOptionValues(remove.getOpt())) {
					boolean removed = compose.removeBlock(name);
					if (removed) {
						System.out.println("removed \"" + name + "\"");
					} else {
						System.err.println("container \""+name+"\" does not exist!");
					}
				}				
			}
			else if (cmd.hasOption(file.getOpt())) {
				
				compose.load(cmd.getOptionValue(file.getOpt()));
				
			} 
			else if (cmd.hasOption(list.getOpt())) {
				
				for (String name : compose.list()) {
					System.out.println(name);
				}
				
			}
			else if (cmd.hasOption(add.getOpt())) {
				
				addContainer();
				
			}
			
		
			compose.save();
			
			
		} catch (ParseException e) {
			System.err.println("Parsing failed.  Reason: " + e.getMessage());
			printHelp(options);
		}

	}
	
	private static void addContainer() throws IOException, URISyntaxException {
		
		Console c = System.console();

		List<String> opts = new ArrayList<>();
		String data_opt = "1";
		opts.add(data_opt);
		String semagrow_opt = "2";
		opts.add(semagrow_opt);
		String fedx_opt = "3";
		opts.add(fedx_opt);
		String opf_opt = "4";
		opts.add(opf_opt);
		String cd_opt = "5";
		opts.add(cd_opt);
		String lsd_opt = "6";
		opts.add(lsd_opt);
		
		String line = null;
		while (true) {
			
			line = c.readLine("Select option ("+data_opt+") add single data node, ("+semagrow_opt+") add semagrow, ("+fedx_opt+") add fedx, "
					+ "("+opf_opt+") add all OPF data nodes, ("+cd_opt+") add all Cross Domain data nodes, ("+lsd_opt+") add all Life Science Domain data nodes: ").trim();
			
			if (opts.contains(line)) {
				break;
			} else {
				System.err.println("Invalid option " + line);
			}
			
		}
		
		Properties properties = new Properties();
		properties.load(new FileReader(System.getProperty("user.dir") + "/default_images"));
		
		if (line.equals(data_opt)) {
			
			insertDataNode(c, properties.getProperty("default_data_image"));
			
		} else if (line.equals(semagrow_opt)) {
			
			insertFederationNode(c, properties.getProperty("default_semagrow_image"), "semagrow");
			
		} else if  (line.equals(fedx_opt)) {
			
			insertFederationNode(c, properties.getProperty("default_fedx_image"), "fedx");
			
		} else if  (line.equals(opf_opt)) {
			
			createPredifinedData(properties.getProperty("default_data_image"), System.getProperty("user.dir") + "/opf_datasets");
			
		} else if  (line.equals(cd_opt)) {
			
			createPredifinedData(properties.getProperty("default_data_image"), System.getProperty("user.dir") + "/cross_domain_datasets");
			
		} else if  (line.equals(lsd_opt)) {
			
			createPredifinedData(properties.getProperty("default_data_image"), System.getProperty("user.dir") + "/life_science_domain_datasets");
			
		}
		
	}
	
	public static void insertDataNode(Console c, String default_image) {
		
		System.out.println("Insert container details. Leave empty for default or optional.");
		
		//name must be set and be unique in compose
		String name = null;
		while (true) {
			name = c.readLine("Insert container name: ").trim();
			if (name.equals("")) {
				System.err.println("You must set the name!");
			} else if (compose.contains(name)) {
				System.err.println("Name \""+name+"\" is already in compose! Please select another name.");
			} else {
				break;
			}
		}
		String image = c.readLine("Insert container image (default \""+default_image+"\"): ").trim();
		if (image.equals("")) {
			image = default_image;
		}
		String mount_dir = c.readLine("Insert the directory that contains the data to be loaded (optional): ").trim();
		String node = c.readLine("Insert node to deploy to (optional): ").trim();
		String download_url = c.readLine("Insert URL to downolad the dataset to load (optional): ").trim();
		
		DataBlock block = new DataBlock(name, image, mount_dir, node, download_url);
		
		compose.addBlock(block);
		
	}
	
	public static void insertFederationNode(Console c, String default_image, String federation_type) {
		
		System.out.println("Insert federation details. Leave empty for default or optional.");
		
		//name must be set and be unique in compose
		String name = null;
		while (true) {
			name = c.readLine("Insert container name: ").trim();
			if (name.equals("")) {
				System.err.println("You must set the name!");
			} else if (compose.contains(name)) {
				System.err.println("Name \""+name+"\" is already in compose! Please select another name.");
			} else {
				break;
			}
		}
		
		String image = c.readLine("Insert federation image (default \""+default_image+"\"): ").trim();
		if (image.equals("")) {
			image = default_image;
		}
		
		//must set the mount directory
		String mount_dir = null;
		while (true) {
			mount_dir = c.readLine("Insert the directory that contains the federation configuration: ").trim();
			if (mount_dir.equals("")) {
				System.err.println("You must set the directory!");
			} else {
				break;
			}
		}
				
		String node = c.readLine("Insert node to deploy to (optional): ").trim();
		
		FederationBlock block = new FederationBlock(name, image, mount_dir, node, federation_type);
		
		compose.addBlock(block);
		
	}
	
	private static void printHelp(Options options) {
		HelpFormatter formater = new HelpFormatter();
		formater.printHelp("federation-composer.sh", options);
	}
	
	private static void createPredifinedData(String default_img, String filename) throws IOException, URISyntaxException {
		
		Path path = Paths.get(filename);
		for (String line : Files.readAllLines(path, StandardCharsets.UTF_8)) {
			String[] splited = line.split(",");
			if (compose.contains(splited[0])) {
				System.err.println("Container with name \""+splited[0]+"\" already exists in compose and will not be added.");
			} else {
				DataBlock block = new DataBlock(splited[0], default_img, "", "", splited[1]);
				compose.addBlock(block);
			}
		}
	}

}
