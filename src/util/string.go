package util

import (
	"strings"

	"github.com/google/uuid"
)

func GetUUID() string {
	return uuid.NewString()
}

func IsEmptyOrWhitespace(s string) bool {
	return strings.TrimSpace(s) == ""
}
