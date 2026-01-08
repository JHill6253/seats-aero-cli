package api

import (
	"fmt"
	"strconv"
	"strings"
)

// Search performs a cached search for availability
func (c *Client) Search(params SearchParams) (*SearchResponse, error) {
	queryParams := make(map[string]string)

	// Origin airports (comma-separated)
	if len(params.OriginAirports) > 0 {
		queryParams["origin_airport"] = strings.Join(params.OriginAirports, ",")
	}

	// Destination airports (comma-separated)
	if len(params.DestinationAirports) > 0 {
		queryParams["destination_airport"] = strings.Join(params.DestinationAirports, ",")
	}

	// Date range
	if params.StartDate != "" {
		queryParams["start_date"] = params.StartDate
	}
	if params.EndDate != "" {
		queryParams["end_date"] = params.EndDate
	}

	// Cabin filter
	if params.Cabin != "" {
		queryParams["cabin"] = params.Cabin
	}

	// Source filter (comma-separated)
	if len(params.Sources) > 0 {
		queryParams["source"] = strings.Join(params.Sources, ",")
	}

	// Direct flights only
	if params.DirectOnly {
		queryParams["direct"] = "true"
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

	var response SearchResponse
	if err := c.get("/search", queryParams, &response); err != nil {
		return nil, fmt.Errorf("search failed: %w", err)
	}

	return &response, nil
}

// SearchAll retrieves all results from a search, handling pagination
func (c *Client) SearchAll(params SearchParams) ([]Availability, error) {
	var allResults []Availability

	params.Take = 100 // Max per page
	params.Skip = 0

	for {
		resp, err := c.Search(params)
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
