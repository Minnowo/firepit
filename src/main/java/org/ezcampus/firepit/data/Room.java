package org.ezcampus.firepit.data;

import java.util.HashMap;

public class Room
{
	private HashMap<String, Client> roomMembers;
	
	public String roomName;
	
	public String roomId;
	
	public Room(String roomId) {
		
		this.roomMembers = new HashMap<String, Client>();
		
		this.roomId = roomId;
		
		this.roomName = roomId;
	}
	
	public boolean hasClient(String sessionId) { 
		return roomMembers.containsKey(sessionId);
	}
	
	public void addClient(Client c) { 
		this.roomMembers.put(c.clientId, c);
	}
	
}
