package api

import "time"

// Route represents a flight route
type Route struct {
	ID                 string `json:"ID"`
	OriginAirport      string `json:"OriginAirport"`
	DestinationAirport string `json:"DestinationAirport"`
	NumDaysOut         int    `json:"NumDaysOut"`
	Distance           int    `json:"Distance"`
	Source             string `json:"Source"`
}

// Availability represents summarized availability for a route on a given date
type Availability struct {
	ID         string    `json:"ID"`
	RouteID    string    `json:"RouteID"`
	Route      Route     `json:"Route"`
	Date       string    `json:"Date"`
	ParsedDate time.Time `json:"-"`

	// Economy
	YAvailable      bool   `json:"YAvailable"`
	YMileageCost    string `json:"YMileageCost"`
	YRemainingSeats int    `json:"YRemainingSeats"`
	YAirlines       string `json:"YAirlines"`
	YDirect         bool   `json:"YDirect"`

	// Premium Economy
	WAvailable      bool   `json:"WAvailable"`
	WMileageCost    string `json:"WMileageCost"`
	WRemainingSeats int    `json:"WRemainingSeats"`
	WAirlines       string `json:"WAirlines"`
	WDirect         bool   `json:"WDirect"`

	// Business
	JAvailable      bool   `json:"JAvailable"`
	JMileageCost    string `json:"JMileageCost"`
	JRemainingSeats int    `json:"JRemainingSeats"`
	JAirlines       string `json:"JAirlines"`
	JDirect         bool   `json:"JDirect"`

	// First
	FAvailable      bool   `json:"FAvailable"`
	FMileageCost    string `json:"FMileageCost"`
	FRemainingSeats int    `json:"FRemainingSeats"`
	FAirlines       string `json:"FAirlines"`
	FDirect         bool   `json:"FDirect"`

	Source    string    `json:"Source"`
	CreatedAt time.Time `json:"CreatedAt"`
	UpdatedAt time.Time `json:"UpdatedAt"`
}

// AvailabilitySegment represents a single flight segment
type AvailabilitySegment struct {
	ID                 string    `json:"ID"`
	RouteID            string    `json:"RouteID"`
	AvailabilityID     string    `json:"AvailabilityID"`
	AvailabilityTripID string    `json:"AvailabilityTripID"`
	FlightNumber       string    `json:"FlightNumber"`
	Distance           int       `json:"Distance"`
	FareClass          string    `json:"FareClass"`
	AircraftName       string    `json:"AircraftName"`
	AircraftCode       string    `json:"AircraftCode"`
	OriginAirport      string    `json:"OriginAirport"`
	DestinationAirport string    `json:"DestinationAirport"`
	DepartsAt          time.Time `json:"DepartsAt"`
	ArrivesAt          time.Time `json:"ArrivesAt"`
	Source             string    `json:"Source"`
	Order              int       `json:"Order"`
}

// Trip represents a specific itinerary with flight details
type Trip struct {
	ID                   string                `json:"ID"`
	RouteID              string                `json:"RouteID"`
	AvailabilityID       string                `json:"AvailabilityID"`
	AvailabilitySegments []AvailabilitySegment `json:"AvailabilitySegments"`
	TotalDuration        int                   `json:"TotalDuration"`
	Stops                int                   `json:"Stops"`
	Carriers             string                `json:"Carriers"`
	RemainingSeats       int                   `json:"RemainingSeats"`
	MileageCost          int                   `json:"MileageCost"`
	TotalTaxes           int                   `json:"TotalTaxes"`
	TaxesCurrency        string                `json:"TaxesCurrency"`
	TaxesCurrencySymbol  string                `json:"TaxesCurrencySymbol"`
	FlightNumbers        string                `json:"FlightNumbers"`
	DepartsAt            time.Time             `json:"DepartsAt"`
	ArrivesAt            time.Time             `json:"ArrivesAt"`
	Cabin                string                `json:"Cabin"`
	Source               string                `json:"Source"`
	CreatedAt            time.Time             `json:"CreatedAt"`
	UpdatedAt            time.Time             `json:"UpdatedAt"`
}

