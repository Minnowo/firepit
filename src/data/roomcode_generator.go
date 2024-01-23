package data

import (
	"strconv"

	"github.com/EZCampusDevs/firepit/util"
)

// behavior for generating room codes
type RoomCodeGenerator interface {
	// gets a new room code or error
	GetRoomCode() (string, error)
}

// generates a room code using a random number
// with n bytes in the given base
type UintNRoomCodeGenerator struct {
	base uint16
	n    uint16
}

// Creates a new UintNRoomCodeGenerator with the given n and base
// Panics if n < 0 || n > 8
// Panics if base < 2 || base > 36
func NewUintNRoomCodeGenerator(n uint16, base uint16) UintNRoomCodeGenerator {
	if n < 0 || n > 8 {
		panic("UintNRoomCodeGenerator cannot have n < 0 or n > 8")
	}

	if base < 2 || base > 36 {
		panic("UintNRoomCodeGenerator cannot have base < 2 or base > 36")
	}

	return UintNRoomCodeGenerator{
		n:    n,
		base: base,
	}
}

func (u UintNRoomCodeGenerator) GetRoomCode() (string, error) {

	i, err := util.NBytesInt(uint(u.n))

	if err != nil {
		return "", err
	}

	return strconv.FormatUint(i, int(u.base)), nil
}
