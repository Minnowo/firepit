package org.ezcampus.firepit.data.message;

import org.ezcampus.firepit.data.Room;

import com.fasterxml.jackson.annotation.JsonProperty;

public class RoomInfoMessage extends SocketMessage
{
	@JsonProperty("room")
	public Room room;

	public RoomInfoMessage(Room room)
	{
		this.room = room;
	}


	@Override
	public int getMessageType()
	{
		return ROOM_INFO;
	}

}
