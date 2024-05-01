package websocket

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/EZCampusDevs/firepit/util"
	"github.com/gorilla/websocket"
	"github.com/labstack/gommon/log"
)

var (
	// pongWait is how long we will await a pong response from client
	pongWait = 10 * time.Second

	// pingInterval has to be less than pongWait, We cant multiply by 0.9 to get 90% of time
	// Because that can make decimals, so instead *9 / 10 to get 90%
	// The reason why it has to be less than PingRequency is becuase otherwise it will send a new Ping before getting response
	pingInterval = (pongWait * 9) / 10
)

type ClientList []*Client
type ClientInfoList []*ClientInfo
type ClientMap map[*Client]uint32
type ClientStatus uint32

const (
	STATUS_CLIENT_OK           ClientStatus = iota
	STATUS_CLIENT_DISCONNECTED ClientStatus = iota
	STATUS_CLIENT_RECONNECT    ClientStatus = iota
)

func (c *ClientMap) ToClientInfoSlice() ClientInfoList {

	keys := make([]*ClientInfo, len(*c))

	i := 0
	for c := range *c {

		keys[i] = c.info
		i++
	}

	keys = keys[0:i]

	return keys
}

// Used to convert the ClientSet into something json serializable
func (c *ClientList) ToClientInfoSlice() ClientInfoList {

	keys := make([]*ClientInfo, len(*c))

	i := 0
	for _, k := range *c {

		if k == nil || k.status != STATUS_CLIENT_OK {
			continue
		}

		keys[i] = k.info
		i++
	}

	keys = keys[0:i]

	return keys
}

// The main client info
type ClientInfo struct {
	Name              string    `json:"client_name" query:"name"`
	DisplayId         string    `json:"client_id" query:"id"`
	Occupation        string    `json:"client_occupation" query:"occup"`
	Number            uint32    `json:"order"`
	RoomId            string    `json:"-" query:"rid"`
	ReconnectionToken string    `json:"-" query:"rtoken"`
	SpeakerRank       uint32    `json:"-"`
	DisconnectedAt    time.Time `json:"-"`
}

// Determines if the info is valid to form a websocket connection
func (c *ClientInfo) IsValid() bool {
	return !util.IsEmptyOrWhitespace(c.Name) && !util.IsEmptyOrWhitespace(c.RoomId)
}

// Client is a websocket client, basically a frontend visitor
type Client struct {

	// the websocket connection
	connection *websocket.Conn

	// manager is the manager used to manage the client
	manager *Manager

	// information about the client
	info *ClientInfo

	// the room the client is in
	room *Room

	status ClientStatus

	// send is used to avoid concurrent writes on the WebSocket
	send chan Event

	once sync.Once
}

// Creates a new client
func NewClient(conn *websocket.Conn, manager *Manager, info *ClientInfo) *Client {
	return &Client{
		connection: conn,
		manager:    manager,
		info:       info,
		status:     STATUS_CLIENT_OK,
		send:       make(chan Event),
		once:       sync.Once{},
	}
}

func NewClientInRoom(conn *websocket.Conn, manager *Manager, info *ClientInfo, room *Room) *Client {
	return &Client{
		connection: conn,
		manager:    manager,
		info:       info,
		room:       room,
		status:     STATUS_CLIENT_OK,
		send:       make(chan Event),
		once:       sync.Once{},
	}
}

// Close the clients connection
func (c *Client) Disconnect() {

	c.connection.Close()
	c.room = nil
}

// Sends the WhoAmI event to the client
func (c *Client) BroadcastWhoAmI() {

	if c.status != STATUS_CLIENT_OK {

		return
	}

	event, err := NewWhoAmIEvent(c)

	if err == nil {
		c.send <- *event
	} else {
		log.Error(err)
	}
}

// Sends the pong message to the client
func (c *Client) pongHandler(_ string) error {
	return c.connection.SetReadDeadline(time.Now().Add(pongWait))
}

// The client message read routine
// Handles all messages from the client to the server
func (c *Client) readMessages() {

	// cleanup
	defer func() {

		c.once.Do(func() {

			if e := c.manager.roomManager.RemoveRoomClient(c.info.RoomId, c); e != nil {

				log.Error(e)
			}
		})
	}()

	// Configure Wait time for Pong response, use Current time + pongWait
	// This has to be done here to set the first initial timer.
	if err := c.pongHandler(""); err != nil {

		log.Error(err)

		return
	}

	// Configure how to handle Pong responses
	c.connection.SetPongHandler(c.pongHandler)

	for {
		_, payload, err := c.connection.ReadMessage()

		if c.room == nil {

			log.Errorf("Client %s does not have a room. Aborting connection", c.info.DisplayId)

			c.connection.Close()

			return
		}

		if err != nil {

			// If Connection is closed, we will Recieve an error here
			// We only want to log Strange errors, but not simple Disconnection
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {

				log.Errorf("error reading message: %v", err)
			}

			log.Debug("Exiting message sink for client ", c, " error: ", err)

			return
		}

		var request Event

		if err := json.Unmarshal(payload, &request); err != nil {

			log.Errorf("error marshalling message: %v", err)

			continue
		}

		// Route the Event and handle the client's message
		if err := c.manager.routeEvent(request, c); err != nil {

			log.Error("Error handeling Message: ", err)
		}
	}
}

// The client message write routine
// Handles all messages from the server to the client
func (c *Client) writeMessages() {

	var err error
	var data []byte

	// ping timer, sends out pings on this interval
	ticker := time.NewTicker(pingInterval)

	// cleanup
	defer func() {

		ticker.Stop()

		c.once.Do(func() {

			if e := c.manager.roomManager.RemoveRoomClient(c.info.RoomId, c); e != nil {

				log.Error(e)
			}
		})
	}()

	// for select channel pattern
	for {
		select {

		// when we get tick events ping the client
		case <-ticker.C:

			log.Debugf("Sending ping to %s", c.info.Name)

			if err := c.connection.WriteMessage(websocket.PingMessage, []byte{}); err != nil {

				log.Error("Cannot ping client: ", err)

				return
			}

		// when we get an event, send it to the client
		case message, ok := <-c.send:

			// channel has been closed
			if !ok {

				// Manager has closed this connection channel, so communicate that to frontend
				if err = c.connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Warn("connection closed: ", err)
				}

				return
			}

			data, err = json.Marshal(message)

			if err != nil {

				log.Error(err)

				return
			}

			if err = c.connection.WriteMessage(websocket.TextMessage, data); err != nil {

				log.Error(err)
			}

			log.Debug("sent message: ", message.Type)
		}

	}
}
