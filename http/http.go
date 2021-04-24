package http

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/serizawa-jp/podcast-server/podcast"
)

func NewHTTPHandler(
	podcaster *podcast.Podcaster,
	targetDir string,
) (http.Handler, error) {
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(middleware.Logger())
	e.Use(noCacheMiddleware)

	e.GET("/", func(c echo.Context) error {
		c.Blob(http.StatusOK, "application/xml", []byte(podcaster.GetFeed()))
		return nil
	})
	e.GET("/sync", func(c echo.Context) error {
		podcaster.Sync()
		c.String(http.StatusOK, "")
		return nil
	})

	e.GET("/rss.xml", func(c echo.Context) error {
		c.Blob(http.StatusOK, "application/xml", []byte(podcaster.GetFeed()))
		return nil
	})

	e.Static("/static", targetDir)

	return e, nil
}
