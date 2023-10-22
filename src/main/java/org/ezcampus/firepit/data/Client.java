package org.ezcampus.firepit.data;

public class Client
{
	public String displayName;
	
	public String clientId;
	
	public Client(String clientId) {
		
		this.clientId = clientId;
		this.displayName = "Anonymous";
	}
	
	public Client(String clientId, String displayName) {
		
		this.clientId = clientId;
		this.displayName = displayName;
	}
}