// SearchResponse represents the response from the cached search endpoint
type SearchResponse struct {
	Data    []Availability `json:"data"`
	Count   int            `json:"count"`
	Cursor  int64          `json:"cursor"`
	HasMore bool           `json:"hasMore"`
}

// AvailabilityResponse represents the response from the bulk availability endpoint
type AvailabilityResponse struct {
	Data    []Availability `json:"data"`
	Count   int            `json:"count"`
	Cursor  int64          `json:"cursor"`
	HasMore bool           `json:"hasMore"`
}

// TripsResponse represents the response from the trips endpoint
type TripsResponse struct {
	Data  []Trip `json:"data"`
	Count int    `json:"count"`
}

// RoutesResponse represents the response from the routes endpoint
type RoutesResponse struct {
	Data  []Route `json:"data"`
	Count int     `json:"count"`
}

// SearchParams contains parameters for cached search
type SearchParams struct {
	OriginAirports      []string
	DestinationAirports []string
	StartDate           string
	EndDate             string
	Cabin               string   // Y, W, J, F or empty for all
	Sources             []string // Empty for all sources
	DirectOnly          bool
	Take                int
	Skip                int
	Cursor              int64
}

// AvailabilityParams contains parameters for bulk availability
type AvailabilityParams struct {
	Source       string
	Cabin        string
	OriginRegion string
	DestRegion   string
	StartDate    string
	EndDate      string
	Take         int
	Skip         int
	Cursor       int64
}

// RoutesParams contains parameters for routes endpoint
type RoutesParams struct {
	Source string
	Origin string
}

// CabinClass represents cabin class codes
type CabinClass string

const (
	CabinEconomy        CabinClass = "Y"
	CabinPremiumEconomy CabinClass = "W"
	CabinBusiness       CabinClass = "J"
	CabinFirst          CabinClass = "F"
)

// ValidCabins returns all valid cabin codes
func ValidCabins() []string {
	return []string{"Y", "W", "J", "F"}
}

// ValidSources returns all valid mileage program sources
func ValidSources() []string {
	return []string{
		"eurobonus",
		"virginatlantic",
		"aeromexico",
		"american",
		"delta",
		"etihad",
		"united",
		"emirates",
		"aeroplan",
		"alaska",
		"velocity",
		"qantas",
		"connectmiles",
		"azul",
		"smiles",
		"flyingblue",
		"jetblue",
		"qatar",
		"turkish",
		"singapore",
		"ethiopian",
		"saudia",
		"finnair",
		"lufthansa",
	}
}

// SourceDisplayName returns a human-readable name for a source
func SourceDisplayName(source string) string {
	names := map[string]string{
		"eurobonus":      "SAS EuroBonus",
		"virginatlantic": "Virgin Atlantic",
		"aeromexico":     "Aeromexico",
		"american":       "American Airlines",
		"delta":          "Delta SkyMiles",
		"etihad":         "Etihad Guest",
		"united":         "United MileagePlus",
		"emirates":       "Emirates Skywards",
		"aeroplan":       "Air Canada Aeroplan",
		"alaska":         "Alaska Mileage Plan",
		"velocity":       "Virgin Australia",
		"qantas":         "Qantas",
		"connectmiles":   "Copa ConnectMiles",
		"azul":           "Azul TudoAzul",
		"smiles":         "GOL Smiles",
		"flyingblue":     "Flying Blue",
		"jetblue":        "JetBlue TrueBlue",
		"qatar":          "Qatar Privilege Club",
		"turkish":        "Turkish Miles&Smiles",
		"singapore":      "Singapore KrisFlyer",
		"ethiopian":      "Ethiopian ShebaMiles",
		"saudia":         "Saudi AlFursan",
		"finnair":        "Finnair Plus",
		"lufthansa":      "Lufthansa Miles&More",
	}
	if name, ok := names[source]; ok {
		return name
	}
	return source
}

// CabinDisplayName returns a human-readable name for a cabin class
func CabinDisplayName(cabin string) string {
	names := map[string]string{
		"Y": "Economy",
		"W": "Premium Economy",
		"J": "Business",
		"F": "First",
	}
	if name, ok := names[cabin]; ok {
		return name
	}
	return cabin
}
