package gui

import (
	"fmt"

	"github.com/MhemedAbderrahmen/lazyguepard/internal/guepard"
	"github.com/MhemedAbderrahmen/lazyguepard/internal/state"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Layout manages the main UI layout
type Layout struct {
	App             *tview.Application
	Pages           *tview.Pages
	MainFlex        *tview.Flex
	DeploymentsList *tview.List
	BranchesView    *tview.TextView
	CommitsView     *tview.TextView
	StatusView      *tview.TextView
	currentPanel    int
	panels          []tview.Primitive
	Client          *guepard.Client
	State           *state.AppState
}

// NewLayout creates a new layout
func NewLayout(app *tview.Application, client *guepard.Client) *Layout {
	layout := &Layout{
		App:    app,
		Pages:  tview.NewPages(),
		Client: client,
		State:  state.NewAppState(),
	}

	layout.setupPanels()
	layout.setupLayout()
	layout.setupKeybindings()

	return layout
}
func (l *Layout) setupPanels() {
	l.DeploymentsList = tview.NewList().
		ShowSecondaryText(false)
	l.DeploymentsList.SetBorder(true).
		SetTitle(" Deployments (1) ").
		SetTitleColor(tcell.ColorYellow)

	// Handle deployment selection
	l.DeploymentsList.SetSelectedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		l.onDeploymentSelected(index)
	})

	// Branches panel (top-right)
	l.BranchesView = tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true)
	l.BranchesView.SetBorder(true).
		SetTitle(" Branches (2) ")

	// Commits panel (bottom-left)
	l.CommitsView = tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true)
	l.CommitsView.SetBorder(true).
		SetTitle(" Commits (3) ")

	// Status panel (bottom-right)
	l.StatusView = tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true)
	l.StatusView.SetBorder(true).
		SetTitle(" Status (4) ")

	// Store panels for navigation
	l.panels = []tview.Primitive{
		l.DeploymentsList,
		l.BranchesView,
		l.CommitsView,
		l.StatusView,
	}
	l.currentPanel = 0
}
func (l *Layout) setupLayout() {
	topRow := tview.NewFlex().
		AddItem(l.DeploymentsList, 0, 1, true).
		AddItem(l.BranchesView, 0, 1, false)

	// Bottom row: Commits | Status
	bottomRow := tview.NewFlex().
		AddItem(l.CommitsView, 0, 1, false).
		AddItem(l.StatusView, 0, 1, false)

	// Main layout: top and bottom rows
	l.MainFlex = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(topRow, 0, 1, true).
		AddItem(bottomRow, 0, 1, false)

	l.Pages.AddPage("main", l.MainFlex, true, true)
}
func (l *Layout) setupKeybindings() {
	l.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab:
			l.nextPanel()
			return nil
		case tcell.KeyBacktab: // Shift+Tab
			l.previousPanel()
			return nil
		case tcell.KeyCtrlC, tcell.KeyEsc:
			l.App.Stop()
			return nil
		}

		// Number keys 1-4 for direct panel access
		switch event.Rune() {
		case '1':
			l.switchToPanel(0)
			return nil
		case '2':
			l.switchToPanel(1)
			return nil
		case '3':
			l.switchToPanel(2)
			return nil
		case '4':
			l.switchToPanel(3)
			return nil
		case 'q':
			l.App.Stop()
			return nil
		case 'r':
			l.refreshDeployments()
			return nil
		}

		return event
	})
}

func (l *Layout) nextPanel() {
	l.currentPanel = (l.currentPanel + 1) % len(l.panels)
	l.highlightCurrentPanel()
	l.App.SetFocus(l.panels[l.currentPanel])
}
func (l *Layout) previousPanel() {
	l.currentPanel = (l.currentPanel - 1 + len(l.panels)) % len(l.panels)
	l.highlightCurrentPanel()
	l.App.SetFocus(l.panels[l.currentPanel])
}

func (l *Layout) switchToPanel(index int) {
	if index >= 0 && index < len(l.panels) {
		l.currentPanel = index
		l.highlightCurrentPanel()
		l.App.SetFocus(l.panels[l.currentPanel])
	}
}

func (l *Layout) highlightCurrentPanel() {
	// Reset all panel title colors
	l.DeploymentsList.SetTitleColor(tcell.ColorWhite)
	l.BranchesView.SetTitleColor(tcell.ColorWhite)
	l.CommitsView.SetTitleColor(tcell.ColorWhite)
	l.StatusView.SetTitleColor(tcell.ColorWhite)

	// Highlight current panel
	switch l.currentPanel {
	case 0:
		l.DeploymentsList.SetTitleColor(tcell.ColorYellow)
	case 1:
		l.BranchesView.SetTitleColor(tcell.ColorYellow)
	case 2:
		l.CommitsView.SetTitleColor(tcell.ColorYellow)
	case 3:
		l.StatusView.SetTitleColor(tcell.ColorYellow)
	}
}

func (l *Layout) LoadDeployments() {
	l.DeploymentsList.Clear()
	l.DeploymentsList.AddItem("Loading deployments...", "", 0, nil)

	go func() {
		output, err := l.Client.ListDeployments()
		l.App.QueueUpdateDraw(func() {
			l.DeploymentsList.Clear()
			if err != nil {
				l.DeploymentsList.AddItem(fmt.Sprintf("Error: %v", err), "", 0, nil)
				return
			}

			deployments := guepard.ParseDeployments(output)
			l.State.Deployments = deployments

			if len(deployments) == 0 {
				l.DeploymentsList.AddItem("No deployments found", "", 0, nil)
				return
			}

			for i, dep := range deployments {
				l.DeploymentsList.AddItem(dep.Name, "", rune('a'+i), nil)
			}
		})
	}()
}

func (l *Layout) refreshDeployments() {
	l.StatusView.SetText("[yellow]Refreshing deployments...")
	l.LoadDeployments()
	l.StatusView.SetText("[green]✓ Refreshed")
}
func (l *Layout) onDeploymentSelected(index int) {
	if index >= len(l.State.Deployments) {
		return
	}

	deployment := l.State.Deployments[index]
	l.State.SelectedDeployment = &deployment

	l.StatusView.SetText(fmt.Sprintf("[green]Selected:[white] %s\n\n[dim]Loading branches and commits...", deployment.Name))

	// TODO: Load branches and commits for this deployment
	l.BranchesView.SetText("[yellow]Branches will appear here")
	l.CommitsView.SetText("[yellow]Commits will appear here")
}
