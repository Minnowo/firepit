package org.ezcampus.firepit.websocket;


import java.io.IOException;

import org.ezcampus.firepit.data.Client;
import org.ezcampus.firepit.data.JsonService;
import org.ezcampus.firepit.data.Room;
import org.ezcampus.firepit.data.RoomSessionController;
import org.ezcampus.firepit.data.message.BadMessage;
import org.ezcampus.firepit.data.message.OkMessage;
import org.ezcampus.firepit.data.message.RenameMessage;
import org.ezcampus.firepit.data.message.SocketMessage;
import org.tinylog.Logger;

import jakarta.inject.Inject;
import jakarta.websocket.CloseReason;
import jakarta.websocket.CloseReason.CloseCodes;
import jakarta.websocket.OnClose;
import jakarta.websocket.OnMessage;
import jakarta.websocket.OnOpen;
import jakarta.websocket.Session;
import jakarta.websocket.server.PathParam;
import jakarta.websocket.server.ServerEndpoint;

@ServerEndpoint("/websocket/")
public class WSRoom {
	
	@Inject
	RoomSessionController roomSessionController;
	
	@Inject 
	JsonService jsonService;
	
	Client client;
	
	Room clientRoom;
	
	@OnOpen
    public void onOpen(@PathParam("rid") String roomid, Session session) throws IOException {
		
		clientRoom = roomSessionController.getRoom(roomid);
		
		if(clientRoom == null) {
			
			Logger.info("Someone tried to connect to a room which does not exist!");
			
			session.close(new CloseReason(CloseCodes.UNEXPECTED_CONDITION, "Room does not exist"));
			
			return;
		}

		client = new Client(session);
		
		clientRoom.addClient(client);
		
		Logger.info("New client connection {}",  session.getId());
		
		session.getBasicRemote().sendText(jsonService.toJson(new OkMessage("Room has been joined")));
    }

    @OnMessage
    public String onMessage(String json, Session session) throws IOException {
    	if(clientRoom == null) {
			
			session.close(new CloseReason(CloseCodes.UNEXPECTED_CONDITION, "Room does not exist"));
			
			return null;
		}
  
    	
    	SocketMessage m = jsonService.fromJson(json, SocketMessage.class);
    	
    	if(m == null) {
    		return jsonService.toJson(new BadMessage("Invalid message request"));
    	}
    	
    	this.clientRoom.broadCast(jsonService.toJson(new RenameMessage("this is a new name")), session);
    	
    	Logger.info("Client {} said {}", session.getId(), m.getMessageType());

        return jsonService.toJson(new OkMessage());
    }

    

    @OnClose
    public void onClose(CloseReason reason, Session session) {
    	
    	Logger.info("Client {} has closed connection with reason {}", session.getId(), reason.getReasonPhrase());
    	
    	if(clientRoom == null) {
			
			return;
		}
    	
    	clientRoom.removeClient(session.getId());
    }
}