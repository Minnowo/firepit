package org.ezcampus.firepit.api;

import org.ezcampus.firepit.data.Room;
import org.ezcampus.firepit.data.RoomSessionController;

import jakarta.inject.Inject;
import jakarta.ws.rs.GET;
import jakarta.ws.rs.Path;
import jakarta.ws.rs.Produces;
import jakarta.ws.rs.core.MediaType;
import jakarta.ws.rs.core.Response;

@Path("/room")
public class APIRoom
{
	@Inject
	RoomSessionController roomSessionController;
	
	@GET
	@Path("/new")
	@Produces(MediaType.TEXT_PLAIN)
	public Response createRoom() {
		
		Room r = roomSessionController.createRoom();
		
		return Response.status(Response.Status.OK).entity(r.roomId).build();
	}

}
