package tui

import (
	"github.com/charmbracelet/lipgloss"
)

// Colors
var (
	primaryColor   = lipgloss.Color("#7C3AED") // Purple
	secondaryColor = lipgloss.Color("#06B6D4") // Cyan
	successColor   = lipgloss.Color("#10B981") // Green
	warningColor   = lipgloss.Color("#F59E0B") // Amber
	errorColor     = lipgloss.Color("#EF4444") // Red
	mutedColor     = lipgloss.Color("#6B7280") // Gray
	textColor      = lipgloss.Color("#F9FAFB") // White
	bgColor        = lipgloss.Color("#1F2937") // Dark gray
)

// Styles
var (
	// App styles
	AppStyle = lipgloss.NewStyle().
			Padding(1, 2)

	// Title styles
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(primaryColor).
			MarginBottom(1)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			MarginBottom(1)

	// Form styles
	LabelStyle = lipgloss.NewStyle().
			Foreground(textColor).
			Width(20)

	InputStyle = lipgloss.NewStyle().
			Foreground(textColor).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(mutedColor).
			Padding(0, 1)

	FocusedInputStyle = lipgloss.NewStyle().
				Foreground(textColor).
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(primaryColor).
				Padding(0, 1)

	// Button styles
	ButtonStyle = lipgloss.NewStyle().
			Foreground(textColor).
			Background(mutedColor).
			Padding(0, 2).
			MarginRight(1)

	FocusedButtonStyle = lipgloss.NewStyle().
				Foreground(textColor).
				Background(primaryColor).
				Padding(0, 2).
				MarginRight(1)

	// Table styles
	TableHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(primaryColor).
				BorderStyle(lipgloss.NormalBorder()).
				BorderBottom(true).
				BorderForeground(mutedColor)

	TableRowStyle = lipgloss.NewStyle().
			Foreground(textColor)

	SelectedRowStyle = lipgloss.NewStyle().
				Foreground(textColor).
				Background(primaryColor)

	// Status styles
	SuccessStyle = lipgloss.NewStyle().
			Foreground(successColor)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(errorColor)

	WarningStyle = lipgloss.NewStyle().
			Foreground(warningColor)

	// Help styles
	HelpStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			MarginTop(1)

	// Cabin class styles
	EconomyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9CA3AF"))

	PremiumStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#60A5FA"))

	BusinessStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#A78BFA"))

	FirstStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FBBF24"))

	// Availability indicator styles
	AvailableStyle = lipgloss.NewStyle().
			Foreground(successColor)

	UnavailableStyle = lipgloss.NewStyle().
				Foreground(mutedColor)

	// Info box
	InfoBoxStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(secondaryColor).
			Padding(1, 2).
			MarginTop(1)

	// Spinner
	SpinnerStyle = lipgloss.NewStyle().
			Foreground(primaryColor)
)

// CabinStyle returns the appropriate style for a cabin class
func CabinStyle(cabin string) lipgloss.Style {
	switch cabin {
	case "Y":
		return EconomyStyle
	case "W":
		return PremiumStyle
	case "J":
		return BusinessStyle
	case "F":
		return FirstStyle
	default:
		return lipgloss.NewStyle()
	}
}
