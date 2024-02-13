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

// These are the event types which are valid
const (
	// used for a client to change their name
	EVENT__SET_CLIENT_NAME = 10

	// used for a client to change the speaker
	EVENT__CLIENT_SET_SPEAKER = 30

	// used when a client leaves the room
	EVENT__CLIENT_LEAVE_ROOM = 40

	// used when a client joins the room
	EVENT__CLIENT_JOIN_ROOM = 50

	// used to tell the client who they are
	EVENT__CLIENT_WHO_AM_I = 100

	// used to provide info about the joined room
	EVENT__ROOM_INFO = 60

	EVENT__SERVER_OK_MESSAGE  = 200
	EVENT__SERVER_BAD_MESSAGE = 400
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
	ClientInfo        *ClientInfo `json:"client"`
	ReconnectionToken string      `json:"reconnection_token"`
}

type RoomInfoEvent struct {
	Room *RoomJSON `json:"room"`
}

// Creates a new set speaker event
func NewSetSpeakerEvent(c *Client) (*Event, error) {

	var event SetSpeakerEvent

	event.SpeakerID = c.info.DisplayId

	jsonData, err := json.Marshal(&event)

	if err != nil {
		log.Error(err)
		return nil, err
	}

	return &Event{
		Type:    EVENT__CLIENT_SET_SPEAKER,
		Payload: jsonData,
	}, nil
}

// Creates a new room info event
func NewRoomInfoEvent(room *Room) (*Event, error) {

	var event RoomInfoEvent

	event.Room = &RoomJSON{
		ID:                room.ID,
		Name:              room.Name,
		Clients:           room.Clients.ToClientInfoSlice(),
		Capacity:          room.Capacity,
		RequireOccupation: room.RequireOccupation,
	}

	if room.Speaker != nil {
		event.Room.Speaker = room.Speaker.info
	}

	jsonData, err := json.Marshal(&event)

	if err != nil {
		log.Error(err)
		return nil, err
	}

	return &Event{
		Type:    EVENT__ROOM_INFO,
		Payload: jsonData,
	}, nil
}

// Creates a event with the given type which had only the data for a client
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

// Creates a new leave room event
func NewLeaveRoomEvent(c *Client) (*Event, error) {
	return NewCommonClientEvent(EVENT__CLIENT_LEAVE_ROOM, c)
}

// Creates a new join room event
func NewJoinRoomEvent(c *Client) (*Event, error) {
	return NewCommonClientEvent(EVENT__CLIENT_JOIN_ROOM, c)
}

// Creates a new who am i event
func NewWhoAmIEvent(c *Client) (*Event, error) {

	var event WhoAmIEvent

	event.ClientInfo = c.info
	event.ReconnectionToken = c.info.ReconnectionToken

	jsonData, err := json.Marshal(&event)

	if err != nil {
		log.Error(err)
		return nil, err
	}

	return &Event{
		Type:    EVENT__CLIENT_WHO_AM_I,
		Payload: jsonData,
	}, nil
}
