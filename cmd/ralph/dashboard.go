package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"ralph2/internal/core"
	"ralph2/internal/server"
	"ralph2/pkg/utils"
)

var dashboardCmd = &cobra.Command{
	Use:   "dashboard",
	Short: "Start the web dashboard",
	Long:  `Start the local web server to view the Ralph dashboard in your browser.`,
	Run: func(cmd *cobra.Command, args []string) {
		startDashboard()
	},
}

func init() {
	rootCmd.AddCommand(dashboardCmd)
}

func startDashboard() {
	fmt.Println("Starting Ralph2 Dashboard...")

	// 1. Initialize EventBus
	bus := utils.NewEventBus()

	// 2. Initialize FSM (Optional: Dashboard might just view state, but let's have an FSM running)
	// For now, we just run the server. In a real app, we might want to attach to an existing ralph run
	// or start a new one. This story focuses on the server itself.
	fsm := core.NewStateManager(bus)
	_ = fsm // usage to avoid compiler error if we don't use it yet

	// 3. Initialize Server
	srv := server.NewServer(bus, 5000)

	// 4. Start Server in goroutine
	go func() {
		if err := srv.Start(); err != nil {
			fmt.Printf("Server error: %v\n", err)
		}
	}()

	fmt.Println("Dashboard running at http://localhost:5000")

	// 5. Handle Interrupts
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// 6. Graceful Shutdown
	fmt.Println("Shutting down...")
	// TODO: implement context timeout shutdown logic used in Start()
}
