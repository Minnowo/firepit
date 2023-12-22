package org.ezcampus.firepit.api;

import java.util.HashMap;

import org.tinylog.Logger;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;

import jakarta.ws.rs.GET;
import jakarta.ws.rs.Path;
import jakarta.ws.rs.Produces;
import jakarta.ws.rs.core.MediaType;
import jakarta.ws.rs.core.Response;

@Path("")
public class APICommon
{
	private static final ObjectMapper JSON_MAPPER = new ObjectMapper();
	
	final static String LICENSE_URL = "https://www.gnu.org/licenses/agpl-3.0.html";
	final static String SOURCE_URL = "https://github.com/EZCampusDevs/firepit";
	
	private static String HEARTBEAT_STRING = null;
	
	public Response heartBeat() 
	{
		try
		{
			if(HEARTBEAT_STRING == null) {
				HashMap<String, String> hm = new HashMap<String, String>();
				hm.put("detail", "EZCampus Firepit Backend");
				hm.put("license", LICENSE_URL);
				hm.put("source", SOURCE_URL);
				
				HEARTBEAT_STRING = JSON_MAPPER.writeValueAsString(hm);
			}
			
			return Response.status(Response.Status.OK).entity(HEARTBEAT_STRING).build();
		}
		catch (JsonProcessingException e)
		{
			Logger.error(e);
		}
		
		return Response.status(Response.Status.INTERNAL_SERVER_ERROR).build();
	}
	
	@GET
	@Produces(MediaType.APPLICATION_JSON)
	public Response heartBeat1() 
	{
		return heartBeat();
	}
	
	@GET
	@Path("/heartbeat")
	@Produces(MediaType.APPLICATION_JSON)
	public Response heartBeat2() 
	{
		return heartBeat();
	}
}
