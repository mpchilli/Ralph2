package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"ralph2/pkg/utils"
)

func TestServerStartsAndStops(t *testing.T) {
	bus := utils.NewEventBus()
	srv := NewServer(bus, 5001)

	// Start in goroutine
	done := make(chan bool)
	go func() {
		if err := srv.Start(); err != nil {
			t.Logf("Server stopped: %v", err)
		}
		done <- true
	}()

	// Wait for startup
	time.Sleep(100 * time.Millisecond)

	// Verify index
	resp, err := http.Get("http://localhost:5001/")
	if err != nil {
		t.Fatalf("Failed to call index: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
	resp.Body.Close()

	// Shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	
	if err := srv.Shutdown(ctx); err != nil {
		t.Errorf("Shutdown failed: %v", err)
	}

	<-done
}

func TestSSEHandler(t *testing.T) {
	bus := utils.NewEventBus()
	s := NewServer(bus, 0)

	// Setup request
	req := httptest.NewRequest(http.MethodGet, "/events", nil)
	rec := httptest.NewRecorder()
	
	// Create context with cancel to simulate disconnect
	ctx, cancel := context.WithCancel(context.Background())
	req = req.WithContext(ctx)
	
	c := s.e.NewContext(req, rec)

	// Run handler in goroutine
	done := make(chan error)
	go func() {
		done <- s.handleEvents(c)
	}()

	// Wait for subscription
	time.Sleep(100 * time.Millisecond)

	// Publish
	bus.Publish("state_change", "UNIT_TEST_STATE")

	// Wait for processing
	time.Sleep(100 * time.Millisecond)

	// Cancel context to stop handler
	cancel()
	
	// Wait for handler to return
	select {
	case <-done:
		// success
	case <-time.After(1 * time.Second):
		t.Fatal("Handler did not return after cancellation")
	}

	// Verify Content Type using Echo constant
	if ct := rec.Header().Get(echo.HeaderContentType); ct != "text/event-stream" {
		t.Errorf("Expected Content-Type text/event-stream, got %s", ct)
	}

	// Check body
	body := rec.Body.String()
	if !strings.Contains(body, "UNIT_TEST_STATE") {
		t.Errorf("Body content missing event: %s", body)
	}
	if !strings.Contains(body, "state_change") {
		t.Errorf("Body content missing topic: %s", body)
	}
}
