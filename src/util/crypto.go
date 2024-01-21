package util

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
)

func NBytesInt(bytes uint) (uint64, error) {

	if bytes == 0 || bytes > 8 {
		return 0, fmt.Errorf("Bytes must be 1:8")
	}

	var randomBytes [8]byte

	_, err := rand.Read(randomBytes[:bytes])

	if err != nil {
		return 0, err
	}

	num := binary.LittleEndian.Uint64(randomBytes[:])

	return num, nil
}
