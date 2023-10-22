package org.ezcampus.firepit.data.message;

import com.fasterxml.jackson.annotation.JsonProperty;

public class SetSpeakerMessage extends SocketMessage
{

	@JsonProperty("speaker_name")
	public String newSpeakerName;

	@Override
	public int getMessageType()
	{
		return CLIENT_SET_SPEAKER;
	}

	

}
