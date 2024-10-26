package main

import (
	"github.com/EZCampusDevs/firepit/data"
	"github.com/EZCampusDevs/firepit/handler"
	"github.com/EZCampusDevs/firepit/handler/websocket"
	"github.com/EZCampusDevs/firepit/logging"
	"github.com/EZCampusDevs/firepit/ui"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	var e *echo.Echo
	var m *websocket.Manager

	e = echo.New()
	m = websocket.NewManager()

	logging.InitFromEnv()

	e.Use(middleware.Recover())

	e.GET("/ws", m.ServeWebsocket)

	if data.IS_DEBUG {
		e.GET("/helpme", m.PrintDebugStuff)
	} else {
		// e.Use(middleware.HTTPSRedirect())
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"https://astoryofand.com",
				"https://*.astoryofand.com"},
		}))
	}

	e.Static("/static", "static")

	roomGroup := e.Group("/room")
	roomGroup.GET("/new", m.GetRoomManager().GETCreateRoom)
	roomGroup.GET("/check/:rid", m.GetRoomManager().GETHasRoom)

	e.RouteNotFound("/", handler.GETHeartbeat)

	ui.Register(e)

	e.Logger.Fatal(e.Start(":3000"))
}
