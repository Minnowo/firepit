package org.ezcampus.firepit.data;

import java.io.IOException;
import java.util.List;
import java.util.concurrent.CopyOnWriteArrayList;

import org.ezcampus.firepit.data.message.SocketMessage;
import org.tinylog.Logger;

import jakarta.inject.Inject;
import jakarta.websocket.Session;

public class Room
{
	@Inject
	JsonService jsonService;

	private List<Client> roomMembers;

	public String roomName;

	public String roomId;

	public Room(String roomId)
	{

		this.roomMembers = new CopyOnWriteArrayList<Client>();

		this.roomId = roomId;

		this.roomName = roomId;
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
		if (!hasClient(c.clientId))
			this.roomMembers.add(c);
	}

	public void removeClient(String sessionId)
	{
		this.roomMembers.removeIf(x -> x.clientId.equals(sessionId));
	}

	public void broadCast(String message, Session sender)
	{
		Logger.debug("Broadcasting to room {} with {} members", this.roomId, this.roomMembers.size());

		for (Client c : this.roomMembers)
		{

			if (c.clientId.equals(sender.getId()))
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
