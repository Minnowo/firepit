package org.ezcampus.firepit.data.message;

import com.fasterxml.jackson.annotation.JsonProperty;

public class OkMessage extends SocketMessage
{
	@JsonProperty("reason")
	public String reason;
	
	public OkMessage()
	{
		this.reason = "";
	}

	public OkMessage(String reason)
	{
		this.reason = reason;
	}

	@Override
	public int getMessageType()
	{
		return SERVER_OK_MESSAGE;
	}

}
