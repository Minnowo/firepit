package org.ezcampus.firepit.api.models.request;

import com.fasterxml.jackson.annotation.JsonProperty;

public class CreateRoomQuery {
    @JsonProperty("room_name")
    private String roomName;
    
    @JsonProperty("room_capacity")
    private int roomCapacity;
    
    //*** Getters and Setters

    public String getRoomName() {
        return this.roomName;
    }

    public int getRoomCapacity() {
        return this.roomCapacity;
    }
}
