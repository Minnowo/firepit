package org.ezcampus.firepit.data;

import com.fasterxml.jackson.annotation.JsonIgnore;
import com.fasterxml.jackson.annotation.JsonProperty;

import jakarta.websocket.Session;

public class Client
{
	@JsonProperty("client_name")
	public String displayName;
	
	@JsonProperty("client_id")
	public String clientDisplayId;
	
	@JsonIgnore
	public String clientId;
	
	@JsonIgnore
	public Session clientSession;
	
	public Client(Session clientSession) {
		
		this.clientSession = clientSession;
		this.clientId = clientSession.getId();
		this.displayName = "Anonymous";
		this.clientDisplayId = java.util.UUID.randomUUID().toString();
	}
}
