package websocket

import (
	"github.com/EZCampusDevs/firepit/util"
	"github.com/gorilla/websocket"
)

type ClientSet map[*Client]bool
type ClientList []*Client
type ClientInfoList []*ClientInfo

func (c ClientSet) ToSlice() ClientInfoList {

	keys := make([]*ClientInfo, len(c))

	i := 0
	for k := range c {
		keys[i] = k.info
		i++
	}

	return keys
}

type ClientInfo struct {
	Name       string `json:"client_name" query:"client_name"`
	DisplayId  string `json:"client_id" query:"client_id"`
	Occupation string `json:"client_occupation" query:"client_occupation"`
	RoomId     uint64 `query:"rid"`
}

func (c *ClientInfo) IsValid() bool {
	return !util.IsEmptyOrWhitespace(c.Name)
}

// Client is a websocket client, basically a frontend visitor
type Client struct {
	// the websocket connection
	connection *websocket.Conn

	// manager is the manager used to manage the client
	manager *Manager

	info *ClientInfo

	// send is used to avoid concurrent writes on the WebSocket
	send chan Event
}

// NewClient is used to initialize a new Client with all required values initialized
func NewClient(conn *websocket.Conn, manager *Manager, info *ClientInfo) *Client {
	return &Client{
		connection: conn,
		manager:    manager,
		info:       info,
		send:       make(chan Event),
	}
}
