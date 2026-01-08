package api

import (
	"fmt"
	"strconv"
)

// GetAvailability retrieves bulk availability for a mileage program
func (c *Client) GetAvailability(params AvailabilityParams) (*AvailabilityResponse, error) {
	queryParams := make(map[string]string)

	// Source is required for bulk availability
	if params.Source != "" {
		queryParams["source"] = params.Source
	}

	// Cabin filter
	if params.Cabin != "" {
		queryParams["cabin"] = params.Cabin
	}

	// Region filters
	if params.OriginRegion != "" {
		queryParams["origin_region"] = params.OriginRegion
	}
	if params.DestRegion != "" {
		queryParams["destination_region"] = params.DestRegion
	}

	// Date range
	if params.StartDate != "" {
		queryParams["start_date"] = params.StartDate
	}
	if params.EndDate != "" {
		queryParams["end_date"] = params.EndDate
	}

	// Pagination
	if params.Take > 0 {
		queryParams["take"] = strconv.Itoa(params.Take)
	}
	if params.Skip > 0 {
		queryParams["skip"] = strconv.Itoa(params.Skip)
	}
	if params.Cursor > 0 {
		queryParams["cursor"] = strconv.FormatInt(params.Cursor, 10)
	}

	var response AvailabilityResponse
	if err := c.get("/availability", queryParams, &response); err != nil {
		return nil, fmt.Errorf("get availability failed: %w", err)
	}

	return &response, nil
}

// GetAvailabilityAll retrieves all availability results, handling pagination
func (c *Client) GetAvailabilityAll(params AvailabilityParams) ([]Availability, error) {
	var allResults []Availability

	params.Take = 100 // Max per page
	params.Skip = 0

	for {
		resp, err := c.GetAvailability(params)
		if err != nil {
			return nil, err
		}

		allResults = append(allResults, resp.Data...)

		if !resp.HasMore || len(resp.Data) == 0 {
			break
		}

		params.Skip += len(resp.Data)
		if params.Cursor == 0 {
			params.Cursor = resp.Cursor
		}
	}

	return allResults, nil
}
