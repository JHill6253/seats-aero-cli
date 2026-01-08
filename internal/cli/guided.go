package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"

	"github.com/JHill6253/seats-aero-cli/internal/api"
	"github.com/JHill6253/seats-aero-cli/internal/config"
	"github.com/JHill6253/seats-aero-cli/internal/export"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7C3AED")).
			MarginBottom(1)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6B7280"))
)

// Action represents the main menu action
type Action string

const (
	ActionSearch       Action = "search"
	ActionAvailability Action = "availability"
	ActionRoutes       Action = "routes"
	ActionTrips        Action = "trips"
	ActionExit         Action = "exit"
)

// ExportFormat represents export options
type ExportFormat string

const (
	ExportNone ExportFormat = "none"
	ExportJSON ExportFormat = "json"
	ExportCSV  ExportFormat = "csv"
)

// RunGuided runs the interactive guided CLI
func RunGuided(cfg *config.Config) error {
	fmt.Println(titleStyle.Render("seats.aero CLI"))
	fmt.Println(subtitleStyle.Render("Search for award flight availability"))
	fmt.Println()

	// Check API key first
	if err := cfg.Validate(); err != nil {
		fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("#EF4444")).Render("Error: " + err.Error()))
		fmt.Println()
		fmt.Println("Set your API key:")
		fmt.Println("  export SEATS_AERO_API_KEY=\"your-api-key\"")
		fmt.Println()
		fmt.Println("Or create a config file at ~/.config/seats-aero/config.yaml")
		return nil
	}

	for {
		var action Action

		err := huh.NewSelect[Action]().
			Title("What would you like to do?").
			Options(
				huh.NewOption("Search for flights", ActionSearch),
				huh.NewOption("View bulk availability", ActionAvailability),
				huh.NewOption("List routes", ActionRoutes),
				huh.NewOption("Get trip details", ActionTrips),
				huh.NewOption("Exit", ActionExit),
			).
			Value(&action).
			Run()

		if err != nil {
			if err == huh.ErrUserAborted {
				return nil
			}
			return err
		}

		switch action {
		case ActionSearch:
			if err := runGuidedSearch(cfg); err != nil {
				fmt.Printf("Error: %v\n\n", err)
			}
		case ActionAvailability:
			if err := runGuidedAvailability(cfg); err != nil {
				fmt.Printf("Error: %v\n\n", err)
			}
		case ActionRoutes:
			if err := runGuidedRoutes(cfg); err != nil {
				fmt.Printf("Error: %v\n\n", err)
			}
		case ActionTrips:
			if err := runGuidedTrips(cfg); err != nil {
				fmt.Printf("Error: %v\n\n", err)
			}
		case ActionExit:
			fmt.Println("Goodbye!")
			return nil
		}
	}
}

func runGuidedSearch(cfg *config.Config) error {
	var (
		origin      string
		destination string
		startDate   string
		endDate     string
		cabin       string
		source      string
	)

	// Pre-fill from config
	if len(cfg.PreferredAirports) > 0 {
		origin = strings.Join(cfg.PreferredAirports, ", ")
	}
	if len(cfg.DefaultCabins) > 0 {
		cabin = strings.Join(cfg.DefaultCabins, ", ")
	}
	if len(cfg.DefaultSources) > 0 {
		source = strings.Join(cfg.DefaultSources, ", ")
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Origin airport(s)").
				Description("Comma-separated, e.g., SFO, LAX").
				Value(&origin).
				Validate(func(s string) error {
					if strings.TrimSpace(s) == "" {
						return fmt.Errorf("origin is required")
					}
					return nil
				}),

			huh.NewInput().
				Title("Destination airport(s)").
				Description("Comma-separated, e.g., NRT, HND").
				Value(&destination).
				Validate(func(s string) error {
					if strings.TrimSpace(s) == "" {
						return fmt.Errorf("destination is required")
					}
					return nil
				}),

			huh.NewInput().
				Title("Start date").
				Description("YYYY-MM-DD format").
				Placeholder("2024-06-01").
				Value(&startDate),

			huh.NewInput().
				Title("End date (optional)").
				Description("YYYY-MM-DD format").
				Placeholder("2024-06-15").
				Value(&endDate),
		),
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Cabin class").
				Options(
					huh.NewOption("All cabins", ""),
					huh.NewOption("Economy", "economy"),
					huh.NewOption("Premium Economy", "premium"),
					huh.NewOption("Business", "business"),
					huh.NewOption("First", "first"),
				).
				Value(&cabin),

			huh.NewInput().
				Title("Mileage program (optional)").
				Description("e.g., aeroplan, united, alaska").
				Value(&source),
		),
	)

	err := form.Run()
	if err != nil {
		if err == huh.ErrUserAborted {
			return nil
		}
		return err
	}

	// Execute search
	fmt.Println("\nSearching...")

	client := api.NewClient(cfg.GetAPIKey())
	params := api.SearchParams{
		OriginAirports:      parseCSV(origin),
		DestinationAirports: parseCSV(destination),
		StartDate:           strings.TrimSpace(startDate),
		EndDate:             strings.TrimSpace(endDate),
		Cabin:               strings.TrimSpace(cabin),
		Sources:             parseCSVLower(source),
	}

	resp, err := client.Search(params)
	if err != nil {
		return fmt.Errorf("search failed: %w", err)
	}

	fmt.Println()
	printSearchResults(resp.Data)
	fmt.Println()

	// Export prompt
	if len(resp.Data) > 0 {
		if err := promptExport(resp.Data); err != nil {
			return err
		}
	}

	return nil
}

