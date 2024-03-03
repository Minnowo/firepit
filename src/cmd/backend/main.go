package main

import (
	"os"
	"strconv"
	"time"

	"github.com/EZCampusDevs/firepit/data"
	"github.com/EZCampusDevs/firepit/database"
	"github.com/EZCampusDevs/firepit/handler"
	"github.com/EZCampusDevs/firepit/handler/websocket"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

func initLogging(app *echo.Echo) {

	IS_DEBUG := os.Getenv("DEBUG")
	LOG_LEVEL := os.Getenv("LOG_LEVEL")

	log.SetHeader("${time_rfc3339} ${level}")
	log.SetLevel(log.INFO)

	app.Logger.SetHeader("${time_rfc3339} ${level}")

	if level, err := strconv.ParseUint(LOG_LEVEL, 10, 8); err == nil {
		app.Logger.SetLevel(log.Lvl(level))
		log.Info("Read LOG_LEVEL from env: ", level)
	} else {
		log.Warn("Could not read LOG_LEVEL from env. Log level is: ", app.Logger.Level())
	}

	if debug_, err := strconv.ParseBool(IS_DEBUG); err == nil {

		app.Debug = debug_

		if debug_ {
			app.Logger.SetLevel(log.DEBUG)
		}

		log.Info("Read DEBUG from Env: ", debug_)
	} else {
		app.Debug = false
		log.Warn("Could not read DEBUG from env. Running in release mode.")
	}

	log.SetLevel(app.Logger.Level())
}

func getDBConf() *database.DBConfig {
	return &database.DBConfig{
		Username:     "root",
		Password:     "root",
		Hostname:     "127.0.0.1",
		Port:         3306,
		DatabaseName: "firepit-mariadb",
	}
}

func main() {

	var e *echo.Echo
	var m *websocket.Manager
	var c *database.DBConfig
	var a *handler.AuthHandler

	e = echo.New()
	m = websocket.NewManager()
	c = getDBConf()
	a = &handler.AuthHandler{
		AuthSecret:   []byte(os.Getenv(data.ENV_JWT_KEY)),
		TokenTimeout: 60 * time.Minute,
	}

	initLogging(e)

	a.ValidateFatal()

	database.DBInit(c)

	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.GET("/ws", m.ServeWebsocket)

	e.Static("/static", "static")

	roomGroup := e.Group("/room")
	roomGroup.GET("/new", m.GetRoomManager().GETCreateRoom)
	roomGroup.GET("/:rid", m.GetRoomManager().GETHasRoom)

	quoteGroup := e.Group("/quote")
	quoteGroup.GET("", handler.GETRandomQuote)

	auth := e.Group("/auth")
	auth.POST("/token", a.POSTCreateJWT)
	auth.POST("/create", a.POSTCreateUser)

	jwtMiddleware := echojwt.WithConfig(echojwt.Config{
		SigningKey: a.AuthSecret,
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(handler.JWTClaims)
		},
		SigningMethod: "HS512", // see auth.go for the AuthCreateJWT function
	})

	authed := e.Group("/authed", jwtMiddleware)
	authed.GET("/refresh", a.GETRefreshJWT)
	authed.POST("/quote", handler.GETCreateNewQuote)

	e.RouteNotFound("/", handler.GETHeartbeat)

	e.Logger.Fatal(e.Start(":3000"))
}
