package websocket

import (
	"errors"
	"net/http"

	"github.com/EZCampusDevs/firepit/util"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

var (
	// configure the websockets
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}

	// error for bad events from client
	ErrEventNotSupported = errors.New("this event type is not supported")
)

// The main websocket manager
type Manager struct {
	roomManager    *RoomManager
	messageHandler map[int]EventHandler
}

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
	m.messageHandler[EVENT__CLIENT_SET_SPEAKER] = handleSetSpeaker
}

// Calls the correct handler for the given event
func (m *Manager) routeEvent(event Event, c *Client) error {

	handler, ok := m.messageHandler[event.Type]

	if !ok {
		return ErrEventNotSupported
	}

	return handler(event, c)
}

func (m *Manager) PrintDebugStuff(c echo.Context) error {

	var info ClientInfo

	if err := c.Bind(&info); err != nil {
		log.Debug("Client tried to connect with bad Info")
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
	}

	if !m.roomManager.HasRoom(info.RoomId) {
		log.Warn("Client tried to join room which did not exist")
		return echo.NewHTTPError(http.StatusBadRequest, "No room exists")
	}

	room, err := m.roomManager.GetRoomById(info.RoomId)

	if err != nil {
		log.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Server error!")
	}

	return c.String(200, room.String())
}

// Function for creating new websocket connections
// Will only accept connections with a valid name, and an existing room id
func (m *Manager) ServeWebsocket(c echo.Context) error {

	var info ClientInfo

	if err := c.Bind(&info); err != nil {
		log.Debug("Client tried to connect with bad Info")
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
	}

	if !info.IsValid() {
		log.Debug("Client tried to connect with bad Info")
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
	}

	if !m.roomManager.HasRoom(info.RoomId) {
		log.Warn("Client tried to join room which did not exist")
		return echo.NewHTTPError(http.StatusBadRequest, "No room exists")
	}

	// upgrade to websocket
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)

	if err != nil {
		log.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Server error!")
	}

	room, err := m.roomManager.GetRoomById(info.RoomId)

	if err != nil {
		log.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Server error!")
	}

	// create a public id for the user
	info.DisplayId = util.GetUUID()

	// create the client
	client := NewClientInRoom(ws, m, &info, room)

	if err := uuid.Validate(info.ReconnectionToken); err == nil {

		client.status = STATUS_CLIENT_RECONNECT

	} else {

		client.info.ReconnectionToken = util.GetUUID()
	}

	// handle the clients read and write
	go client.readMessages()
	go client.writeMessages()

	log.Debug("Adding new client", client.info.DisplayId)
	log.Debug("                 ", info)

	// add the client to their room
	client.room.registerClient <- client

	return nil
}
