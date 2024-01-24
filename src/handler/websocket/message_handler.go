package websocket

import (
	"encoding/json"

	"github.com/EZCampusDevs/firepit/data"
	"github.com/labstack/gommon/log"
)

// Handles the client to server message for setting the speaker
func handleSetSpeaker(e_ Event, c *Client) error {
	log.Debug("EVENT__CLIENT_SET_SPEAKER MESSAGE: ", e_)

	// pause the room thread so we can copy the speaker ptr safely
	c.room.state <- data.CHAN__PAUSED
	speaker := c.room.Speaker
	c.room.state <- data.CHAN__RUNNING

	// check the client is the speaker
	if speaker != c {
		log.Debug("Client cannot change speaker because they are not the speaker")
		return nil
	}

	var e SetSpeakerEvent = SetSpeakerEvent{}

	if err := json.Unmarshal(e_.Payload, &e); err != nil {
		log.Error(e_.Payload)
		return err
	}

	c.room.setSpeakerById <- e.SpeakerID

	return nil
}
