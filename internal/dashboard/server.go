package dashboard

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	
	"ralph2/pkg/utils"
)

type Server struct {
	e   *echo.Echo
	bus *utils.EventBus
	mu  sync.Mutex
	// In production, sync.Map or fan-out manager for many clients
}

func NewServer(bus *utils.EventBus) *Server {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	s := &Server{
		e:   e,
		bus: bus,
	}

	e.GET("/events", s.handleEvents)
	e.GET("/", s.handleIndex)

	return s
}

func (s *Server) Start(addr string) error {
	return s.e.Start(addr)
}

func (s *Server) handleIndex(c echo.Context) error {
	return c.HTML(http.StatusOK, "<h1>Ralph2 Dashboard</h1><p>Connect to /events for live updates.</p>")
}

func (s *Server) handleEvents(c echo.Context) error {
	c.Response().Header().Set(echo.HeaderContentType, "text/event-stream")
	c.Response().Header().Set(echo.HeaderCacheControl, "no-cache")
	c.Response().Header().Set(echo.HeaderConnection, "keep-alive")

	
	// Subscribe to state_change events
	updates := s.bus.Subscribe("state_change")
	defer s.bus.Unsubscribe(updates)

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-c.Request().Context().Done():
			return nil
		case event := <-updates:
			data, err := json.Marshal(event)
			if err != nil {
				continue
			}
			fmt.Fprintf(c.Response(), "data: %s\n\n", data)
			c.Response().Flush()
		case <-ticker.C:
			// Heartbeat comment
			fmt.Fprintf(c.Response(), ": keepalive\n\n")
			c.Response().Flush()
		}
	}
}
