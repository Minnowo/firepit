package websocket

import (
	"fmt"
	"net/http"
	"sync"
	"time"

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

// Gets a room by id
func (r *RoomManager) GetRoomById(rid string) (*Room, error) {

	r.RLock()
	defer r.RUnlock()

	if room, ok := r.rooms[rid]; ok {

		return room, nil
	}

	return nil, fmt.Errorf("Room does not exist")
}

// Adds the given client to the room
func (r *RoomManager) AddClientToRoom(rid string, c *Client) error {

	r.RLock()
	defer r.RUnlock()

	if room, ok := r.rooms[rid]; ok {

		room.registerClient <- c
		return nil
	}
	return fmt.Errorf("Room does not exist")
}

// Removes the given client from the room
func (r *RoomManager) RemoveRoomClient(rid string, c *Client) error {

	r.RLock()
	defer r.RUnlock()

	if room, ok := r.rooms[rid]; ok {

		room.unregisterClient <- c
		return nil
	}
	return fmt.Errorf("Room does not exist")
}

// Sets the clients room object
func (r *RoomManager) SetClientRoomPtr(rid string, c *Client) error {

	r.RLock()
	defer r.RUnlock()

	if room, ok := r.rooms[rid]; ok {

		c.room = room
		return nil
	}

	return fmt.Errorf("Room does not exist")
}

//
// Room Code and logic
//

var (
	RoomEmptyCheckInterval = 2 * time.Minute
)

// The room type used for json
type RoomJSON struct {
	ID                string         `json:"room_code"`
	Name              string         `json:"room_name"`
	Clients           ClientInfoList `json:"room_members"`
	Speaker           *ClientInfo    `json:"room_speaker"`
	Capacity          uint32         `json:"room_capacity"`
	RequireOccupation bool           `json:"room_occupation"`
}

// The main room type, used for logic
type Room struct {
	ID                string
	Name              string
	Size              uint32
	Clients           ClientList
	Speaker           *Client
	Capacity          uint32
	RequireOccupation bool
	LastEmptyTime     time.Time

	state            chan byte
	setSpeakerById   chan string
	registerClient   chan *Client
	unregisterClient chan *Client
	broadcastInfo    chan *Client
	broadcast        chan *Event
}

// Create a new room
func NewRoom(name string, speaker *Client) *Room {
	const capacity uint32 = 30

	return &Room{
		Name:             name,
		Speaker:          speaker,
		Capacity:         capacity,
		Clients:          make(ClientList, capacity),
		state:            make(chan byte),
		setSpeakerById:   make(chan string),
		registerClient:   make(chan *Client),
		unregisterClient: make(chan *Client),
		broadcast:        make(chan *Event),
		broadcastInfo:    make(chan *Client),
	}
}

// Handle channel communication for the room
func (r *Room) RunRoom() {

	ticker := time.NewTicker(RoomEmptyCheckInterval)

	defer ticker.Stop()

	for {
		select {
		case c := <-r.registerClient:
			log.Debugf("Adding client %s", c.info.DisplayId)
			r._addClient(c)
		case c := <-r.unregisterClient:
			log.Debugf("Removing client %s", c.info.DisplayId)
			r._removeClient(c)
		case c := <-r.setSpeakerById:
			log.Debugf("Setting speaker by id %s", c)
			r._setSpeakerById(c)
		case c := <-r.broadcastInfo:
			log.Debugf("Broadcast room info to %s", c.info.DisplayId)
			r._broadCastRoomInfo(c)
		case e := <-r.broadcast:
			log.Debugf("Broadcasting %d", e.Type)
			r._broadcast(e)

		case time := <-ticker.C:

			if len(r.Clients) > 0 {

				r.LastEmptyTime = time

				log.Debug("There are clients in the room!")

				continue
			}

			log.Debug("There are no clients in the room!")

			if time.Sub(r.LastEmptyTime) > RoomEmptyCheckInterval*2 {
				log.Infof("Room %s has been empty for a long time! It is now dead.", r.ID)
				return
			}

		// lets us handle pause, resume, and kill the thread
		case s := <-r.state:

			log.Debugf("Room state is %s", data.ChannelStateToString(s))

			switch s {
			default:
			case data.CHAN__RUNNING:
				continue
			case data.CHAN__PAUSED:
				break
			case data.CHAN__DEAD:
				r._cleanupRoom()
				return
			}

		lock:
			for {

				s = <-r.state

				log.Debugf("Room state is %s", data.ChannelStateToString(s))

				switch s {
				case data.CHAN__RUNNING:
					break lock
				case data.CHAN__PAUSED:
					break
				case data.CHAN__DEAD:
					r._cleanupRoom()
					return
				}
			}
		}

		log.Debug("RunRoom tick")
	}
}

// Called to cleanup a room that has died
func (r *Room) _cleanupRoom() {

	for _, client := range r.Clients {

		if client != nil {
			client.Disconnect()
		}
	}

	r.Speaker = nil
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

// Set the speaker by id
func (r *Room) _setSpeakerById(id string) {

	for _, c := range r.Clients {

		if c == nil || c.status != STATUS_CLIENT_OK {
			continue
		}

		if c.info.DisplayId != id {
			continue
		}

		if r.Speaker != nil {
			c.info.SpeakerRank = r.Speaker.info.SpeakerRank + 1
		} else {
			c.info.SpeakerRank = 1
		}

		log.Debugf("Found new speaker %s", c.info.DisplayId)

		r.Speaker = c

		break
	}

	r._broadcastSetSpeaker()
}

// Adds the given client to the room
func (r *Room) _addClient(c *Client) {

	if r.Capacity == r.Size {

		log.Warnf("Client %s cannot join room %d because it is full", c.info.DisplayId, r.ID)

		c.connection.Close()

		return
	}

	if c.status == STATUS_CLIENT_RECONNECT {

		log.Infof("Client %s is trying to reconnect!", c.info.DisplayId)

		if r._tryReconnectClient(c) {

			log.Infof("Client %s has reconnected!", c.info.DisplayId)

			r._broadcastClientJoinedRoom(c)

			if r.Size == 0 {
				r.Speaker = c
			}

			r.Size++

			r._broadCastRoomInfo(c)

			return
		}

		log.Errorf("Client %s failed to reconnect!", c.info.DisplayId)

		c.Disconnect()

		return
	}

	r._broadcastClientJoinedRoom(c)

	if r.Size == 0 {

		r.Speaker = c
		c.info.SpeakerRank = 1
	} else {

		c.info.SpeakerRank = 0
	}

	// if all clients in the array are not null
	// we will replace a disconnected client
	var replaceDisconnectedIndex int = -1

	// find a new spot for the client
	for i, client := range r.Clients {

		if client == nil {

			c.info.Number = uint32(i)

			r.Clients[i] = c

			r.Size += 1

			replaceDisconnectedIndex = -1

			break
		}

		if replaceDisconnectedIndex == -1 && client.status != STATUS_CLIENT_OK {

			replaceDisconnectedIndex = i
		}
	}

	// will be -1 if there is no non-null client to replace
	if replaceDisconnectedIndex != -1 {

		c.info.Number = uint32(replaceDisconnectedIndex)

		r.Clients[replaceDisconnectedIndex] = c

		r.Size += 1
	}

	log.Debugf("Client joined room; Can reconnect with %d; Now has %d members", c.info.ReconnectionToken, r.Size)
	log.Info(r.Clients)

	r._broadCastRoomInfo(c)
}

// try and reconnect the client
// if a disconnected client has the same ReconnectionToken it will be replaced with the given client
func (r *Room) _tryReconnectClient(c *Client) bool {

	for i, client := range r.Clients {

		if client == nil || client.status == STATUS_CLIENT_OK {
			continue
		}

		if client.info.ReconnectionToken == c.info.ReconnectionToken {

			c.status = STATUS_CLIENT_OK

			c.info = client.info

			client.info = nil

			r.Clients[i] = c

			return true
		}
	}

	return false
}

// Remove the given client from the room
func (r *Room) _removeClient(c *Client) {

	if c == nil {
		return
	}

	log.Infof("Client %s is being removed from the room", c.info.DisplayId)

	for i := uint32(0); i < r.Size; i++ {

		cl := r.Clients[i]

		if cl != c {

			continue
		}

		switch c.status {
		case STATUS_CLIENT_LEFT:
		case STATUS_CLIENT_REMOVED_BY_SERVER:
			r.Clients[i] = nil
			break

		case STATUS_CLIENT_DISCONNECTED:
			break

		case STATUS_CLIENT_OK:

			log.Warnf("Speaker %s is leaving with ok stsatus", c.info.DisplayId)

			c.status = STATUS_CLIENT_DISCONNECTED
			break
		}

		log.Infof("Client %s has been disconnected", c.info.DisplayId)

		c.Disconnect()

		r._broadcastClientLeaveRoom(c)

		r.Size--

		log.Debugf("Client left room; Now has %d members", r.Size)

		break
	}

	if c == r.Speaker {

		r._findNextBestSpeaker()
	}
}

func (r *Room) _findNextBestSpeaker() {

	if r.Size == 0 {

		r.Speaker = nil
		return
	}

	for _, cl := range r.Clients {

		if cl == nil || cl.status != STATUS_CLIENT_OK {
			continue
		}

		if r.Speaker.status != STATUS_CLIENT_OK {

			r.Speaker = cl
			continue
		}

		if cl.info.SpeakerRank > r.Speaker.info.SpeakerRank {

			r.Speaker = cl
		}
	}

	r._broadcastSetSpeaker()
}

// Broadcast the given even to all room members
func (r *Room) _broadcast(e *Event) {

	for i := uint32(0); i < r.Size; i++ {

		client := r.Clients[i]

		if client == nil || client.status != STATUS_CLIENT_OK {
			continue
		}

		log.Debugf("Broadcasting to client %s", client.info.DisplayId)

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
