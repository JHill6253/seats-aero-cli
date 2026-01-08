package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/JHill6253/seats-aero-cli/internal/config"
	"github.com/JHill6253/seats-aero-cli/internal/tui"
)

var (
	cfgFile string
	cfg     *config.Config
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "seats",
	Short: "CLI and TUI for seats.aero award flight search",
	Long: `seats is a command-line interface and terminal UI for searching
award flight availability using the seats.aero API.

Run without arguments to launch the interactive TUI, or use subcommands
for direct CLI access.

Examples:
  seats                                    # Launch TUI
  seats search --from SFO --to NRT         # Search flights
  seats availability --source aeroplan     # Bulk availability
  seats routes --source united             # List routes`,
	Run: func(cmd *cobra.Command, args []string) {
		// Default behavior: launch TUI
		if err := runTUI(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/seats-aero/config.yaml)")
}

func initConfig() {
	var err error
	cfg, err = config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: %v\n", err)
	}
}

// GetConfig returns the loaded configuration
func GetConfig() *config.Config {
	return cfg
}

func runTUI() error {
	if cfg == nil {
		var err error
		cfg, err = config.Load()
		if err != nil {
			return err
		}
	}

	return tui.Run(cfg)
}

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
	Long:  `View and manage seats CLI configuration.`,
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	Run: func(cmd *cobra.Command, args []string) {
		if cfg == nil {
			fmt.Println("No configuration loaded")
			return
		}

		fmt.Println("Configuration:")
		fmt.Printf("  API Key: %s\n", maskAPIKey(cfg.GetAPIKey()))
		fmt.Printf("  Default Sources: %v\n", cfg.DefaultSources)
		fmt.Printf("  Default Cabins: %v\n", cfg.DefaultCabins)
		fmt.Printf("  Preferred Airports: %v\n", cfg.PreferredAirports)

		if path, err := config.ConfigPath(); err == nil {
			fmt.Printf("\nConfig file path: %s\n", path)
		}
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configShowCmd)
}

func maskAPIKey(key string) string {
	if key == "" {
		return "(not set)"
	}
	if len(key) <= 8 {
		return "****"
	}
	return key[:4] + "..." + key[len(key)-4:]
}
