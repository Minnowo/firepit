package websocket

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/EZCampusDevs/firepit/util"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
	ErrEventNotSupported = errors.New("this event type is not supported")
)

// Manager is used to hold references to all Clients Registered, and Broadcasting etc
type Manager struct {
	roomManager    *RoomManager
	messageHandler map[int]EventHandler
}

// NewManager is used to initalize all the values inside the manager
func NewManager() *Manager {
	m := &Manager{
		roomManager:    NewRoomManager(),
		messageHandler: make(map[int]EventHandler),
	}
	m.setupEventHandlers()
	return m
}

func (m *Manager) GetRoomManager() *RoomManager {
	return m.roomManager
}

// setupEventHandlers configures and adds all handlers
func (m *Manager) setupEventHandlers() {
	m.messageHandler[SERVER_BAD_MESSAGE] = func(e Event, c *Client) error {
		log.Info("BAD MESSAGE: ", e)
		return nil
	}
	m.messageHandler[SERVER_OK_MESSAGE] = func(e Event, c *Client) error {
		log.Info("OK MESSAGE: ", e)
		return nil
	}
	m.messageHandler[CLIENT_SET_SPEAKER] = func(e Event, c *Client) error {
		log.Info("Set SPEAKER: ", e)
		return nil
	}
}

// routeEvent is used to make sure the correct event goes into the correct handler
func (m *Manager) routeEvent(event Event, c *Client) error {
	// Check if Handler is present in Map
	if handler, ok := m.messageHandler[event.Type]; ok {
		// Execute the handler and return any err
		if err := handler(event, c); err != nil {
			return err
		}
		return nil
	} else {
		return ErrEventNotSupported
	}
}

func (c *Client) readMessages() {

	defer func() {
		c.manager.roomManager.RemoveRoomClient(c.info.RoomId, c)
	}()

	log.Debug("Starting message sink for client ", c)

	for {
		// ReadMessage is used to read the next message in queue
		// in the connection
		_, payload, err := c.connection.ReadMessage()

		if err != nil {

			// If Connection is closed, we will Recieve an error here
			// We only want to log Strange errors, but not simple Disconnection
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Infof("error reading message: %v", err)
			}

			log.Debug("Exiting message sink for client ", c, " error: ", err)

			break
		}

		var request Event

		if err := json.Unmarshal(payload, &request); err != nil {
			log.Errorf("error marshalling message: %v", err)
			continue
		}

		// Route the Event
		if err := c.manager.routeEvent(request, c); err != nil {
			log.Error("Error handeling Message: ", err)
		}
	}
}

// writeMessages is a process that listens for new messages to output to the Client
func (c *Client) writeMessages() {

	var err error
	var data []byte

	defer func() {
		c.manager.roomManager.RemoveRoomClient(c.info.RoomId, c)
	}()

	for {
		select {
		case message, ok := <-c.send:

			// Ok will be false Incase the egress channel is closed
			if !ok {
				// Manager has closed this connection channel, so communicate that to frontend
				if err = c.connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Warn("connection closed: ", err)
				}
				// Return to close the goroutine
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

			log.Info("sent message")
		}

	}
}

func (m *Manager) ServeWebsocket(c echo.Context) error {

	var info ClientInfo

	if err := c.Bind(&info); err != nil {
		log.Debug("Client tried to connect with bad Info")
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
	}

	if util.IsEmptyOrWhitespace(info.Name) {
		log.Debug("Client tried to connect with bad Info")
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
	}

	if !m.roomManager.HasRoom(info.RoomId) {
		log.Warn("Client tried to join room which did not exist")
		return echo.NewHTTPError(http.StatusBadRequest, "No room exists")
	}

	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)

	if err != nil {
		log.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Server error!")
	}

	info.DisplayId = util.GetUUID()
	client := NewClient(ws, m, &info)

	go client.readMessages()
	go client.writeMessages()

	log.Info("Adding new client", client)
	log.Info("                 ", info)

	event, err := NewWhoAmIEvent(client)

	if err == nil {
		client.send <- *event
	}

	m.roomManager.BroadcastRoomInfo(info.RoomId, client)
	m.roomManager.AddClientToRoom(info.RoomId, client)

	return nil
}
