package org.ezcampus.firepit.websocket;


import java.io.IOException;

import org.ezcampus.firepit.data.Client;
import org.ezcampus.firepit.data.Room;
import org.ezcampus.firepit.data.RoomSessionController;
import org.tinylog.Logger;

import com.fasterxml.jackson.databind.ObjectMapper;

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
	
	private final ObjectMapper JSON_MAPPER = new ObjectMapper();
	
	@OnOpen
    public void onOpen(@PathParam("rid") String roomid, Session session) throws IOException {
		
		Room r = roomSessionController.getRoom(roomid);
		
		if(r == null) {
			
			session.close(new CloseReason(CloseCodes.UNEXPECTED_CONDITION, "Room does not exist"));
			
			return;
		}
			
		Client c = new Client(session.getId());
		
		r.addClient(c);
		
		Logger.info("New client connection {}",  session.getId());
    }

    @OnMessage
    public String onMessage(String name, Session session) {
    	
    	Logger.info("Client {} said {}", session.getId(), name);
    	
        return ("Hello" + name);
    }

    

    @OnClose
    public void helloOnClose(CloseReason reason, Session session) {
    	
    	Logger.info("Client {} has closed connection with reason {}", reason.getReasonPhrase());
    }
}