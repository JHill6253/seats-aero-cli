package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/JHill6253/seats-aero-cli/internal/api"
	"github.com/JHill6253/seats-aero-cli/internal/export"
)

var tripsCmd = &cobra.Command{
	Use:   "trips <availability-id>",
	Short: "Get trip details for an availability",
	Long: `Retrieve detailed flight information for a specific availability.

The availability ID is returned from search or availability commands.

Examples:
  seats trips abc123def456
  seats trips abc123def456 --output json`,
	Args: cobra.ExactArgs(1),
	RunE: runTrips,
}

var tripsOutput string

func init() {
	rootCmd.AddCommand(tripsCmd)

	tripsCmd.Flags().StringVarP(&tripsOutput, "output", "o", "table", "Output format: table, json, csv")
}

func runTrips(cmd *cobra.Command, args []string) error {
	cfg := GetConfig()
	if cfg == nil {
		return fmt.Errorf("configuration not loaded")
	}

	if err := cfg.Validate(); err != nil {
		return err
	}

	availabilityID := args[0]
	client := api.NewClient(cfg.GetAPIKey())

	resp, err := client.GetTrips(availabilityID)
	if err != nil {
		return fmt.Errorf("get trips failed: %w", err)
	}

	switch strings.ToLower(tripsOutput) {
	case "json":
		return export.TripsToJSON(os.Stdout, resp.Data, true)
	case "csv":
		return export.TripsToCSV(os.Stdout, resp.Data)
	default:
		printTripsResults(resp.Data)
	}

	return nil
}

func printTripsResults(trips []api.Trip) {
	if len(trips) == 0 {
		fmt.Println("No trips found.")
		return
	}

	fmt.Printf("Found %d trips:\n\n", len(trips))

	for i, t := range trips {
		fmt.Printf("Trip %d: %s\n", i+1, t.Cabin)
		fmt.Printf("  Flights: %s\n", t.FlightNumbers)
		fmt.Printf("  Carriers: %s\n", t.Carriers)
		fmt.Printf("  Stops: %d\n", t.Stops)
		fmt.Printf("  Duration: %dh %dm\n", t.TotalDuration/60, t.TotalDuration%60)
		fmt.Printf("  Miles: %d\n", t.MileageCost)
		if t.TotalTaxes > 0 {
			fmt.Printf("  Taxes: %s%.2f\n", t.TaxesCurrencySymbol, float64(t.TotalTaxes)/100)
		}
		fmt.Printf("  Seats: %d\n", t.RemainingSeats)
		fmt.Printf("  Departs: %s\n", t.DepartsAt.Format("2006-01-02 15:04"))
		fmt.Printf("  Arrives: %s\n", t.ArrivesAt.Format("2006-01-02 15:04"))

		if len(t.AvailabilitySegments) > 0 {
			fmt.Println("  Segments:")
			for _, seg := range t.AvailabilitySegments {
				fmt.Printf("    %s: %s -> %s (%s %s)\n",
					seg.FlightNumber,
					seg.OriginAirport,
					seg.DestinationAirport,
					seg.AircraftCode,
					seg.DepartsAt.Format("15:04"),
				)
			}
		}
		fmt.Println()
	}
}
