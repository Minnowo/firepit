package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/EZCampusDevs/firepit/data"
	"github.com/EZCampusDevs/firepit/database"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type AuthHandler struct {
	AuthSecret   []byte
	TokenTimeout time.Duration
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func (a *AuthHandler) ValidateFatal() {

	if len(a.AuthSecret) < 16 {

		log.Fatalf("JWT Auth has bad or weak secret. Use %s to specify a secret", data.ENV_JWT_KEY)
	}
}

func (a *AuthHandler) BasicPayloadCheck(c echo.Context, payload *data.AuthPayload) error {

	if err := c.Bind(payload); err != nil {

		return fmt.Errorf("Invalid JSON payload")
	}

	if err := payload.IsValid(); err != nil {

		return err
	}

	return nil
}

func (a *AuthHandler) CreateJWTFromAuthPayload(authPayload *data.AuthPayload) (string, error) {

	expirationTime := time.Now().Add(a.TokenTimeout)

	log.Infof("Creating new JWT. Expires (RFC3339):", expirationTime.Format(time.RFC3339))

	claims := &Claims{

		Username: authPayload.Username,

		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	return token.SignedString(a.AuthSecret)
}

func (a *AuthHandler) POSTCreateUser(c echo.Context) error {

	var authPayload data.AuthPayload

	if err := a.BasicPayloadCheck(c, &authPayload); err != nil {

		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if !database.CreateUser(&authPayload) {

		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Failed to create user"})
	}

	log.Infof("Created new user with name %s", authPayload.Username)

	finalToken, err := a.CreateJWTFromAuthPayload(&authPayload)

	if err != nil {

		log.Error(err)

		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}

	return c.JSON(http.StatusOK, map[string]string{"token": finalToken})
}

func (a *AuthHandler) POSTCreateJWT(c echo.Context) error {

	var authPayload data.AuthPayload

	if err := a.BasicPayloadCheck(c, &authPayload); err != nil {

		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if !database.IsCredentialsValid(&authPayload) {

		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid credentials"})
	}

	finalToken, err := a.CreateJWTFromAuthPayload(&authPayload)

	if err != nil {

		log.Error(err)

		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}

	return c.JSON(http.StatusOK, map[string]string{"token": finalToken})
}

func GETEnsureAuthed(c echo.Context) error {

	return c.String(http.StatusOK, "Hello, World!")
}
