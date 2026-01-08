package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/JHill6253/seats-aero-cli/internal/api"
	"github.com/JHill6253/seats-aero-cli/internal/export"
)

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for award flight availability",
	Long: `Search for cached award flight availability between airports.

Examples:
  seats search --from SFO --to NRT --start-date 2024-06-01
  seats search --from SFO,LAX --to NRT,HND --cabin J --source united,aeroplan
  seats search --from SFO --to NRT --start-date 2024-06-01 --end-date 2024-06-15 --output json
  seats search --from SFO --to NRT --direct-only --output csv > results.csv`,
	RunE: runSearch,
}

var (
	searchFrom      string
	searchTo        string
	searchStartDate string
	searchEndDate   string
	searchCabin     string
	searchSource    string
	searchDirect    bool
	searchOutput    string
)

func init() {
	rootCmd.AddCommand(searchCmd)

	searchCmd.Flags().StringVar(&searchFrom, "from", "", "Origin airport(s), comma-separated (required)")
	searchCmd.Flags().StringVar(&searchTo, "to", "", "Destination airport(s), comma-separated (required)")
	searchCmd.Flags().StringVar(&searchStartDate, "start-date", "", "Start date (YYYY-MM-DD)")
	searchCmd.Flags().StringVar(&searchEndDate, "end-date", "", "End date (YYYY-MM-DD)")
	searchCmd.Flags().StringVar(&searchCabin, "cabin", "", "Cabin class: Y/economy, W/premium, J/business, F/first")
	searchCmd.Flags().StringVar(&searchSource, "source", "", "Mileage program source(s), comma-separated")
	searchCmd.Flags().BoolVar(&searchDirect, "direct-only", false, "Only show direct flights")
	searchCmd.Flags().StringVarP(&searchOutput, "output", "o", "table", "Output format: table, json, csv")

	searchCmd.MarkFlagRequired("from")
	searchCmd.MarkFlagRequired("to")
}

func runSearch(cmd *cobra.Command, args []string) error {
	cfg := GetConfig()
	if cfg == nil {
		return fmt.Errorf("configuration not loaded")
	}

	if err := cfg.Validate(); err != nil {
		return err
	}

	client := api.NewClient(cfg.GetAPIKey())

	params := api.SearchParams{
		OriginAirports:      parseCSV(searchFrom),
		DestinationAirports: parseCSV(searchTo),
		StartDate:           searchStartDate,
		EndDate:             searchEndDate,
		Cabin:               cabinCodeToName(strings.ToUpper(searchCabin)),
		Sources:             parseCSV(searchSource),
		DirectOnly:          searchDirect,
	}

	resp, err := client.Search(params)
	if err != nil {
		return fmt.Errorf("search failed: %w", err)
	}

	switch strings.ToLower(searchOutput) {
	case "json":
		return export.ToJSON(os.Stdout, resp.Data, true)
	case "csv":
		return export.ToCSV(os.Stdout, resp.Data)
	default:
		printSearchResults(resp.Data)
	}

	return nil
}

func parseCSV(s string) []string {
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

func printSearchResults(results []api.Availability) {
	if len(results) == 0 {
		fmt.Println("No results found.")
		return
	}

	fmt.Printf("Found %d results:\n\n", len(results))

	// Print header
	fmt.Printf("%-12s %-5s %-5s %-15s %-8s %-8s %-8s %-8s\n",
		"Date", "From", "To", "Source", "Y", "W", "J", "F")
	fmt.Println(strings.Repeat("-", 80))

	for _, a := range results {
		yInfo := formatCabinInfo(a.YAvailable, a.YMileageCost, a.YRemainingSeats)
		wInfo := formatCabinInfo(a.WAvailable, a.WMileageCost, a.WRemainingSeats)
		jInfo := formatCabinInfo(a.JAvailable, a.JMileageCost, a.JRemainingSeats)
		fInfo := formatCabinInfo(a.FAvailable, a.FMileageCost, a.FRemainingSeats)

		fmt.Printf("%-12s %-5s %-5s %-15s %-8s %-8s %-8s %-8s\n",
			a.Date,
			a.Route.OriginAirport,
			a.Route.DestinationAirport,
			a.Source,
			yInfo,
			wInfo,
			jInfo,
			fInfo,
		)
	}
}

func formatCabinInfo(available bool, miles string, seats int) string {
	if !available {
		return "-"
	}
	if miles == "" || miles == "0" {
		return fmt.Sprintf("(%d)", seats)
	}
	return fmt.Sprintf("%sk", miles[:len(miles)-3])
}

// cabinCodeToName converts cabin codes (Y/W/J/F) to API names (economy/premium/business/first)
// If already a valid name or empty, returns as-is (lowercased)
func cabinCodeToName(code string) string {
	code = strings.TrimSpace(code)
	if code == "" {
		return ""
	}
	switch strings.ToUpper(code) {
	case "Y", "ECONOMY":
		return "economy"
	case "W", "PREMIUM":
		return "premium"
	case "J", "BUSINESS":
		return "business"
	case "F", "FIRST":
		return "first"
	default:
		return strings.ToLower(code)
	}
}
