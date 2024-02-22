package main

import (
	"os"
	"strconv"

	"github.com/EZCampusDevs/firepit/database"
	"github.com/EZCampusDevs/firepit/handler"
	"github.com/EZCampusDevs/firepit/handler/websocket"
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

	e = echo.New()
	m = websocket.NewManager()
	c = getDBConf()

	initLogging(e)

	database.DBInit(c)

	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.GET("/ws", m.ServeWebsocket)

	e.Static("/static", "static")

	roomGroup := e.Group("/room")
	roomGroup.GET("/new", m.GetRoomManager().CreateRoomGET)

	quoteGroup := e.Group("/quote")
	quoteGroup.GET("", handler.GetRandomQuote)
	quoteGroup.POST("", handler.CreateNewQuote)

	e.RouteNotFound("/", handler.Heartbeat)

	e.Logger.Fatal(e.Start(":3000"))
}
