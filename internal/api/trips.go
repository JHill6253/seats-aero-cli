package api

import (
	"fmt"
)

// GetTrips retrieves trip details for an availability ID
func (c *Client) GetTrips(availabilityID string) (*TripsResponse, error) {
	endpoint := fmt.Sprintf("/trips/%s", availabilityID)

	var response TripsResponse
	if err := c.get(endpoint, nil, &response); err != nil {
		return nil, fmt.Errorf("get trips failed: %w", err)
	}

	return &response, nil
}
