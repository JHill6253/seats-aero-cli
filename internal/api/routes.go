package api

import (
	"fmt"
)

// GetRoutes retrieves available routes for a source
func (c *Client) GetRoutes(params RoutesParams) (*RoutesResponse, error) {
	queryParams := make(map[string]string)

	if params.Source != "" {
		queryParams["source"] = params.Source
	}

	if params.Origin != "" {
		queryParams["origin"] = params.Origin
	}

	var response RoutesResponse
	if err := c.get("/routes", queryParams, &response); err != nil {
		return nil, fmt.Errorf("get routes failed: %w", err)
	}

	return &response, nil
}
