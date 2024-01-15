package websocket

import (
	"net/http"
	"strconv"
	"sync"

	"github.com/EZCampusDevs/firepit/util"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type RoomList map[uint64]*Room

type RoomManager struct {
	rooms RoomList

	sync.RWMutex
}

func NewRoomManager() *RoomManager {
	return &RoomManager{
		rooms: make(map[uint64]*Room),
	}
}

func (r *RoomManager) HasRoom(rid uint64) bool {

	r.RLock()
	defer r.RUnlock()

	if _, ok := r.rooms[rid]; ok {
		return true
	}
	return false
}

func (r *RoomManager) AddRoom(rid uint64) {

	r.Lock()
	defer r.Unlock()

	r.rooms[rid] = NewRoom(strconv.FormatUint(rid, 10), nil)
}

func (r *RoomManager) AddClientToRoom(rid uint64, c *Client) {

	var err error
	var event *Event

	r.RLock()
	defer r.RUnlock()

	event, err = NewJoinRoomEvent(c)

	if err == nil {
		r.lockedBroadcast(rid, event)
	}

	if room, ok := r.rooms[rid]; ok {
		room.AddClient(c)
	}
}

func (r *RoomManager) RemoveRoomClient(rid uint64, c *Client) {

	r.RLock()
	defer r.RUnlock()

	event, err := NewLeaveRoomEvent(c)

	if err == nil {
		r.lockedBroadcast(rid, event)
	}

	if room, ok := r.rooms[rid]; ok {
		room.RemoveClient(c)
	}
}

func (r *RoomManager) Broadcast(rid uint64, e *Event) {

	r.RLock()
	defer r.RUnlock()

	r.lockedBroadcast(rid, e)
}

func (r *RoomManager) lockedBroadcast(rid uint64, e *Event) {

	if room, ok := r.rooms[rid]; ok {
		room.Broadcast(e)
	}
}

func (r *RoomManager) BroadcastRoomInfo(rid uint64, c *Client) {

	r.RLock()
	defer r.RUnlock()

	if room, ok := r.rooms[rid]; ok {

		rinfo, err := room.GetRoomInfoEvent()

		if err == nil {
			c.send <- *rinfo
		}
	}
}

func (r *RoomManager) CreateRoom() (uint64, error) {

	for {
		rid, err := util.GenerateFull64BitNumber()

		if err != nil {
			return 0, err
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

type RoomJSON struct {
	ID                uint64         `json:"room_code"`
	Name              string         `json:"room_name"`
	Clients           ClientInfoList `json:"room_members"`
	Speaker           *ClientInfo    `json:"room_speaker"`
	Capacity          uint           `json:"room_capacity"`
	RequireOccupation bool           `json:"room_occupation"`
}
type Room struct {
	ID                uint64
	Name              string
	Clients           ClientSet
	Speaker           *Client
	Capacity          uint
	RequireOccupation bool

	// Using a syncMutex here to be able to lock state before editing clients
	// Could also use Channels to block
	sync.RWMutex
}

func NewRoom(name string, speaker *Client) *Room {
	return &Room{
		Name:    name,
		Speaker: speaker,
		Clients: make(ClientSet),
	}
}

func (r *Room) GetRoomInfoEvent() (*Event, error) {

	r.RLock()
	defer r.RUnlock()

	return NewRoomInfoEvent(r)
}
func (r *Room) Broadcast(e *Event) {

	r.RLock()
	defer r.RUnlock()

	for client := range r.Clients {

		client.send <- *e
	}
}

func (r *Room) AddClient(c *Client) {

	r.Lock()
	defer r.Unlock()

	r.Clients[c] = true

	if len(r.Clients) == 1 {

		r.Speaker = c
	}

	log.Debugf("Client joined room; Now has %d members", len(r.Clients))
}

func (r *Room) RemoveClient(c *Client) {

	r.Lock()
	defer r.Unlock()

	if _, ok := r.Clients[c]; !ok {
		return
	}

	c.connection.Close()
	delete(r.Clients, c)

	if r.Speaker == c {

		for key := range r.Clients {
			r.Speaker = key
			break
		}
	}

	log.Debugf("Client left room; Now has %d members", len(r.Clients))
}

func (m *RoomManager) CreateRoomGET(c echo.Context) error {

	rid, err := m.CreateRoom()

	if err != nil {
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Server error")
	}

	return c.String(http.StatusOK, strconv.FormatUint(rid, 10))
}
