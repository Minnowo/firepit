package org.ezcampus.firepit.api;

import java.io.IOException;

import org.ezcampus.firepit.api.models.request.CreateRoomQuery;
import org.ezcampus.firepit.data.Room;
import org.ezcampus.firepit.data.RoomSessionController;

//* External Library Imports (Jackson, tinylog & Jakarta) */

import com.fasterxml.jackson.databind.ObjectMapper;
import org.tinylog.Logger;
import jakarta.inject.Inject;
import jakarta.ws.rs.Consumes;
import jakarta.ws.rs.GET;
import jakarta.ws.rs.POST;
import jakarta.ws.rs.Path;
import jakarta.ws.rs.Produces;
import jakarta.ws.rs.QueryParam;
import jakarta.ws.rs.core.MediaType;
import jakarta.ws.rs.core.Response;

@Path("/room")
public class APIRoom {
	@Inject
	RoomSessionController roomSessionController;

	private final ObjectMapper jsonMap = new ObjectMapper();

	@POST
	@Path("/new")
	@Consumes(MediaType.APPLICATION_JSON)
	@Produces(MediaType.APPLICATION_JSON)
	public Response createRoom(String JSON_PAYLOAD) {

		CreateRoomQuery requestData;
		try {
			requestData = jsonMap.readValue(JSON_PAYLOAD, CreateRoomQuery.class);
		} catch (IOException e) {
			Logger.debug("Got bad json: {}", e);
			return Response.status(Response.Status.BAD_REQUEST).entity("Invalid JSON payload").build();
		}

		//* --- Assertions for JSON where successful ---

		Room r = roomSessionController.createRoom(
			requestData.getRoomName(), 
			requestData.getRoomCapacity(),
			requestData.getRequireOccupation()
		);

		return Response.status(Response.Status.OK).entity(r.roomId).build();
	}

}
