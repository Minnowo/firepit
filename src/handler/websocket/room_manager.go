package websocket

import (
	"fmt"
	"sync"

	"github.com/EZCampusDevs/firepit/data"
	"github.com/rs/zerolog/log"
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
			continue
		}

		r.AddRoom(rid)

		log.Debug().Str("room", rid).Int("roomCount", len(r.rooms)).Msg("Created room")

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
