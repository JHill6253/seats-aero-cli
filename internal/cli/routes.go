package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/JHill6253/seats-aero-cli/internal/api"
	"github.com/JHill6253/seats-aero-cli/internal/export"
)

var routesCmd = &cobra.Command{
	Use:   "routes",
	Short: "List available routes",
	Long: `List routes available for a mileage program.

Examples:
  seats routes --source aeroplan
  seats routes --source united --origin SFO`,
	RunE: runRoutes,
}

var (
	routesSource string
	routesOrigin string
	routesOutput string
)

func init() {
	rootCmd.AddCommand(routesCmd)

	routesCmd.Flags().StringVar(&routesSource, "source", "", "Mileage program source")
	routesCmd.Flags().StringVar(&routesOrigin, "origin", "", "Filter by origin airport")
	routesCmd.Flags().StringVarP(&routesOutput, "output", "o", "table", "Output format: table, json, csv")
}

func runRoutes(cmd *cobra.Command, args []string) error {
	cfg := GetConfig()
	if cfg == nil {
		return fmt.Errorf("configuration not loaded")
	}

	if err := cfg.Validate(); err != nil {
		return err
	}

	client := api.NewClient(cfg.GetAPIKey())

	params := api.RoutesParams{
		Source: strings.ToLower(routesSource),
		Origin: strings.ToUpper(routesOrigin),
	}

	resp, err := client.GetRoutes(params)
	if err != nil {
		return fmt.Errorf("get routes failed: %w", err)
	}

	switch strings.ToLower(routesOutput) {
	case "json":
		return export.RoutesToJSON(os.Stdout, resp.Data, true)
	case "csv":
		return export.RoutesToCSV(os.Stdout, resp.Data)
	default:
		printRoutesResults(resp.Data)
	}

	return nil
}

func printRoutesResults(routes []api.Route) {
	if len(routes) == 0 {
		fmt.Println("No routes found.")
		return
	}

	fmt.Printf("Found %d routes:\n\n", len(routes))

	// Print header
	fmt.Printf("%-5s %-5s %-8s %-20s\n", "From", "To", "Distance", "Source")
	fmt.Println(strings.Repeat("-", 45))

	for _, r := range routes {
		fmt.Printf("%-5s %-5s %-8d %-20s\n",
			r.OriginAirport,
			r.DestinationAirport,
			r.Distance,
			api.SourceDisplayName(r.Source),
		)
	}
}
