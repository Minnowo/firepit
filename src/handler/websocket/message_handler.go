package websocket

import (
	"encoding/json"

	"github.com/EZCampusDevs/firepit/data"
	"github.com/rs/zerolog/log"
)

// Handles the client to server message for setting the speaker
func handleSetSpeaker(e_ Event, c *Client) error {

	log.Debug().Msg("Handling set speaker message")

	// pause the room thread so we can copy the speaker ptr safely
	c.room.state <- data.CHAN__PAUSED
	speaker := c.room.Speaker
	c.room.state <- data.CHAN__RUNNING

	// check the client is the speaker
	if speaker != c {

		log.Debug().Str("client", c.info.DisplayId).Msg("Client cannot change speaker because they are not the speaker")

		return nil
	}

	var e SetSpeakerEvent = SetSpeakerEvent{}

	if err := json.Unmarshal(e_.Payload, &e); err != nil {

		log.Error().Str("client", c.info.DisplayId).Bytes("msg", e_.Payload).Msg("Failed to Unmarshal json from client message")

		return err
	}

	c.room.setSpeakerById <- e.ClientInfo.DisplayId

	return nil
}
