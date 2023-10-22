package org.ezcampus.firepit.data.message;

import com.fasterxml.jackson.annotation.JsonProperty;

public class ClientPositionMessage extends SocketMessage
{
	@JsonProperty("new_name")
	public String newName;


	@Override
	public int getMessageType()
	{
		return CLIENT_SET_SPEAKER;
	}
}
