package models

// Zone represents a Sannti Cloud zone/region
type Zone struct {
UUID        string `json:"uuid"`
Name        string `json:"name"`
CountryName string `json:"countryName"`
IsActive    bool   `json:"isActive"`
}

// Instance represents a compute instance
type Instance struct {
UUID                string `json:"uuid"`
Name                string `json:"name"`
DisplayName         string `json:"displayName"`
State               string `json:"state"`
ZoneName            string `json:"zoneName"`
TemplateName        string `json:"templateName"`
ServiceOfferingName string `json:"serviceOfferingName"`
Created             string `json:"created"`
MemoryMB            string `json:"memory"`
CPUNumber           int    `json:"cpuNumber"`
CPUSpeed            int    `json:"cpuSpeed"`
IPAddress           string `json:"publicIpAddress"`
CPUCore             string `json:"cpuCore"`
PrivateIP           string `json:"instancePrivateIp"`
NetworkName         string `json:"networkName"`
VolumeSize          string `json:"volumeSize"`
Status              string `json:"status"`
}

// ListInstanceResponse wraps the instance list response
type ListInstanceResponse struct {
ListInstanceResponse []Instance `json:"listInstanceResponse"`
Count                int        `json:"count"`
}

// CreateInstanceRequest represents a request to create an instance
type CreateInstanceRequest struct {
Name                string `json:"name"`
TemplateUUID        string `json:"templateUuid"`
ComputeOfferingUUID string `json:"computeOfferingUuid"`
NetworkUUID         string `json:"networkUuid"`
ZoneUUID            string `json:"zoneUuid"`
Region              string `json:"-"` // Internal field, not sent to API
RootDiskSize        int64  `json:"rootDiskSize,omitempty"`
SSHKeyName          string `json:"sshKeyName,omitempty"`
SecurityGroupName   string `json:"securitygroupName,omitempty"`
}

// ComputeOffering represents a compute size/flavor
type ComputeOffering struct {
UUID          string `json:"uuid"`
Name          string `json:"name"`
DisplayText   string `json:"displayText"`
NumberOfCores string `json:"numberOfCores"`
ClockSpeed    string `json:"clockSpeed"`
Memory        string `json:"memory"`
StorageType   string `json:"storageType"`
IsActive      bool   `json:"isActive"`
}

// Template represents an OS image
type Template struct {
UUID        string `json:"uuid"`
Name        string `json:"name"`
Description string `json:"description"`
OsTypeName  string `json:"osCategoryName"`
ZoneName    string `json:"zoneName"`
IsReady     bool   `json:"isActive"`
}

// Network represents a virtual network
type Network struct {
UUID        string `json:"uuid"`
Name        string `json:"name"`
DisplayText string `json:"displayText"`
ZoneName    string `json:"zoneName"`
State       string `json:"state"`
Cidr        string `json:"cidr"`
Gateway     string `json:"gateway"`
Type        string `json:"type"`
}

// IPAddress represents a public IP address
type IPAddress struct {
UUID        string `json:"uuid"`
IpAddress   string `json:"publicIpAddress"`
State       string `json:"state"`
ZoneName    string `json:"zoneName"`
IsStaticNat bool   `json:"isSourcenat"`
}

// FirewallRule represents a firewall rule
type FirewallRule struct {
UUID      string `json:"uuid"`
Protocol  string `json:"protocol"`
StartPort string `json:"startPort"`
EndPort   string `json:"endPort"`
CidrList  string `json:"cidrList"`
State     string `json:"status"`
}

// KubernetesVersion represents an available Kubernetes version
type KubernetesVersion struct {
UUID         string `json:"uuid"`
Name         string `json:"name"`
Description  string `json:"description"`
IsActive     bool   `json:"isActive"`
MinCPUNumber int64  `json:"minCpuNumber"`
MinMemory    int64  `json:"minMemory"`
}

// KubernetesCluster represents a Kubernetes cluster
type KubernetesCluster struct {
UUID             string `json:"uuid"`
Name             string `json:"name"`
Description      string `json:"description"`
State            string `json:"state"`
ZoneName         string `json:"zoneName"`
Size             int    `json:"size"`
ControlNodes     int    `json:"controlNodes"`
KubernetesVersion string `json:"kubernetesVersion"`
}

// CreateKubernetesRequest represents a request to create a Kubernetes cluster
type CreateKubernetesRequest struct {
Name                string `json:"name"`
Description         string `json:"description,omitempty"`
ZoneUUID            string `json:"zoneUuid"`
Region              string `json:"-"` // Internal field
KubernetesVersionUUID string `json:"kubernetesSupportedVersionUuid"`
ComputeOfferingUUID string `json:"computeOfferingUuid"`
NetworkUUID         string `json:"transNetworkUuid"`
Size                int    `json:"size"`
ControlNodes        int    `json:"controlNodes"`
HAEnabled           bool   `json:"haEnabled"`
SSHKeyName          string `json:"sshKeyName,omitempty"`
NodeRootDiskSize    int64  `json:"nodeRootDiskSize,omitempty"`
}
