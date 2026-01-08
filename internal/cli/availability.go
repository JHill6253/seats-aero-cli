package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/JHill6253/seats-aero-cli/internal/api"
	"github.com/JHill6253/seats-aero-cli/internal/export"
)

var availabilityCmd = &cobra.Command{
	Use:   "availability",
	Short: "Get bulk availability for a mileage program",
	Long: `Retrieve bulk availability data for a specific mileage program.

This is useful for exploring availability across many routes at once.

Examples:
  seats availability --source aeroplan
  seats availability --source united --cabin J,F
  seats availability --source delta --origin-region north-america --dest-region europe`,
	RunE: runAvailability,
}

var (
	availSource       string
	availCabin        string
	availOriginRegion string
	availDestRegion   string
	availStartDate    string
	availEndDate      string
	availOutput       string
)

func init() {
	rootCmd.AddCommand(availabilityCmd)

	availabilityCmd.Flags().StringVar(&availSource, "source", "", "Mileage program source (required)")
	availabilityCmd.Flags().StringVar(&availCabin, "cabin", "", "Cabin class: Y, W, J, F (comma-separated for multiple)")
	availabilityCmd.Flags().StringVar(&availOriginRegion, "origin-region", "", "Origin region filter")
	availabilityCmd.Flags().StringVar(&availDestRegion, "dest-region", "", "Destination region filter")
	availabilityCmd.Flags().StringVar(&availStartDate, "start-date", "", "Start date (YYYY-MM-DD)")
	availabilityCmd.Flags().StringVar(&availEndDate, "end-date", "", "End date (YYYY-MM-DD)")
	availabilityCmd.Flags().StringVarP(&availOutput, "output", "o", "table", "Output format: table, json, csv")

	availabilityCmd.MarkFlagRequired("source")
}

func runAvailability(cmd *cobra.Command, args []string) error {
	cfg := GetConfig()
	if cfg == nil {
		return fmt.Errorf("configuration not loaded")
	}

	if err := cfg.Validate(); err != nil {
		return err
	}

	client := api.NewClient(cfg.GetAPIKey())

	params := api.AvailabilityParams{
		Source:       strings.ToLower(availSource),
		Cabin:        strings.ToUpper(availCabin),
		OriginRegion: availOriginRegion,
		DestRegion:   availDestRegion,
		StartDate:    availStartDate,
		EndDate:      availEndDate,
	}

	resp, err := client.GetAvailability(params)
	if err != nil {
		return fmt.Errorf("get availability failed: %w", err)
	}

	switch strings.ToLower(availOutput) {
	case "json":
		return export.ToJSON(os.Stdout, resp.Data, true)
	case "csv":
		return export.ToCSV(os.Stdout, resp.Data)
	default:
		printAvailabilityResults(resp.Data)
	}

	return nil
}

func printAvailabilityResults(results []api.Availability) {
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
