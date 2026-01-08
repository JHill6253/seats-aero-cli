package export

import (
	"encoding/csv"
	"io"
	"strconv"

	"github.com/JHill6253/seats-aero-cli/internal/api"
)

// ToCSV exports availability data as CSV
func ToCSV(w io.Writer, data []api.Availability) error {
	writer := csv.NewWriter(w)
	defer writer.Flush()

	// Write header
	header := []string{
		"ID",
		"Date",
		"Origin",
		"Destination",
		"Source",
		"Y_Available",
		"Y_Miles",
		"Y_Seats",
		"Y_Direct",
		"W_Available",
		"W_Miles",
		"W_Seats",
		"W_Direct",
		"J_Available",
		"J_Miles",
		"J_Seats",
		"J_Direct",
		"F_Available",
		"F_Miles",
		"F_Seats",
		"F_Direct",
	}
	if err := writer.Write(header); err != nil {
		return err
	}

	// Write data rows
	for _, a := range data {
		row := []string{
			a.ID,
			a.Date,
			a.Route.OriginAirport,
			a.Route.DestinationAirport,
			a.Source,
			strconv.FormatBool(a.YAvailable),
			a.YMileageCost,
			strconv.Itoa(a.YRemainingSeats),
			strconv.FormatBool(a.YDirect),
			strconv.FormatBool(a.WAvailable),
			a.WMileageCost,
			strconv.Itoa(a.WRemainingSeats),
			strconv.FormatBool(a.WDirect),
			strconv.FormatBool(a.JAvailable),
			a.JMileageCost,
			strconv.Itoa(a.JRemainingSeats),
			strconv.FormatBool(a.JDirect),
			strconv.FormatBool(a.FAvailable),
			a.FMileageCost,
			strconv.Itoa(a.FRemainingSeats),
			strconv.FormatBool(a.FDirect),
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}

// TripsToCSV exports trip data as CSV
func TripsToCSV(w io.Writer, data []api.Trip) error {
	writer := csv.NewWriter(w)
	defer writer.Flush()

	// Write header
	header := []string{
		"ID",
		"Cabin",
		"Carriers",
		"Flight_Numbers",
		"Stops",
		"Duration_Min",
		"Miles",
		"Taxes",
		"Seats",
		"Departs",
		"Arrives",
		"Source",
	}
	if err := writer.Write(header); err != nil {
		return err
	}

	// Write data rows
	for _, t := range data {
		row := []string{
			t.ID,
			t.Cabin,
			t.Carriers,
			t.FlightNumbers,
			strconv.Itoa(t.Stops),
			strconv.Itoa(t.TotalDuration),
			strconv.Itoa(t.MileageCost),
			strconv.Itoa(t.TotalTaxes),
			strconv.Itoa(t.RemainingSeats),
			t.DepartsAt.Format("2006-01-02 15:04"),
			t.ArrivesAt.Format("2006-01-02 15:04"),
			t.Source,
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}

// RoutesToCSV exports route data as CSV
func RoutesToCSV(w io.Writer, data []api.Route) error {
	writer := csv.NewWriter(w)
	defer writer.Flush()

	// Write header
	header := []string{
		"ID",
		"Origin",
		"Destination",
		"Distance",
		"Source",
	}
	if err := writer.Write(header); err != nil {
		return err
	}

	// Write data rows
	for _, r := range data {
		row := []string{
			r.ID,
			r.OriginAirport,
			r.DestinationAirport,
			strconv.Itoa(r.Distance),
			r.Source,
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}
