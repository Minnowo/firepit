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

	if data.IS_DEBUG {

		app.Debug = data.IS_DEBUG

		if data.IS_DEBUG {
			app.Logger.SetLevel(log.DEBUG)
		}

		log.Info("Read DEBUG from Env: ", data.IS_DEBUG)
	} else {
		app.Debug = false
		log.Warn("Could not read DEBUG from env. Running in release mode.")
	}

	log.SetLevel(app.Logger.Level())
}

func getDBConf() *database.DBConfig {

	var username string
	var password string
	var hostname string
	var databasename string = "firepit-mariadb"
	var port int = 3306

	if u, e := os.LookupEnv(data.ENV_DATABASE_USERNAME_KEY); e {
		username = u
	} else {
		log.Fatalf("Must set %s to the database username", data.ENV_DATABASE_USERNAME_KEY)
	}

	if p, e := os.LookupEnv(data.ENV_DATABASE_PASSWORD_KEY); e {
		password = p
	} else {
		log.Fatalf("Must set %s to the database password", data.ENV_DATABASE_PASSWORD_KEY)
	}

	if h, e := os.LookupEnv(data.ENV_DATABASE_HOSTNAME_KEY); e {
		hostname = h
	} else {
		log.Fatalf("Must set %s to the database hostname", data.ENV_DATABASE_HOSTNAME_KEY)
	}

	if n, e := os.LookupEnv(data.ENV_DATABASE_NAME_KEY); e {
		databasename = n
	} else {
		log.Warnf("%s was not set, using default value of %s", data.ENV_DATABASE_NAME_KEY, databasename)
	}

	if p, e := os.LookupEnv(data.ENV_DATABASE_PORT_KEY); e {

		i, err := strconv.Atoi(p)

		if err == nil {
			port = i
		} else {
			log.Warnf("%s had an invalid integer. Using default value of %d", data.ENV_DATABASE_PORT_KEY, port)
		}

	} else {
		log.Warnf("%s was not set, using default value of %d", data.ENV_DATABASE_PORT_KEY, port)
	}

	return &database.DBConfig{
		Username:     username,
		Password:     password,
		Hostname:     hostname,
		Port:         port,
		DatabaseName: databasename,
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

	e.GET("/ws", m.ServeWebsocket)

	if data.IS_DEBUG {
		e.GET("/helpme", m.PrintDebugStuff)
	} else {
		e.Use(middleware.HTTPSRedirect())
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"https://astoryofand.com",
				"https://*.astoryofand.com"},
		}))
	}

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
