package org.ezcampus.firepit.data;

import java.io.IOException;
import java.util.List;
import java.util.concurrent.CopyOnWriteArrayList;

import org.ezcampus.firepit.data.message.MessageType;
import org.ezcampus.firepit.data.message.SetSpeakerMessage;
import org.ezcampus.firepit.data.message.SocketMessage;
import org.tinylog.Logger;

import com.fasterxml.jackson.annotation.JsonProperty;

import jakarta.inject.Inject;
import jakarta.websocket.Session;

public class Room
{
	@JsonProperty("room_members")
	public List<Client> roomMembers;

	@JsonProperty("room_speaker")
	public Client speaker = null;
	
	@JsonProperty("room_name")
	public String roomName;

	@JsonProperty("room_capacity")
	public Integer roomCapacity;

	@JsonProperty("room_code")
	public String roomId;

	//* Do client's in this room need to specify their department/occupation */
	@JsonProperty("require_occupation")
	public Boolean requireOccupation;

	public Room(String roomId, String roomName, int roomCapacity, boolean requireOccupation)
	{

		this.roomMembers = new CopyOnWriteArrayList<Client>();

		this.roomId = roomId;

		// User Inputted Attributes for Room	
		this.roomName = roomName;
		this.roomCapacity = (Integer)roomCapacity;
		this.requireOccupation = (Boolean)requireOccupation;
	}

	public boolean hasClient(String sessionId)
	{
		for (Client c : this.roomMembers)
			if (c.clientId.equals(sessionId))
				return true;

		return false;
	}

	public void addClient(Client c)
	{
		if(this.roomMembers.size() == 0) {
			this.roomMembers.add(c);
			this.speaker = c;
		}
		else if (!hasClient(c.clientId)) {
			this.roomMembers.add(c);
		}
	}

	public void removeClient(String sessionId)
	{
		this.roomMembers.removeIf(x -> x.clientId.equals(sessionId));
	}
	
	
	public boolean setSpeakerFromPublicId(String publicId) {
		
		for(Client c : roomMembers) {
			
			if(c.clientDisplayId.equals(publicId)) {
				
				this.speaker = c;
				
				return true;
			}
		}
		
		return false;
	}
	

	public void broadCast(String message, Session sender)
	{
		Logger.debug("Broadcasting to room {} with {} members", this.roomId, this.roomMembers.size());

		for (Client c : this.roomMembers)
		{

			if (sender != null && c.clientId.equals(sender.getId()))
				continue;

			Logger.info("Sending message to client {}", c.clientId);

			try
			{
				if (c.clientSession.isOpen())
				{
					c.clientSession.getBasicRemote().sendText(message);
				}
				else
				{
					Logger.info("Cannot send message to client {} because the connection is closed", c.clientId);
				}

			}
			catch (IOException e)
			{
				Logger.error(e);
			}
		}
	}
}
