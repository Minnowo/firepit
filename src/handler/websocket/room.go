package websocket

import (
	"fmt"
	"net/http"
	"time"

	"github.com/EZCampusDevs/firepit/data"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
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
		ID:               name,
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

func (r *Room) String() string {

	speaker := "None"

	if r.Speaker != nil {
		speaker = fmt.Sprintf("%s (ID: %s)", r.Speaker.info.Name, r.Speaker.info.DisplayId)
	}

	clients := ""
	for client := range r.Clients {
		clients += fmt.Sprintf("\n    - %s (ID: %s)", client.info.Name, client.info.DisplayId)
	}

	reconnects := ""

	for token, client := range r.Reconnects {
		reconnects += fmt.Sprintf("\n    - Token: %s, ClientID: %s", token, client.DisplayId)
	}

	return fmt.Sprintf(`Room:
  ID: %s
  Name: %s
  Clients: %s
  Reconnects: %s
  Speaker: %s
  Capacity: %d
  RequireOccupation: %t
  ClientOrder: %d
  LastEmptyTime: %s`,
		r.ID,
		r.Name,
		clients,
		reconnects,
		speaker,
		r.Capacity,
		r.RequireOccupation,
		r.ClientOrder,
		r.LastEmptyTime.Format(time.RFC1123))
}

// Handle channel communication for the room
func (r *Room) RunRoom() {

	ticker := time.NewTicker(RoomEmptyCheckInterval)

	defer ticker.Stop()

	for {
		select {

		case c := <-r.registerClient:

			log.Info().Str("client", c.info.DisplayId).Str("room", r.ID).Msg("Adding client")

			r._addClient(c)

			continue

		case c := <-r.unregisterClient:

			log.Info().Str("client", c.info.DisplayId).Str("room", r.ID).Msg("Removing client")

			r._removeClient(c)

			continue

		case c := <-r.setSpeakerById:

			log.Info().Str("id", c).Str("room", r.ID).Msg("Setting speaker")

			r._setSpeakerById(c)

			continue

		case c := <-r.broadcastInfo:

			log.Info().Str("id", c.info.DisplayId).Str("room", r.ID).Msg("Broadcast room info to client")

			r._broadCastRoomInfo(c)

			continue

		case e := <-r.broadcast:

			log.Info().Int("type", e.Type).Str("room", r.ID).Msg("Broadcasting message")

			r._broadcast(e)

			continue

		case time := <-ticker.C:

			log.Debug().Msg("RunRoom tick")

			for _, i := range r.Reconnects {

				if time.Sub(i.DisconnectedAt) > RoomDisconnectClearInterval {

					log.Info().Str("client", i.DisplayId).Msg("Removing client from the reconnects list")

					delete(r.Reconnects, i.ReconnectionToken)
				}
			}

			if len(r.Clients) > 0 {

				r.LastEmptyTime = time

				log.Info().Str("room", r.ID).Int("clientCount", len(r.Clients)).Msg("Room still has clients")

				continue
			}

			log.Info().Str("room", r.ID).Msg("Room is empty")

			if time.Sub(r.LastEmptyTime) > RoomEmptyCheckInterval*2 {

				log.Info().Str("room", r.ID).Msg("Room has been empty for a long time, killing it")

				return
			}

			continue

		// lets us handle pause, resume, and kill the thread
		case s := <-r.state:

			log.Info().Str("room", r.ID).Str("state", data.ChannelStateToString(s)).Msg("Room state has been changed")

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

				log.Info().Str("room", r.ID).Str("state", data.ChannelStateToString(s)).Msg("Room state has been changed")

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
		log.Error().Err(err)
	}

}

// Broadcast that a client has joined the room
func (r *Room) _broadcastClientJoinedRoom(c *Client) {

	event, err := NewJoinRoomEvent(c)

	if err == nil {
		r._broadcast(event)
	} else {
		log.Error().Err(err)
	}
}

// Broadcast that a client has left the room
func (r *Room) _broadcastClientLeaveRoom(c *Client) {

	event, err := NewLeaveRoomEvent(c)

	if err == nil {
		r._broadcast(event)
	} else {
		log.Error().Err(err)
	}
}

// Broadcast that the speaker has been set
func (r *Room) _broadcastSetSpeaker() {

	if r.Speaker == nil {

		log.Warn().Str("room", r.ID).Msg("Cannot set speaker because speaker is nil")

		return
	}

	event, err := NewCommonClientEvent(EVENT__CLIENT_SET_SPEAKER, r.Speaker)

	if err == nil {
		r._broadcast(event)
	} else {
		log.Error().Err(err)
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

		log.Info().Str("room", r.ID).Str("client", c.info.DisplayId).Msg("Found new speaker for room")

		r.Speaker = c

		break
	}

	r._broadcastSetSpeaker()
}

// Adds the given client to the room
func (r *Room) _addClient(c *Client) {

	clientCount := uint32(len(r.Clients))

	if r.Capacity == clientCount {

		log.Warn().Str("room", r.ID).Str("client", c.info.DisplayId).Msg("Client cannot join room because it is full")

		c.connection.Close()

		return
	}

	if c.status == STATUS_CLIENT_RECONNECT {

		log.Info().Str("room", r.ID).Str("client", c.info.DisplayId).Msg("Client is trying to reconnect")

		info, ok := r.Reconnects[c.info.ReconnectionToken]

		if !ok {

			log.Warn().Str("room", r.ID).Str("client", c.info.DisplayId).Msg("Client failed to reconnect")

			c.Disconnect()

			return
		}

		delete(r.Reconnects, c.info.ReconnectionToken)

		log.Info().Str("room", r.ID).Str("client", c.info.DisplayId).Msg("Client has reconnected")

		c.info = info
		c.status = STATUS_CLIENT_OK

	} else {

		r.ClientOrder++

		c.info.Number = r.ClientOrder
	}

	// inform the client who they are
	c.BroadcastWhoAmI()

	// tell other people c has joined
	r._broadcastClientJoinedRoom(c)

	// add c to the room
	r.Clients[c] = 0

	log.Info().Str("room", r.ID).Str("client", c.info.DisplayId).Msg("Client has joined the room")

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

	log.Info().Str("room", r.ID).Str("client", c.info.DisplayId).Msg("Client is being removed from the room")

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

		log.Debug().Str("room", r.ID).Str("client", c.info.DisplayId).Int("type", e.Type).Msg("Broadcasting message")

		c.send <- *e
	}
}

// GET request handler for creating a new room
// Room ids are 8byte integers now
func (m *RoomManager) GETCreateRoom(c echo.Context) error {

	rid, err := m.CreateRoomID()

	if err != nil {

		log.Error().Err(err).Msg("Could not create room")

		return c.String(http.StatusInternalServerError, "Server error")
	}

	return c.String(http.StatusOK, rid)
}

func (m *RoomManager) GETHasRoom(c echo.Context) error {

	rid := c.Param("rid")

	return c.JSON(http.StatusOK, map[string]bool{"room_exists": m.HasRoom(rid)})
}
