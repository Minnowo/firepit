package util

import (
	"crypto/rand"
	"encoding/binary"
)

func GenerateFull64BitNumber() (uint64, error) {

	var randomBytes [8]byte

	_, err := rand.Read(randomBytes[:])

	if err != nil {
		return 0, err
	}

	full64BitNumber := binary.LittleEndian.Uint64(randomBytes[:])

	return full64BitNumber, nil
}
