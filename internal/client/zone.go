package client

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/sannticloud/sannti-cli/internal/models"
)

// ZoneCache holds the cached zone list
type ZoneCache struct {
	zones map[string]models.Zone // map[regionName]Zone
	mu    sync.RWMutex
}

var zoneCache = &ZoneCache{
	zones: make(map[string]models.Zone),
}

// ListZones retrieves all available zones (regions)
func (c *Client) ListZones() ([]models.Zone, error) {
	respBody, err := c.Get("/zone/zonelist")
	if err != nil {
		return nil, err
	}

	var response struct {
		ListZoneResponse []models.Zone `json:"listZoneResponse"`
		Count            int           `json:"count"`
	}

	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to parse zones response: %w", err)
	}

	// Update cache
	zoneCache.mu.Lock()
	for _, zone := range response.ListZoneResponse {
		if zone.IsActive {
			zoneCache.zones[zone.Name] = zone
		}
	}
	zoneCache.mu.Unlock()

	return response.ListZoneResponse, nil
}

// GetZoneUUID resolves a region name to zone UUID
// This is the critical function that maintains the region abstraction
func (c *Client) GetZoneUUID(regionName string) (string, error) {
	// Check cache first
	zoneCache.mu.RLock()
	if zone, exists := zoneCache.zones[regionName]; exists {
		zoneCache.mu.RUnlock()
		return zone.UUID, nil
	}
	zoneCache.mu.RUnlock()

	// Cache miss - fetch zones
	zones, err := c.ListZones()
	if err != nil {
		return "", fmt.Errorf("failed to fetch zones: %w", err)
	}

	// Try again after fetch
	zoneCache.mu.RLock()
	defer zoneCache.mu.RUnlock()
	
	if zone, exists := zoneCache.zones[regionName]; exists {
		return zone.UUID, nil
	}

	// Build list of available regions for error message
	var availableRegions []string
	for _, zone := range zones {
		if zone.IsActive {
			availableRegions = append(availableRegions, zone.Name)
		}
	}

	return "", fmt.Errorf("region '%s' not found. Available regions: %v. Run 'sannti region list' for details", regionName, availableRegions)
}
