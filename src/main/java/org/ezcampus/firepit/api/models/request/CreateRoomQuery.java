package org.ezcampus.firepit.api.models.request;

import com.fasterxml.jackson.annotation.JsonProperty;

public class CreateRoomQuery {
    @JsonProperty("room_name")
    private String roomName;
    
    @JsonProperty("room_capacity")
    private int roomCapacity;

    @JsonProperty("require_occupation")
    private boolean requireOccupation;
    
    //*** Getters and Setters

    public String getRoomName() {
        return this.roomName;
    }

    public int getRoomCapacity() {
        return this.roomCapacity;
    }

    public boolean getRequireOccupation() {
        return this.requireOccupation;
    }
}
