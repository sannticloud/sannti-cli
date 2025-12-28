package client

import (
"encoding/json"
"fmt"
"net/url"

"github.com/sannticloud/sannti-cli/internal/models"
)

// ListKubernetesVersions retrieves available Kubernetes versions
func (c *Client) ListKubernetesVersions(regionName string) ([]models.KubernetesVersion, error) {
path := "/costestimate/kubernetes-version-list"

// Resolve region to zone UUID (obrigatório para este endpoint)
if regionName == "" {
return nil, fmt.Errorf("region is required for listing Kubernetes versions")
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
ListKubernetesVersion []models.KubernetesVersion `json:"listKubernetesVersion"`
Count                 int                        `json:"count"`
}

if err := json.Unmarshal(respBody, &response); err != nil {
return nil, fmt.Errorf("failed to parse kubernetes versions response: %w", err)
}

return response.ListKubernetesVersion, nil
}

// ListKubernetesClusters retrieves all Kubernetes clusters
func (c *Client) ListKubernetesClusters(clusterUUID string) ([]models.KubernetesCluster, error) {
path := "/kubernetes/listCluster"

if clusterUUID != "" {
path = fmt.Sprintf("%s?clusterUuid=%s", path, url.QueryEscape(clusterUUID))
}

respBody, err := c.Get(path)
if err != nil {
return nil, err
}

var clusters []models.KubernetesCluster
if err := json.Unmarshal(respBody, &clusters); err != nil {
return nil, fmt.Errorf("failed to parse kubernetes clusters response: %w", err)
}

return clusters, nil
}

// CreateKubernetesCluster creates a new Kubernetes cluster
func (c *Client) CreateKubernetesCluster(req models.CreateKubernetesRequest) (*models.KubernetesCluster, error) {
zoneUUID, err := c.GetZoneUUID(req.Region)
if err != nil {
return nil, err
}
req.ZoneUUID = zoneUUID

respBody, err := c.Post("/kubernetes/createKubernetes", req)
if err != nil {
return nil, err
}

var cluster models.KubernetesCluster
if err := json.Unmarshal(respBody, &cluster); err != nil {
return nil, fmt.Errorf("failed to parse create kubernetes cluster response: %w", err)
}

return &cluster, nil
}

// DeleteKubernetesCluster deletes a Kubernetes cluster
func (c *Client) DeleteKubernetesCluster(clusterUUID string) error {
path := fmt.Sprintf("/kubernetes/destroyKubernetes?clusterUuid=%s", url.QueryEscape(clusterUUID))

_, err := c.Delete(path)
if err != nil {
return fmt.Errorf("failed to delete kubernetes cluster: %w", err)
}

return nil
}
