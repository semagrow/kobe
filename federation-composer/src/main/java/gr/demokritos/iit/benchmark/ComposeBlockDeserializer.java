/**
 * 
 */
package gr.demokritos.iit.benchmark;

import java.io.IOException;

import com.fasterxml.jackson.core.JsonParser;
import com.fasterxml.jackson.databind.DeserializationContext;
import com.fasterxml.jackson.databind.JsonDeserializer;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.node.ObjectNode;

/**
 * @author Yiannis Mouchakis
 *
 */
public class ComposeBlockDeserializer extends JsonDeserializer<ComposeBlock> {

    @Override
    public ComposeBlock deserialize(JsonParser jp, DeserializationContext context) throws IOException {

        ObjectMapper mapper = (ObjectMapper) jp.getCodec();
        ObjectNode root = mapper.readTree(jp);
        if (root.has("federation_type")) {
        	return mapper.readValue(root.toString(), FederationBlock.class);
        } else {
        	return mapper.readValue(root.toString(), DataBlock.class);
        }
    }
}