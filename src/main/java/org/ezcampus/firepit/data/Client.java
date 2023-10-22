package org.ezcampus.firepit.data;

import jakarta.websocket.Session;

public class Client
{
	public String displayName;
	
	public String clientId;
	
	public Session clientSession;
	
	public Client(Session clientSession) {
		
		this.clientSession = clientSession;
		this.clientId = clientSession.getId();
		this.displayName = "Anonymous";
	}
	
	public Client(String clientId, String displayName) {
		
		this.clientId = clientId;
		this.displayName = displayName;
	}
}
