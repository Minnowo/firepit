package data

import (
	"os"
	"strings"
)

const (
	ENV_JWT_KEY               = "JWT_SECRET"
	ENV_DATABASE_HOSTNAME_KEY = "DB_HOSTNAME"
	ENV_DATABASE_NAME_KEY     = "DB_NAME"
	ENV_DATABASE_PASSWORD_KEY = "DB_PASSWORD"
	ENV_DATABASE_USERNAME_KEY = "DB_USERNAME"
	ENV_DATABASE_PORT_KEY     = "DB_PORT"
	ENV_DEBUG_KEY             = "DEBUG"
)

var _DEBUG_ENV_VALUE = strings.ToLower(os.Getenv(ENV_DEBUG_KEY))

var IS_DEBUG bool = _DEBUG_ENV_VALUE == "1" || _DEBUG_ENV_VALUE == "true"
