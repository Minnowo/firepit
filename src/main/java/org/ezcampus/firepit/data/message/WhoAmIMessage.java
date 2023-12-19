package org.ezcampus.firepit.data.message;

import org.ezcampus.firepit.data.Client;

import com.fasterxml.jackson.annotation.JsonProperty;

public class WhoAmIMessage extends SocketMessage
{
	@JsonProperty("client")
	public Client client;

	public WhoAmIMessage(Client c)
	{
		this.client = c;
	}


	@Override
	public int getMessageType()
	{
		return CLIENT_WHO_AM_I_MESSAGE;
	}

}
