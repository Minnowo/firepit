package org.ezcampus.firepit.data.message;

import org.ezcampus.firepit.data.Client;

import com.fasterxml.jackson.annotation.JsonProperty;

public class ClientJoinRoomMessage extends SocketMessage
{
	@JsonProperty("client")
	public Client client;

	public ClientJoinRoomMessage(Client c)
	{
		this.client = c;
	}


	@Override
	public int getMessageType()
	{
		return CLIENT_JOIN_ROOM;
	}

}
