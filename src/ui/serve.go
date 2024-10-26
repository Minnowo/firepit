package ui

import (
	"embed"
	"io"
	"io/fs"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

//go:embed dist
var uiDir embed.FS
var buildDir, _ = fs.Sub(uiDir, "dist")

// Register registers the ui on the root path.
func Register(e *echo.Echo) {
	e.GET("/", serveFile("index.html", "text/html"))
	e.GET("/index.html", serveFile("index.html", "text/html"))
	e.GET("/vite.svg", serveFile("vite.svg", "image/svg+xml"))
	// e.GET("/manifest.json", serveFile("manifest.json", "application/json"))
	// e.GET("/service-worker.js", serveFile("service-worker.js", "text/javascript"))
	// e.GET("/asset-manifest.json", serveFile("asset-manifest.json", "application/json"))
	e.GET("/assets/:resource", echo.WrapHandler(http.FileServer(http.FS(buildDir))))

	// e.GET("/favicon.ico", serveFile("favicon.ico", "image/x-icon"))

	// for _, size := range []string{"16x16", "32x32", "192x192", "256x256"} {
	// 	fileName := fmt.Sprintf("favicon-%s.png", size)
	// 	e.GET("/"+fileName, serveFile(fileName, "image/png"))
	// }
}

func serveFile(name, contentType string) echo.HandlerFunc {

	file, err := buildDir.Open(name)

	if err != nil {
		log.Panic().Err(err).Msgf("could not find %s", file)
	}

	defer file.Close()

	content, err := io.ReadAll(file)

	if err != nil {
		log.Panic().Err(err).Msgf("could not read %s", file)
	}

	return func(c echo.Context) error {

		log.Debug().Str("file", name).Msg("Serving file")

		c.Response().Header().Set("Content-Type", contentType)

		return c.Blob(http.StatusOK, contentType, content)
	}
}
