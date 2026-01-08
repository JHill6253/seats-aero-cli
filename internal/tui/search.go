package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/JHill6253/seats-aero-cli/internal/api"
	"github.com/JHill6253/seats-aero-cli/internal/config"
)

// Field indices
const (
	fieldOrigin = iota
	fieldDestination
	fieldStartDate
	fieldEndDate
	fieldCabin
	fieldSource
	numFields
)

// SearchModel handles the search form
type SearchModel struct {
	inputs    []textinput.Model
	focused   int
	submitted bool
	config    *config.Config
}

// NewSearchModel creates a new search form
func NewSearchModel(cfg *config.Config) SearchModel {
	inputs := make([]textinput.Model, numFields)

	// Origin
	inputs[fieldOrigin] = textinput.New()
	inputs[fieldOrigin].Placeholder = "SFO, LAX"
	inputs[fieldOrigin].CharLimit = 50
	inputs[fieldOrigin].Width = 30
	inputs[fieldOrigin].Prompt = ""
	inputs[fieldOrigin].Focus()

	// Destination
	inputs[fieldDestination] = textinput.New()
	inputs[fieldDestination].Placeholder = "NRT, HND"
	inputs[fieldDestination].CharLimit = 50
	inputs[fieldDestination].Width = 30
	inputs[fieldDestination].Prompt = ""

	// Start Date
	inputs[fieldStartDate] = textinput.New()
	inputs[fieldStartDate].Placeholder = "2024-06-01"
	inputs[fieldStartDate].CharLimit = 10
	inputs[fieldStartDate].Width = 30
	inputs[fieldStartDate].Prompt = ""

	// End Date
	inputs[fieldEndDate] = textinput.New()
	inputs[fieldEndDate].Placeholder = "2024-06-15 (optional)"
	inputs[fieldEndDate].CharLimit = 10
	inputs[fieldEndDate].Width = 30
	inputs[fieldEndDate].Prompt = ""

	// Cabin
	inputs[fieldCabin] = textinput.New()
	inputs[fieldCabin].Placeholder = "J, F (optional)"
	inputs[fieldCabin].CharLimit = 10
	inputs[fieldCabin].Width = 30
	inputs[fieldCabin].Prompt = ""

	// Source
	inputs[fieldSource] = textinput.New()
	inputs[fieldSource].Placeholder = "aeroplan, united (optional)"
	inputs[fieldSource].CharLimit = 100
	inputs[fieldSource].Width = 30
	inputs[fieldSource].Prompt = ""

	// Pre-fill from config
	if len(cfg.PreferredAirports) > 0 {
		inputs[fieldOrigin].SetValue(strings.Join(cfg.PreferredAirports, ", "))
	}
	if len(cfg.DefaultCabins) > 0 {
		inputs[fieldCabin].SetValue(strings.Join(cfg.DefaultCabins, ", "))
	}
	if len(cfg.DefaultSources) > 0 {
		inputs[fieldSource].SetValue(strings.Join(cfg.DefaultSources, ", "))
	}

	return SearchModel{
		inputs:  inputs,
		focused: 0,
		config:  cfg,
	}
}

// Update handles input
func (m SearchModel) Update(msg tea.Msg) (SearchModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "down", "j":
			m.focused = (m.focused + 1) % numFields
			return m, m.updateFocus()

		case "shift+tab", "up", "k":
			m.focused = (m.focused - 1 + numFields) % numFields
			return m, m.updateFocus()

		case "enter":
			if m.focused == numFields-1 || msg.String() == "enter" {
				// Validate required fields
				if m.inputs[fieldOrigin].Value() != "" &&
					m.inputs[fieldDestination].Value() != "" {
					m.submitted = true
					return m, nil
				}
			}
			// Move to next field
			m.focused = (m.focused + 1) % numFields
			return m, m.updateFocus()
		}
	}

	// Update the focused input
	var cmd tea.Cmd
	m.inputs[m.focused], cmd = m.inputs[m.focused].Update(msg)
	return m, cmd
}

func (m SearchModel) updateFocus() tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))
	for i := range m.inputs {
		if i == m.focused {
			cmds[i] = m.inputs[i].Focus()
		} else {
			m.inputs[i].Blur()
		}
	}
	return tea.Batch(cmds...)
}

// View renders the search form
func (m SearchModel) View() string {
	var b strings.Builder

	// Title
	b.WriteString(TitleStyle.Render("seats.aero Search"))
	b.WriteString("\n")
	b.WriteString(SubtitleStyle.Render("Search for award flight availability"))
	b.WriteString("\n\n")

	// Form fields
	labels := []string{
		"Origin airports:",
		"Destination airports:",
		"Start date:",
		"End date:",
		"Cabin class:",
		"Mileage program:",
	}

	for i, input := range m.inputs {
		label := LabelStyle.Render(labels[i])

		var inputView string
		if i == m.focused {
			inputView = FocusedInputStyle.Render(input.View())
		} else {
			inputView = InputStyle.Render(input.View())
		}

		b.WriteString(lipgloss.JoinHorizontal(lipgloss.Left, label, inputView))
		b.WriteString("\n")
	}

	// Submit button hint
	b.WriteString("\n")
	if m.focused == numFields-1 {
		b.WriteString(FocusedButtonStyle.Render("[ Search ]"))
	} else {
		b.WriteString(ButtonStyle.Render("[ Search ]"))
	}

	// Help
	b.WriteString("\n\n")
	b.WriteString(HelpStyle.Render("Tab/j/k: navigate  |  Enter: search  |  q: quit"))
	b.WriteString("\n\n")

	// Legend
	legendBox := lipgloss.JoinVertical(
		lipgloss.Left,
		SubtitleStyle.Render("Cabin classes:"),
		lipgloss.JoinHorizontal(lipgloss.Left,
			EconomyStyle.Render("Y")+" Economy  ",
			PremiumStyle.Render("W")+" Premium  ",
			BusinessStyle.Render("J")+" Business  ",
			FirstStyle.Render("F")+" First",
		),
	)
	b.WriteString(InfoBoxStyle.Render(legendBox))

	return b.String()
}

// GetSearchParams returns the search parameters from the form
func (m SearchModel) GetSearchParams() api.SearchParams {
	return api.SearchParams{
		OriginAirports:      parseInputCSV(m.inputs[fieldOrigin].Value()),
		DestinationAirports: parseInputCSV(m.inputs[fieldDestination].Value()),
		StartDate:           strings.TrimSpace(m.inputs[fieldStartDate].Value()),
		EndDate:             strings.TrimSpace(m.inputs[fieldEndDate].Value()),
		Cabin:               strings.ToUpper(strings.TrimSpace(m.inputs[fieldCabin].Value())),
		Sources:             parseInputCSVLower(m.inputs[fieldSource].Value()),
	}
}

func parseInputCSV(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(strings.ToUpper(p))
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}

func parseInputCSVLower(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(strings.ToLower(p))
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}
