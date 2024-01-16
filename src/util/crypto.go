package util

import (
	"crypto/rand"
	"encoding/binary"
)

// Generates a 64bit number where all 8 bytes are not 0
func GenerateFull64BitNumber() (uint64, error) {

	var randomBytes [8]byte

	_, err := rand.Read(randomBytes[:])

	if err != nil {
		return 0, err
	}

	full64BitNumber := binary.LittleEndian.Uint64(randomBytes[:])

	return full64BitNumber, nil
}