func runGuidedAvailability(cfg *config.Config) error {
	var (
		source string
		cabin  string
	)

	// Build source options
	sourceOptions := []huh.Option[string]{
		huh.NewOption("Select a program...", ""),
	}
	for _, s := range api.ValidSources() {
		sourceOptions = append(sourceOptions, huh.NewOption(api.SourceDisplayName(s), s))
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Mileage program").
				Options(sourceOptions...).
				Value(&source).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("source is required")
					}
					return nil
				}),

			huh.NewSelect[string]().
				Title("Cabin class").
				Options(
					huh.NewOption("All cabins", ""),
					huh.NewOption("Economy", "economy"),
					huh.NewOption("Premium Economy", "premium"),
					huh.NewOption("Business", "business"),
					huh.NewOption("First", "first"),
				).
				Value(&cabin),
		),
	)

	err := form.Run()
	if err != nil {
		if err == huh.ErrUserAborted {
			return nil
		}
		return err
	}

	fmt.Println("\nFetching availability...")

	client := api.NewClient(cfg.GetAPIKey())
	params := api.AvailabilityParams{
		Source: source,
		Cabin:  cabin,
	}

	resp, err := client.GetAvailability(params)
	if err != nil {
		return fmt.Errorf("get availability failed: %w", err)
	}

	fmt.Println()
	printAvailabilityResults(resp.Data)
	fmt.Println()

	if len(resp.Data) > 0 {
		if err := promptExport(resp.Data); err != nil {
			return err
		}
	}

	return nil
}

func runGuidedRoutes(cfg *config.Config) error {
	var (
		source string
		origin string
	)

	// Build source options
	sourceOptions := []huh.Option[string]{
		huh.NewOption("All programs", ""),
	}
	for _, s := range api.ValidSources() {
		sourceOptions = append(sourceOptions, huh.NewOption(api.SourceDisplayName(s), s))
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Mileage program").
				Options(sourceOptions...).
				Value(&source),

			huh.NewInput().
				Title("Origin airport (optional)").
				Description("Filter by origin, e.g., SFO").
				Value(&origin),
		),
	)

	err := form.Run()
	if err != nil {
		if err == huh.ErrUserAborted {
			return nil
		}
		return err
	}

	fmt.Println("\nFetching routes...")

	client := api.NewClient(cfg.GetAPIKey())
	params := api.RoutesParams{
		Source: source,
		Origin: strings.ToUpper(strings.TrimSpace(origin)),
	}

	resp, err := client.GetRoutes(params)
	if err != nil {
		return fmt.Errorf("get routes failed: %w", err)
	}

	fmt.Println()
	printRoutesResults(resp.Data)
	fmt.Println()

	return nil
}

func runGuidedTrips(cfg *config.Config) error {
	var availabilityID string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Availability ID").
				Description("From search or availability results").
				Value(&availabilityID).
				Validate(func(s string) error {
					if strings.TrimSpace(s) == "" {
						return fmt.Errorf("availability ID is required")
					}
					return nil
				}),
		),
	)

	err := form.Run()
	if err != nil {
		if err == huh.ErrUserAborted {
			return nil
		}
		return err
	}

	fmt.Println("\nFetching trip details...")

	client := api.NewClient(cfg.GetAPIKey())
	resp, err := client.GetTrips(strings.TrimSpace(availabilityID))
	if err != nil {
		return fmt.Errorf("get trips failed: %w", err)
	}

	fmt.Println()
	printTripsResults(resp.Data)
	fmt.Println()

	return nil
}

func promptExport(data []api.Availability) error {
	var format ExportFormat

	err := huh.NewSelect[ExportFormat]().
		Title("Export results?").
		Options(
			huh.NewOption("No", ExportNone),
			huh.NewOption("JSON", ExportJSON),
			huh.NewOption("CSV", ExportCSV),
		).
		Value(&format).
		Run()

	if err != nil {
		if err == huh.ErrUserAborted {
			return nil
		}
		return err
	}

	if format == ExportNone {
		return nil
	}

	var filename string
	err = huh.NewInput().
		Title("Filename").
		Value(&filename).
		Validate(func(s string) error {
			if strings.TrimSpace(s) == "" {
				return fmt.Errorf("filename is required")
			}
			return nil
		}).
		Run()

	if err != nil {
		if err == huh.ErrUserAborted {
			return nil
		}
		return err
	}

	// Add extension if not present
	filename = strings.TrimSpace(filename)
	switch format {
	case ExportJSON:
		if !strings.HasSuffix(filename, ".json") {
			filename += ".json"
		}
		f, err := os.Create(filename)
		if err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}
		defer f.Close()
		if err := export.ToJSON(f, data, true); err != nil {
			return fmt.Errorf("failed to export JSON: %w", err)
		}
	case ExportCSV:
		if !strings.HasSuffix(filename, ".csv") {
			filename += ".csv"
		}
		f, err := os.Create(filename)
		if err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}
		defer f.Close()
		if err := export.ToCSV(f, data); err != nil {
			return fmt.Errorf("failed to export CSV: %w", err)
		}
	}

	fmt.Printf("Exported to %s\n\n", filename)
	return nil
}

func parseCSVLower(s string) []string {
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
