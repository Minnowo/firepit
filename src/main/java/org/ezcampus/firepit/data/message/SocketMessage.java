package org.ezcampus.firepit.data.message;

import com.fasterxml.jackson.annotation.JsonIgnore;
import com.fasterxml.jackson.annotation.JsonSubTypes;
import com.fasterxml.jackson.annotation.JsonSubTypes.Type;
import com.fasterxml.jackson.annotation.JsonTypeInfo;
import com.fasterxml.jackson.annotation.JsonTypeInfo.Id;

@JsonTypeInfo(use = Id.NAME, include = JsonTypeInfo.As.PROPERTY, property = "messageType")
@JsonSubTypes({
    @Type(value = RenameMessage.class, name = "10"),
    @Type(value = ClientPositionMessage.class, name = "20"),
    @Type(value = SetSpeakerMessage.class, name = "30"),
    @Type(value = BadMessage.class, name = "400"),
    @Type(value = OkMessage.class, name = "200"),
})
public abstract class SocketMessage implements MessageType
{
	
	// this is put in automatically because of the above @Type
	@JsonIgnore
	public abstract int getMessageType(); 
}
