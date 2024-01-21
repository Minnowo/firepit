package websocket

import (
	"net/http"
	"sync"

	"github.com/EZCampusDevs/firepit/data"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

//
// RoomManager Code and Logic
//

type RoomList map[string]*Room

// The main RoomManager type
type RoomManager struct {
	rooms             RoomList
	roomCodeGenerator data.RoomCodeGenerator

	// To ensure threadsafe handling of rooms
	sync.RWMutex
}

// Creates a new room manager
func NewRoomManager() *RoomManager {
	return &RoomManager{
		rooms:             make(map[string]*Room),
		roomCodeGenerator: data.NewUintNRoomCodeGenerator(3, 16),
	}
}

// Get a unique room id
func (r *RoomManager) CreateRoomID() (string, error) {

	for {
		rid, err := r.roomCodeGenerator.GetRoomCode()

		if err != nil {
			return "", err
		}

		if r.HasRoom(rid) {
			log.Debugf("Room with id %d already exist!", rid)
			continue
		}

		r.AddRoom(rid)

		log.Debugf("Created room with id %d", rid)
		log.Debugf("There are now %d rooms", len(r.rooms))

		return rid, nil
	}
}

// Creates a new room with the given id
func (r *RoomManager) AddRoom(rid string) {

	room := NewRoom(rid, nil)

	// important we start the room loop here
	go room.RunRoom()

	r.Lock()
	defer r.Unlock()

	r.rooms[rid] = room
}

// Checks if there is an existing room with this id
func (r *RoomManager) HasRoom(rid string) bool {

	r.RLock()
	defer r.RUnlock()

	if _, ok := r.rooms[rid]; ok {
		return true
	}
	return false
}

// Adds the given client to the room
func (r *RoomManager) AddClientToRoom(rid string, c *Client) {

	r.RLock()
	defer r.RUnlock()

	if room, ok := r.rooms[rid]; ok {

		room.registerClient <- c
	}
}

// Removes the given client from the room
func (r *RoomManager) RemoveRoomClient(rid string, c *Client) {

	r.RLock()
	defer r.RUnlock()

	if room, ok := r.rooms[rid]; ok {

		room.unregisterClient <- c
	}
}

//
// Room Code and logic
//

// The room type used for json
type RoomJSON struct {
	ID                string         `json:"room_code"`
	Name              string         `json:"room_name"`
	Clients           ClientInfoList `json:"room_members"`
	Speaker           *ClientInfo    `json:"room_speaker"`
	Capacity          int            `json:"room_capacity"`
	RequireOccupation bool           `json:"room_occupation"`
}

// The main room type, used for logic
type Room struct {
	ID                string
	Name              string
	Clients           ClientSet
	Speaker           *Client
	Capacity          int
	RequireOccupation bool

	done             chan struct{}
	setSpeaker       chan *Client
	registerClient   chan *Client
	unregisterClient chan *Client
	broadcastInfo    chan *Client
	broadcast        chan *Event
}

// Create a new room
func NewRoom(name string, speaker *Client) *Room {
	return &Room{
		Name:             name,
		Speaker:          speaker,
		Capacity:         30,
		done:             make(chan struct{}),
		Clients:          make(ClientSet),
		setSpeaker:       make(chan *Client),
		registerClient:   make(chan *Client),
		unregisterClient: make(chan *Client),
		broadcast:        make(chan *Event),
		broadcastInfo:    make(chan *Client),
	}
}

// Handle channel communication for the room
func (r *Room) RunRoom() {

	for {
		select {
		case c := <-r.registerClient:
			log.Debugf("Adding client %s", c.info.Name)
			r._addClient(c)
		case c := <-r.unregisterClient:
			log.Debugf("Removing client %s", c.info.Name)
			r._removeClient(c)
		case c := <-r.setSpeaker:
			log.Debugf("Setting speaker %s", c.info.Name)
			r.Speaker = c
		case c := <-r.broadcastInfo:
			log.Debugf("Broadcast room info to %s", c.info.Name)
			r._broadCastRoomInfo(c)
		case e := <-r.broadcast:
			log.Debugf("Broadcasting %d", e.Type)
			r._broadcast(e)
		case <-r.done:
			log.Infof("Room %d is closing", r.ID)
			return
		}
		log.Debug("RunRoom tick")
	}
}

// Broadcast information about the room
func (r *Room) _broadCastRoomInfo(c *Client) {

	event, err := NewRoomInfoEvent(r)

	if err == nil {
		c.send <- *event
	} else {
		log.Error(err)
	}

}

// Broadcast that a client has joined the room
func (r *Room) _broadcastClientJoinedRoom(c *Client) {

	event, err := NewJoinRoomEvent(c)

	if err == nil {
		r._broadcast(event)
	} else {
		log.Error(err)
	}
}

// Broadcast that a client has left the room
func (r *Room) _broadcastClientLeaveRoom(c *Client) {

	event, err := NewLeaveRoomEvent(c)

	if err == nil {
		r._broadcast(event)
	} else {
		log.Error(err)
	}
}

// Broadcast that the speaker has been set
func (r *Room) _broadcastSetSpeaker() {

	if r.Speaker == nil {
		log.Warn("Cannot set speaker because speaker is nil")
		return
	}

	event, err := NewSetSpeakerEvent(r.Speaker)

	if err == nil {
		r._broadcast(event)
	} else {
		log.Error(err)
	}
}

// Adds the given client to the room
func (r *Room) _addClient(c *Client) {

	if r.Capacity == len(r.Clients) {

		log.Warnf("Client %s cannot join room %d because it is full", c.info.Name, r.ID)
		c.connection.Close()
		return
	}

	log.Debugf("Client joined room; Now has %d members", len(r.Clients))

	r._broadcastClientJoinedRoom(c)
	r.Clients[c] = true

	if len(r.Clients) == 1 {

		r.Speaker = c
	}

	r._broadCastRoomInfo(c)
}

// Remove the given client from the room
func (r *Room) _removeClient(c *Client) {

	if _, ok := r.Clients[c]; !ok {
		return
	}

	c.connection.Close()
	delete(r.Clients, c)

	r._broadcastClientLeaveRoom(c)

	if r.Speaker == c {

		for key := range r.Clients {
			r.Speaker = key
			r._broadcastSetSpeaker()
			break
		}
	}

	log.Debugf("Client left room; Now has %d members", len(r.Clients))
}

// Broadcast the given even to all room members
func (r *Room) _broadcast(e *Event) {

	for client := range r.Clients {

		log.Debugf("Broadcasting to client %d", client.info.Name)
		client.send <- *e
	}
}

// GET request handler for creating a new room
// Room ids are 8byte integers now
func (m *RoomManager) CreateRoomGET(c echo.Context) error {

	rid, err := m.CreateRoomID()

	if err != nil {
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Server error")
	}

	return c.String(http.StatusOK, rid)
}
