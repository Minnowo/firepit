package org.ezcampus.firepit.data.message;

import com.fasterxml.jackson.annotation.JsonProperty;

public class RenameMessage extends SocketMessage
{
	@JsonProperty("new_name")
	public String newName;

	public RenameMessage(String newname)
	{
		this.newName = newname;
	}

	@Override
	public int getMessageType()
	{
		return SET_CLIENT_NAME;
	}

}
