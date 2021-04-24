package http

import (
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/serizawa-jp/podcast-server/podcast"
)

var (
	// Unix epoch time
	epoch = time.Unix(0, 0).Format(time.RFC1123)

	// Taken from https://github.com/mytrile/nocache
	noCacheHeaders = map[string]string{
		"Expires":         epoch,
		"Cache-Control":   "no-cache, private, max-age=0",
		"Pragma":          "no-cache",
		"X-Accel-Expires": "0",
	}
	etagHeaders = []string{
		"ETag",
		"If-Modified-Since",
		"If-Match",
		"If-None-Match",
		"If-Range",
		"If-Unmodified-Since",
	}
)

func NewHTTPHandler(
	podcaster *podcast.Podcaster,
	targetDir string,
) (http.Handler, error) {
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(middleware.Logger())
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			req := c.Request()
			// Delete any ETag headers that may have been set
			for _, v := range etagHeaders {
				if req.Header.Get(v) != "" {
					req.Header.Del(v)
				}
			}

			// Set our NoCache headers
			res := c.Response()
			for k, v := range noCacheHeaders {
				res.Header().Set(k, v)
			}
			return next(c)
		}
	})

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
