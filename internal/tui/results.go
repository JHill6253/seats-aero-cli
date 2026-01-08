package tui

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/JHill6253/seats-aero-cli/internal/api"
	"github.com/JHill6253/seats-aero-cli/internal/export"
)

// ResultsModel handles the results table view
type ResultsModel struct {
	table           table.Model
	results         []api.Availability
	exportRequested bool
	exportFormat    string
	width           int
	height          int
}

// NewResultsModel creates a new results view
func NewResultsModel() ResultsModel {
	columns := []table.Column{
		{Title: "Date", Width: 12},
		{Title: "From", Width: 5},
		{Title: "To", Width: 5},
		{Title: "Source", Width: 12},
		{Title: "Y", Width: 8},
		{Title: "W", Width: 8},
		{Title: "J", Width: 8},
		{Title: "F", Width: 8},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(15),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#7C3AED")).
		BorderBottom(true).
		Bold(true)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("#F9FAFB")).
		Background(lipgloss.Color("#7C3AED")).
		Bold(false)
	t.SetStyles(s)

	return ResultsModel{
		table: t,
	}
}

// SetResults populates the table with results
func (m *ResultsModel) SetResults(results []api.Availability) {
	m.results = results

	rows := make([]table.Row, len(results))
	for i, a := range results {
		rows[i] = table.Row{
			a.Date,
			a.Route.OriginAirport,
			a.Route.DestinationAirport,
			a.Source,
			formatCabin(a.YAvailable, a.YMileageCost, a.YRemainingSeats),
			formatCabin(a.WAvailable, a.WMileageCost, a.WRemainingSeats),
			formatCabin(a.JAvailable, a.JMileageCost, a.JRemainingSeats),
			formatCabin(a.FAvailable, a.FMileageCost, a.FRemainingSeats),
		}
	}

	m.table.SetRows(rows)
}

// SetSize updates the table size
func (m *ResultsModel) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.table.SetHeight(height - 5)
}

// Update handles input
func (m ResultsModel) Update(msg tea.Msg) (ResultsModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "e":
			// Export as JSON
			m.exportFormat = "json"
			m.exportRequested = true
			m.doExport()
			return m, nil
		case "c":
			// Export as CSV
			m.exportFormat = "csv"
			m.exportRequested = true
			m.doExport()
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

// View renders the results
func (m ResultsModel) View() string {
	var b strings.Builder

	// Title
	b.WriteString(TitleStyle.Render("Search Results"))
	b.WriteString("\n")
	b.WriteString(SubtitleStyle.Render(fmt.Sprintf("Found %d results", len(m.results))))
	b.WriteString("\n\n")

	// Table
	b.WriteString(m.table.View())
	b.WriteString("\n\n")

	// Selected item details
	if len(m.results) > 0 {
		idx := m.table.Cursor()
		if idx >= 0 && idx < len(m.results) {
			a := m.results[idx]
			details := m.renderDetails(a)
			b.WriteString(InfoBoxStyle.Render(details))
		}
	}

	// Help
	b.WriteString("\n")
	b.WriteString(HelpStyle.Render("j/k: navigate  |  e: export JSON  |  c: export CSV  |  Esc: back  |  q: quit"))

	return b.String()
}

func (m ResultsModel) renderDetails(a api.Availability) string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("%s -> %s on %s\n",
		a.Route.OriginAirport,
		a.Route.DestinationAirport,
		a.Date,
	))
	b.WriteString(fmt.Sprintf("Source: %s\n\n", api.SourceDisplayName(a.Source)))

	// Cabin details
	cabins := []struct {
		name      string
		available bool
		miles     string
		seats     int
		direct    bool
		airlines  string
		style     lipgloss.Style
	}{
		{"Economy", a.YAvailable, a.YMileageCost, a.YRemainingSeats, a.YDirect, a.YAirlines, EconomyStyle},
		{"Premium", a.WAvailable, a.WMileageCost, a.WRemainingSeats, a.WDirect, a.WAirlines, PremiumStyle},
		{"Business", a.JAvailable, a.JMileageCost, a.JRemainingSeats, a.JDirect, a.JAirlines, BusinessStyle},
		{"First", a.FAvailable, a.FMileageCost, a.FRemainingSeats, a.FDirect, a.FAirlines, FirstStyle},
	}

	for _, c := range cabins {
		if c.available {
			directStr := ""
			if c.direct {
				directStr = " (direct)"
			}
			b.WriteString(c.style.Render(fmt.Sprintf("%s: %s miles, %d seats%s",
				c.name, formatMiles(c.miles), c.seats, directStr)))
			if c.airlines != "" {
				b.WriteString(fmt.Sprintf(" - %s", c.airlines))
			}
			b.WriteString("\n")
		}
	}

	return b.String()
}

func (m ResultsModel) doExport() {
	if len(m.results) == 0 {
		return
	}

	filename := "seats_results"
	var err error

	switch m.exportFormat {
	case "csv":
		filename += ".csv"
		f, ferr := os.Create(filename)
		if ferr != nil {
			return
		}
		defer f.Close()
		err = export.ToCSV(f, m.results)
	default:
		filename += ".json"
		f, ferr := os.Create(filename)
		if ferr != nil {
			return
		}
		defer f.Close()
		err = export.ToJSON(f, m.results, true)
	}

	if err != nil {
		return
	}
}

func formatCabin(available bool, miles string, seats int) string {
	if !available {
		return "-"
	}
	return fmt.Sprintf("%s (%d)", formatMiles(miles), seats)
}

func formatMiles(miles string) string {
	if miles == "" || miles == "0" {
		return "?"
	}
	// Convert "50000" to "50k"
	if len(miles) > 3 {
		return miles[:len(miles)-3] + "k"
	}
	return miles
}
