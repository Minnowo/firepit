package websocket

import (
	"net/http"
	"time"

	"github.com/EZCampusDevs/firepit/data"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

//
// Room Code and logic
//

var (
	RoomEmptyCheckInterval      = 2 * time.Minute
	RoomDisconnectClearInterval = 15 * time.Minute
)

type ReconnectionTokenMap map[string]*ClientInfo

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
	Clients           ClientMap
	Reconnects        ReconnectionTokenMap
	Speaker           *Client
	Capacity          uint32
	RequireOccupation bool
	ClientOrder       uint32
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
	const capacity uint32 = 25

	return &Room{
		Name:             name,
		Speaker:          speaker,
		Capacity:         capacity,
		ClientOrder:      0,
		Clients:          make(ClientMap),
		Reconnects:       make(ReconnectionTokenMap),
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
			log.Info(r.Clients)
			continue

		case c := <-r.unregisterClient:
			log.Debugf("Removing client %s", c.info.DisplayId)
			r._removeClient(c)
			log.Info(r.Clients)
			continue

		case c := <-r.setSpeakerById:
			log.Debugf("Setting speaker by id %s", c)
			r._setSpeakerById(c)
			continue

		case c := <-r.broadcastInfo:
			log.Debugf("Broadcast room info to %s", c.info.DisplayId)
			r._broadCastRoomInfo(c)
			continue

		case e := <-r.broadcast:
			log.Debugf("Broadcasting %d", e.Type)
			r._broadcast(e)
			continue

		case time := <-ticker.C:

			log.Debug("RunRoom tick")

			for _, i := range r.Reconnects {

				if time.Sub(i.DisconnectedAt) > RoomDisconnectClearInterval {

					delete(r.Reconnects, i.ReconnectionToken)
				}
			}

			if len(r.Clients) > 0 {

				r.LastEmptyTime = time

				log.Debugf("There are %s clients in the room!", len(r.Clients))

				continue
			}

			log.Debug("There are no clients in the room!")

			if time.Sub(r.LastEmptyTime) > RoomEmptyCheckInterval*2 {
				log.Infof("Room %s has been empty for a long time! It is now dead.", r.ID)
				return
			}
			continue

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
	}
}

// Called to cleanup a room that has died
func (r *Room) _cleanupRoom() {

	for client := range r.Clients {

		client.Disconnect()
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

	for c := range r.Clients {

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

	clientCount := uint32(len(r.Clients))

	if r.Capacity == clientCount {

		log.Warnf("Client %s cannot join room %d because it is full", c.info.DisplayId, r.ID)

		c.connection.Close()

		return
	}

	if c.status == STATUS_CLIENT_RECONNECT {

		log.Infof("Client %s is trying to reconnect with %s", c.info.DisplayId, c.info.ReconnectionToken)

		info, ok := r.Reconnects[c.info.ReconnectionToken]

		if !ok {

			log.Infof("Reconnection failed")

			c.Disconnect()

			return
		}

		delete(r.Reconnects, c.info.ReconnectionToken)

		log.Infof("RECONNECTION SUCCESSFUL")

		c.info = info
		c.status = STATUS_CLIENT_OK

	} else {

		r.ClientOrder++

		c.info.Number = r.ClientOrder
	}

	// tell other people c has joined
	r._broadcastClientJoinedRoom(c)

	// add c to the room
	r.Clients[c] = 0

	log.Infof("Adding client, current count is %d", clientCount)

	if clientCount == 0 {

		r.Speaker = c
		c.info.SpeakerRank = 1
	} else {

		c.info.SpeakerRank = 0
	}

	// send room info to c
	r._broadCastRoomInfo(c)
}

// Remove the given client from the room
func (r *Room) _removeClient(c *Client) {

	if c == nil {
		return
	}
	log.Infof("Client %s is being removed from the room", c.info.DisplayId)

	// delete from the room
	delete(r.Clients, c)

	c.info.DisconnectedAt = time.Now()

	// say this client can reconnect
	r.Reconnects[c.info.ReconnectionToken] = c.info

	// tell the client they're gone
	c.Disconnect()

	// tell everyone else they're gone
	r._broadcastClientLeaveRoom(c)

	if c == r.Speaker {

		// tell everyone who the new speaker is
		r._findNextBestSpeaker()
	}
}

func (r *Room) _findNextBestSpeaker() {

	if len(r.Clients) == 0 {

		r.Speaker = nil

		return
	}

	for c := range r.Clients {

		r.Speaker = c

		break
	}

	r._broadcastSetSpeaker()
}

// Broadcast the given even to all room members
func (r *Room) _broadcast(e *Event) {

	for c := range r.Clients {

		log.Debugf("Broadcasting to client %s", c.info.DisplayId)

		c.send <- *e
	}
}

// GET request handler for creating a new room
// Room ids are 8byte integers now
func (m *RoomManager) GETCreateRoom(c echo.Context) error {

	rid, err := m.CreateRoomID()

	if err != nil {
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Server error")
	}

	return c.String(http.StatusOK, rid)
}

func (m *RoomManager) GETHasRoom(c echo.Context) error {

	rid := c.Param("rid")

	return c.JSON(http.StatusOK, map[string]bool{"room_exists": m.HasRoom(rid)})
}
