package websocket

import (
	"encoding/json"

	"github.com/labstack/gommon/log"
)

// Event is the Messages sent over the websocket
// Used to differ between different actions
type Event struct {
	// Type is the message type sent
	Type int `json:"messageType"`
	// Payload is the data Based on the Type
	Payload json.RawMessage `json:"payload"`
}

// EventHandler is a function signature that is used to affect messages on the socket and triggered
// depending on the type
type EventHandler func(event Event, c *Client) error

const (
	SET_CLIENT_NAME         = 10
	SET_CLIENT_POSITION     = 20
	CLIENT_SET_SPEAKER      = 30
	CLIENT_LEAVE_ROOM       = 40
	CLIENT_JOIN_ROOM        = 50
	CLIENT_WHO_AM_I_MESSAGE = 100

	ROOM_INFO = 60

	SERVER_OK_MESSAGE  = 200
	SERVER_BAD_MESSAGE = 400
)

type BadMessageEvent struct {
	Reason string `json:"reason"`
}

type OkMessageEvent struct {
	Reason string `json:"reason"`
}

type SetSpeakerEvent struct {
	SpeakerID string `json:"speaker_id"`
}

type JoinRoomEvent struct {
	ClientInfo *ClientInfo `json:"client"`
}

type LeaveRoomEvent struct {
	ClientInfo *ClientInfo `json:"client"`
}

type WhoAmIEvent struct {
	ClientInfo *ClientInfo `json:"client"`
}

type RoomInfoEvent struct {
	Room *RoomJSON `json:"room"`
}

func NewRoomInfoEvent(room *Room) (*Event, error) {

	var event RoomInfoEvent

	event.Room = &RoomJSON{
		ID:                room.ID,
		Name:              room.Name,
		Clients:           room.Clients.ToSlice(),
		Capacity:          room.Capacity,
		RequireOccupation: room.RequireOccupation,
	}

	if room.Speaker != nil {
		event.Room.Speaker = room.Speaker.info
	}

	log.Error(event.Room.Clients)

	jsonData, err := json.Marshal(&event)

	if err != nil {
		log.Error(err)
		return nil, err
	}

	return &Event{
		Type:    ROOM_INFO,
		Payload: jsonData,
	}, nil
}
func NewCommonClientEvent(type_ int, c *Client) (*Event, error) {

	var event JoinRoomEvent

	event.ClientInfo = c.info

	jsonData, err := json.Marshal(&event)

	if err != nil {
		log.Error(err)
		return nil, err
	}

	return &Event{
		Type:    type_,
		Payload: jsonData,
	}, nil
}

func NewLeaveRoomEvent(c *Client) (*Event, error) {
	return NewCommonClientEvent(CLIENT_LEAVE_ROOM, c)
}
func NewJoinRoomEvent(c *Client) (*Event, error) {
	return NewCommonClientEvent(CLIENT_JOIN_ROOM, c)
}
func NewWhoAmIEvent(c *Client) (*Event, error) {
	return NewCommonClientEvent(CLIENT_WHO_AM_I_MESSAGE, c)
}
