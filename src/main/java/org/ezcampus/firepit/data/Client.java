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

	@JsonProperty("client_occupation")
	public String clientOccupation;
	
	@JsonIgnore
	public String clientId;
	
	@JsonIgnore
	public Session clientSession;
	
	public Client(Session clientSession, String dirtyDisplayName, String dirtyOccupation) {

		//Display Name
		if(dirtyDisplayName != null && !dirtyDisplayName.isBlank()) {
			this.displayName = dirtyDisplayName;	
		} else {
			this.displayName = "Anonymous";
		}

		//Client Occupation
		this.clientOccupation = null;
		if(dirtyOccupation != null && !dirtyOccupation.isBlank()) {
			this.clientOccupation = dirtyOccupation;	
		} 

		//Generics
		this.clientSession = clientSession;
		this.clientId = clientSession.getId();
		this.clientDisplayId = java.util.UUID.randomUUID().toString();
	}
}
