package org.ezcampus.firepit.data;

import java.util.ArrayList;

import javax.print.attribute.standard.Chromaticity;

import org.tinylog.Logger;

import jakarta.annotation.PostConstruct;
import jakarta.annotation.PreDestroy;
import jakarta.enterprise.context.ApplicationScoped;

@ApplicationScoped
public class RoomSessionController
{
	private int roomId;
	
	private ArrayList<Room> rooms;
	
	@PostConstruct
	void init()
	{
		roomId = 0;
		
		rooms = new ArrayList<Room>();
	}

	@PreDestroy
	void destroy()
	{
		// ...
	}

	public Room createRoom(String roomName, int roomCapacity, boolean requireOccupation) {
		
		Room r = new Room(Integer.toString(this.roomId), roomName, roomCapacity, requireOccupation);
		
		Logger.info("Creating new room wtih id {}", this.roomId);
		Logger.info(">> ROON NAME: {} ", roomName);
		Logger.info(">> ROON CAP: {} ", String.valueOf(roomCapacity));
			
		this.roomId++;
		
		this.rooms.add(r);
		
		return r;
	}
	
	public Room getRoom(String roomId) {
		
		if(roomId == null || roomId.isBlank())
			
			return null;
		
		for(Room r : this.rooms) 
			
			if(r.roomId.equals(roomId)) 
				
				return r;
			
		return null;		
	}
}
