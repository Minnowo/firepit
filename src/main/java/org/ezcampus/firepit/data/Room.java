package org.ezcampus.firepit.data;

import java.io.IOException;
import java.util.List;
import java.util.Stack;
import java.util.concurrent.ConcurrentLinkedDeque;
import java.util.concurrent.CopyOnWriteArrayList;

import org.ezcampus.firepit.data.message.SetSpeakerMessage;
import org.tinylog.Logger;

import com.fasterxml.jackson.annotation.JsonIgnore;
import com.fasterxml.jackson.annotation.JsonProperty;

import jakarta.inject.Inject;
import jakarta.websocket.Session;

public class Room
{
	@JsonIgnore
	public ConcurrentLinkedDeque<Client> speakerHistory;
	
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
		this.speakerHistory = new ConcurrentLinkedDeque<Client>();
		

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
			this.speakerHistory.add(c);
		}
		else if (!hasClient(c.clientId)) {
			this.roomMembers.add(c);
		}
	}

	public void removeClient( Session sender )
	{
		this.roomMembers.removeIf(x -> x.clientId.equals(sender.getId()));
		
		if(this.roomMembers.size() == 0) {
			this.speaker = null;
			return;
		}
		
		if(this.speaker == null || !this.speaker.clientId.equals(sender.getId())) {
			return;
		}
		
		Logger.info("The current speaker has left the room!");
		
		
		// remove the speaker from the history
		this.speakerHistory.pop();
		
		// find the previous speaker and make them the new speaker
		while(!this.speakerHistory.isEmpty()) {
			
			Client c = this.speakerHistory.pop();
			
			if(c.clientId.equals(sender.getId())) {
				continue;
			}
			
			if(!this.broadCastSetSpeaker(c.clientDisplayId, sender)) {
				continue;
			}
			
			return;
		}
		
		Logger.info("Could not find previous speaker!");
		Logger.info("Trying the first person in the room: {}", this.roomMembers.get(0).clientId);
		
		this.broadCastSetSpeaker(this.roomMembers.get(0).clientDisplayId, sender);
	}
	
	
	public boolean setSpeakerFromPublicId(String publicId) {
		
		for(Client c : roomMembers) {
			
			if(c.clientDisplayId.equals(publicId)) {
				
				this.speaker = c;
				this.speakerHistory.add(c);
				Logger.info("Speaker has been set to {}", c.clientId);
				return true;
			}
		}
		
		return false;
	}
	
	public boolean broadCastSetSpeaker(String publicId, Session sender) {
		
		if(!this.setSpeakerFromPublicId(publicId)) {
			return false;
		}
		
		SetSpeakerMessage ssm = new SetSpeakerMessage();
		
		ssm.newSpeakerId = publicId;
		
		this.broadCast(ssm.toJson(), sender);
		
		return true;
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
