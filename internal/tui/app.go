package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/JHill6253/seats-aero-cli/internal/api"
	"github.com/JHill6253/seats-aero-cli/internal/config"
)

// View represents different screens in the TUI
type View int

const (
	ViewSearch View = iota
	ViewResults
	ViewLoading
	ViewError
)

// Model is the main TUI model
type Model struct {
	config   *config.Config
	client   *api.Client
	view     View
	search   SearchModel
	results  ResultsModel
	err      error
	width    int
	height   int
	quitting bool
}

// searchResultMsg is sent when search results are received
type searchResultMsg struct {
	results []api.Availability
}

// searchErrorMsg is sent when a search fails
type searchErrorMsg struct {
	err error
}

// NewModel creates a new TUI model
func NewModel(cfg *config.Config) Model {
	client := api.NewClient(cfg.GetAPIKey())

	return Model{
		config:  cfg,
		client:  client,
		view:    ViewSearch,
		search:  NewSearchModel(cfg),
		results: NewResultsModel(),
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles messages
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			if m.view == ViewSearch {
				m.quitting = true
				return m, tea.Quit
			}
			// From other views, go back to search
			m.view = ViewSearch
			return m, nil
		case "esc":
			if m.view != ViewSearch {
				m.view = ViewSearch
				return m, nil
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.results.SetSize(msg.Width, msg.Height-10)

	case searchResultMsg:
		m.results.SetResults(msg.results)
		m.view = ViewResults
		return m, nil

	case searchErrorMsg:
		m.err = msg.err
		m.view = ViewError
		return m, nil
	}

	var cmd tea.Cmd

	switch m.view {
	case ViewSearch:
		m.search, cmd = m.search.Update(msg)

		// Check if search was submitted
		if m.search.submitted {
			m.search.submitted = false
			m.view = ViewLoading
			return m, m.doSearch()
		}

	case ViewResults:
		m.results, cmd = m.results.Update(msg)

		// Check for export request
		if m.results.exportRequested {
			m.results.exportRequested = false
			// Handle export
		}
	}

	return m, cmd
}

// View renders the UI
func (m Model) View() string {
	if m.quitting {
		return ""
	}

	var content string

	switch m.view {
	case ViewSearch:
		content = m.search.View()
	case ViewResults:
		content = m.results.View()
	case ViewLoading:
		content = m.loadingView()
	case ViewError:
		content = m.errorView()
	}

	return AppStyle.Render(content)
}

func (m Model) loadingView() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		TitleStyle.Render("Searching..."),
		"",
		SpinnerStyle.Render("Fetching availability from seats.aero..."),
		"",
		HelpStyle.Render("Press Esc to cancel"),
	)
}

func (m Model) errorView() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		TitleStyle.Render("Error"),
		"",
		ErrorStyle.Render(fmt.Sprintf("Error: %v", m.err)),
		"",
		HelpStyle.Render("Press Esc to go back"),
	)
}

func (m Model) doSearch() tea.Cmd {
	return func() tea.Msg {
		params := m.search.GetSearchParams()

		resp, err := m.client.Search(params)
		if err != nil {
			return searchErrorMsg{err: err}
		}

		return searchResultMsg{results: resp.Data}
	}
}

// Run starts the TUI
func Run(cfg *config.Config) error {
	if err := cfg.Validate(); err != nil {
		return err
	}

	p := tea.NewProgram(
		NewModel(cfg),
		tea.WithAltScreen(),
	)

	_, err := p.Run()
	return err
}
