package client

import (
"encoding/json"
"fmt"
"net/url"

"github.com/sannticloud/sannti-cli/internal/models"
)

// ListInstances retrieves all instances, optionally filtered by region
func (c *Client) ListInstances(regionName string) ([]models.Instance, error) {
path := "/instance/instanceList"

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

var response models.ListInstanceResponse
if err := json.Unmarshal(respBody, &response); err != nil {
return nil, fmt.Errorf("failed to parse instances response: %w", err)
}

return response.ListInstanceResponse, nil
}

// GetInstance retrieves a specific instance by UUID
func (c *Client) GetInstance(uuid, regionName string) (*models.Instance, error) {
path := fmt.Sprintf("/instance/instanceList?vmUuid=%s", url.QueryEscape(uuid))

// Add zoneUuid if region provided
if regionName != "" {
zoneUUID, err := c.GetZoneUUID(regionName)
if err != nil {
return nil, err
}
path = fmt.Sprintf("%s&zoneUuid=%s", path, url.QueryEscape(zoneUUID))
}

respBody, err := c.Get(path)
if err != nil {
return nil, err
}

var response models.ListInstanceResponse
if err := json.Unmarshal(respBody, &response); err != nil {
return nil, fmt.Errorf("failed to parse instance response: %w", err)
}

if len(response.ListInstanceResponse) == 0 {
return nil, fmt.Errorf("instance not found: %s", uuid)
}

return &response.ListInstanceResponse[0], nil
}

// CreateInstance creates a new compute instance
func (c *Client) CreateInstance(req models.CreateInstanceRequest) (*models.Instance, error) {
zoneUUID, err := c.GetZoneUUID(req.Region)
if err != nil {
return nil, err
}
req.ZoneUUID = zoneUUID

respBody, err := c.Post("/instance/createInstance", req)
if err != nil {
return nil, err
}

var response models.ListInstanceResponse
if err := json.Unmarshal(respBody, &response); err != nil {
return nil, fmt.Errorf("failed to parse create instance response: %w", err)
}

if len(response.ListInstanceResponse) == 0 {
return nil, fmt.Errorf("no instance returned from API")
}

return &response.ListInstanceResponse[0], nil
}

// StartInstance starts a stopped instance
func (c *Client) StartInstance(uuid string) error {
path := fmt.Sprintf("/instance/startInstance?uuid=%s", url.QueryEscape(uuid))

_, err := c.Get(path)
if err != nil {
return fmt.Errorf("failed to start instance: %w", err)
}

return nil
}

// StopInstance stops a running instance
func (c *Client) StopInstance(uuid string) error {
path := fmt.Sprintf("/instance/stopInstance?uuid=%s&forceStop=false", url.QueryEscape(uuid))

_, err := c.Get(path)
if err != nil {
return fmt.Errorf("failed to stop instance: %w", err)
}

return nil
}

// DeleteInstance deletes an instance
func (c *Client) DeleteInstance(uuid string) error {
path := fmt.Sprintf("/instance/destroyInstance?uuid=%s&expunge=true", url.QueryEscape(uuid))

_, err := c.Get(path)
if err != nil {
return fmt.Errorf("failed to delete instance: %w", err)
}

return nil
}

// ListComputeOfferings retrieves available compute offerings (sizes)
func (c *Client) ListComputeOfferings(regionName string) ([]models.ComputeOffering, error) {
path := "/compute/computeOfferingList"

if regionName == "" {
return nil, fmt.Errorf("region is required for listing compute offerings")
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
ListComputeOfferingResponse []models.ComputeOffering `json:"listComputeOfferingResponse"`
Count                       int                      `json:"count"`
}

if err := json.Unmarshal(respBody, &response); err != nil {
return nil, fmt.Errorf("failed to parse compute offerings response: %w", err)
}

return response.ListComputeOfferingResponse, nil
}

// ListTemplates retrieves available templates (images)
func (c *Client) ListTemplates(regionName string) ([]models.Template, error) {
path := "/template/templateList"

if regionName == "" {
return nil, fmt.Errorf("region is required for listing templates")
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
ListTemplateResponse []models.Template `json:"listTemplateResponse"`
Count                int               `json:"count"`
}

if err := json.Unmarshal(respBody, &response); err != nil {
return nil, fmt.Errorf("failed to parse templates response: %w", err)
}

return response.ListTemplateResponse, nil
}
