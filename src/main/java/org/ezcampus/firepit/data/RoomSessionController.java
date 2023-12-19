package org.ezcampus.firepit.data;

import java.util.ArrayList;
import java.util.Random;

import org.tinylog.Logger;
import jakarta.annotation.PostConstruct;
import jakarta.annotation.PreDestroy;
import jakarta.enterprise.context.ApplicationScoped;

@ApplicationScoped
public class RoomSessionController
{
	private ArrayList<Room> rooms;
	
	@PostConstruct
	void init()
	{
		rooms = new ArrayList<Room>();
	}

	@PreDestroy
	void destroy()
	{
		// ...
	}

	//* -------- EXTERNAL ROOM CODE GENERATION METHOD -------- */

	private static final String ALPHA_NUMERIC_STRING = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789";
    private Random random = new Random();

    private String generateRandomAlphaNumeric(int length) {
        StringBuilder builder = new StringBuilder();
        while (length-- != 0) {
            int character = (int)(random.nextDouble() * ALPHA_NUMERIC_STRING.length());
            builder.append(ALPHA_NUMERIC_STRING.charAt(character));
        }
        return builder.toString();
    }
	//* ------------------------------------------------------------  */

	public Room createRoom(String roomName, int roomCapacity, boolean requireOccupation) {
		
		String generatedRoomId = generateRandomAlphaNumeric(6);

		Room r = new Room(generatedRoomId, roomName, roomCapacity, requireOccupation);
		
		Logger.info("\nX-X =====| ROOM CREATION |===== X-X\n");
		Logger.info("\nINTERNAL ROOM ID >> {}", generatedRoomId);
		Logger.info("\n>> ROON NAME >> {} ", roomName);
		Logger.info("\n>> ROOM CAP >> {} ", String.valueOf(roomCapacity));
		Logger.info("\n>> REQUIRE OCCUPATION >> {} ", String.valueOf(requireOccupation));
			
		this.rooms.add(r);
		return r;
	}
	
	public Room getRoom(String roomId) {
		
		System.out.println("GOT ROOM_ID: "+roomId);

		if(roomId == null || roomId.isBlank())
			
			return null;
		
		for(Room r : this.rooms){
			
			System.out.println(">>> ROOM" + r.roomName + " "+ r.roomId);

			if(r.roomId.trim().equals(roomId.trim())) { // Avoiding white-space issues for java's strict comparison
				return r;
			} 
		}

		return null;		
	}
}
