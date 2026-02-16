package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"ralph2/pkg/utils"
	"ralph2/web"
)

type Server struct {
	e        *echo.Echo
	eventBus *utils.EventBus
	port     int
}

func NewServer(bus *utils.EventBus, port int) *Server {
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	
	s := &Server{
		e:        e,
		eventBus: bus,
		port:     port,
	}

	// Serve static files
	distFS, err := web.GetDistFS()
	if err == nil {
		// Serve index.html for root, and other assets as needed
		// echo.WrapHandler + http.FileServer is one way, 
		// but Echo has e.StaticFS which is easier for root
		e.GET("/*", echo.WrapHandler(http.StripPrefix("/", http.FileServer(http.FS(distFS)))))
	} else {
		e.GET("/", s.handleIndex)
	}

	// Register API routes
	e.GET("/events", s.handleEvents)

	return s
}

func (s *Server) Start() error {
	return s.e.Start(fmt.Sprintf(":%d", s.port))
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.e.Shutdown(ctx)
}

func (s *Server) handleIndex(c echo.Context) error {
	return c.String(http.StatusOK, "<h1>Ralph2 Web Dashboard</h1><p>Waiting for embedded UI...</p>")
}

func (s *Server) handleEvents(c echo.Context) error {
	// SSE headers
	c.Response().Header().Set(echo.HeaderContentType, "text/event-stream")
	c.Response().Header().Set(echo.HeaderCacheControl, "no-cache")
	c.Response().Header().Set(echo.HeaderConnection, "keep-alive")

	// Subscribe to all events for now
	sub := s.eventBus.Subscribe("state_change")
	// TODO: Subscribe to logging events too when available

	// Cleanup on disconnect
	defer s.eventBus.Unsubscribe(sub)

	// Stream loop
	for {
		select {
		case <-c.Request().Context().Done():
			return nil
		case event := <-sub:
			data, err := json.Marshal(event)
			if err != nil {
				continue
			}
			fmt.Fprintf(c.Response(), "data: %s\n\n", data)
			c.Response().Flush()
		}
	}
}
