package export

import (
	"encoding/json"
	"io"

	"github.com/JHill6253/seats-aero-cli/internal/api"
)

// ToJSON exports availability data as JSON
func ToJSON(w io.Writer, data []api.Availability, pretty bool) error {
	encoder := json.NewEncoder(w)
	if pretty {
		encoder.SetIndent("", "  ")
	}
	return encoder.Encode(data)
}

// TripsToJSON exports trip data as JSON
func TripsToJSON(w io.Writer, data []api.Trip, pretty bool) error {
	encoder := json.NewEncoder(w)
	if pretty {
		encoder.SetIndent("", "  ")
	}
	return encoder.Encode(data)
}

// RoutesToJSON exports route data as JSON
func RoutesToJSON(w io.Writer, data []api.Route, pretty bool) error {
	encoder := json.NewEncoder(w)
	if pretty {
		encoder.SetIndent("", "  ")
	}
	return encoder.Encode(data)
}
