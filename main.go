package main

import (
	"fmt"
	"os"

	"github.com/MhemedAbderrahmen/lazyguepard/internal/guepard"
	"github.com/MhemedAbderrahmen/lazyguepard/internal/gui"
	"github.com/rivo/tview"
)

func main() {
	// Initialize Guepard client
	client, err := guepard.NewClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		fmt.Println("Make sure Guepard CLI is installed and in your PATH")
		os.Exit(1)
	}

	// Create TUI application
	app := tview.NewApplication()

	// Create layout
	layout := gui.NewLayout(app, client)

	// Set initial content
	layout.StatusView.SetText("[green]🐆 LazyGuepard[white]\n\nControls:\n• ↑/↓: Navigate\n• Enter: Select\n• Tab: Switch panels\n• r: Refresh\n• q: Quit")

	// Load deployments
	layout.LoadDeployments()

	// Run the application
	if err := app.SetRoot(layout.Pages, true).Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
