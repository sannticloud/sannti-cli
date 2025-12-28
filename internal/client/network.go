package client

import (
"encoding/json"
"fmt"
"net/url"

"github.com/sannticloud/sannti-cli/internal/models"
)

// ListNetworks retrieves all networks, optionally filtered by region
func (c *Client) ListNetworks(regionName string) ([]models.Network, error) {
path := "/network/networkList"

if regionName != "" {
zoneUUID, err := c.GetZoneUUID(regionName)
if err != nil {
return nil, err
}
path = fmt.Sprintf("%s?zoneUuid=%s", path, url.QueryEscape(zoneUUID))
}

respBody, err := c.Get(path)
if err != nil {
return nil, err
}

var response struct {
ListNetworkResponse []models.Network `json:"listNetworkResponse"`
Count               int              `json:"count"`
}

if err := json.Unmarshal(respBody, &response); err != nil {
return nil, fmt.Errorf("failed to parse networks response: %w", err)
}

return response.ListNetworkResponse, nil
}

// ListIPAddresses retrieves IP addresses for a region
func (c *Client) ListIPAddresses(regionName string) ([]models.IPAddress, error) {
path := "/ipaddress/ipAddressList"

if regionName == "" {
return nil, fmt.Errorf("region is required for listing IP addresses")
}

zoneUUID, err := c.GetZoneUUID(regionName)
if err != nil {
return nil, err
}
path = fmt.Sprintf("%s?zoneUuid=%s", path, url.QueryEscape(zoneUUID))

respBody, err := c.Get(path)
if err != nil {
return nil, err
}

var response struct {
ListIpAddressResponse []models.IPAddress `json:"listIpAddressResponse"`
Count                 int                `json:"count"`
}

if err := json.Unmarshal(respBody, &response); err != nil {
return nil, fmt.Errorf("failed to parse IP addresses response: %w", err)
}

return response.ListIpAddressResponse, nil
}

// ListFirewallRules retrieves firewall rules
func (c *Client) ListFirewallRules(regionName string) ([]models.FirewallRule, error) {
path := "/firewallrule/firewallRuleList"

if regionName != "" {
zoneUUID, err := c.GetZoneUUID(regionName)
if err != nil {
return nil, err
}
path = fmt.Sprintf("%s?zoneUuid=%s", path, url.QueryEscape(zoneUUID))
}

respBody, err := c.Get(path)
if err != nil {
return nil, err
}

var response struct {
ListFirewallRuleResponse []models.FirewallRule `json:"listFirewallRuleResponse"`
Count                    int                   `json:"count"`
}

if err := json.Unmarshal(respBody, &response); err != nil {
return nil, fmt.Errorf("failed to parse firewall rules response: %w", err)
}

return response.ListFirewallRuleResponse, nil
}
