# Sannti CLI - Architecture

## 🏗️ Design Principles

### 1. User-Facing Abstraction
- **Users see**: `region`, `compute`, `image`, `size`
- **API uses**: `zone`, `instance`, `template`, `offering`
- **Mapping**: Dynamic, never hardcoded

### 2. Region Abstraction Layer
The CLI maintains a clean abstraction between user-facing terminology and internal API vocabulary.
```
User Input: --region br-southeast-1
    ↓
Internal: GetZoneUUID("br-southeast-1")
    ↓
Cache: map[regionName]zoneUUID
    ↓
API Call: zoneUuid=a68cfa70-c3ac-43c8-adaf-52995aeb326e
```

**Key Features:**
- In-memory cache of region → zone UUID mapping
- Automatic refresh on cache miss
- Clear error messages with available regions
- Thread-safe with sync.RWMutex

### 3. Authentication Flow
```
Config Priority:
1. Environment Variables (SANNTI_ACCESS_KEY, SANNTI_SECRET_KEY)
2. Config File (~/.sannti/config.yaml)
3. Error if neither exists

HTTP Request:
  Headers:
    - apikey: <access-key>
    - secretkey: <secret-key>
```

## 📦 Package Structure
```
sannti-cli/
├── cmd/                    # Cobra commands
│   ├── root.go            # Base command + global flags
│   ├── configure.go       # Credential configuration
│   ├── version.go         # Version display
│   ├── region.go          # Region operations
│   ├── compute.go         # Instance management
│   ├── network.go         # Network operations
│   ├── ip.go              # IP address management
│   ├── firewall.go        # Firewall rules
│   └── kubernetes.go      # K8s version listing
│
├── internal/client/        # API client layer
│   ├── client.go          # HTTP client + auth
│   ├── zone.go            # Region mapping + cache
│   ├── compute.go         # Instance endpoints
│   ├── network.go         # Network endpoints
│   └── kubernetes.go      # K8s endpoints
│
├── internal/config/        # Configuration management
│   └── config.go          # Viper integration
│
├── internal/output/        # Output formatting
│   └── formatter.go       # Table/JSON/YAML formatters
│
└── internal/models/        # Data structures
    └── models.go          # API response models
```

## 🔄 Request Flow

### Example: `sannti compute list --region br-southeast-1`

1. **Command Parsing** (cmd/compute.go)
```go
   region := regionFlag
   cfg := config.LoadConfig()
   c := client.NewClient(cfg.AccessKey, cfg.SecretKey)
```

2. **Client Call** (internal/client/compute.go)
```go
   instances, err := c.ListInstances("br-southeast-1")
```

3. **Region Resolution** (internal/client/zone.go)
```go
   zoneUUID := c.GetZoneUUID("br-southeast-1")
   // Returns: "a68cfa70-c3ac-43c8-adaf-52995aeb326e"
```

4. **HTTP Request** (internal/client/client.go)
```go
   GET /restapi/instance/instanceList?zoneUuid=<uuid>
   Headers:
     - apikey: xxx
     - secretkey: yyy
```

5. **Response Parsing** (internal/client/compute.go)
```go
   var response models.ListInstanceResponse
   json.Unmarshal(respBody, &response)
   return response.ListInstanceResponse, nil
```

6. **Output Formatting** (internal/output/formatter.go)
```go
   output.Print(instances, FormatTable, headers, rowFunc)
```

## 📊 API Endpoint Patterns

The CLI uses consistent patterns for API endpoints:

**Instance Operations:**
- List: `GET /instance/instanceList?zoneUuid=<uuid>`
- Get: `GET /instance/instanceList?vmUuid=<uuid>&zoneUuid=<uuid>`
- Create: `POST /instance/createInstance`
- Start: `GET /instance/startInstance?uuid=<uuid>`
- Stop: `GET /instance/stopInstance?uuid=<uuid>&forceStop=false`
- Delete: `GET /instance/destroyInstance?uuid=<uuid>&expunge=true`

**Resource Discovery:**
- Images: `GET /template/templateList?zoneUuid=<uuid>`
- Sizes: `GET /compute/computeOfferingList?zoneUuid=<uuid>`
- Networks: `GET /network/networkList?zoneUuid=<uuid>`
- IPs: `GET /ipaddress/ipAddressList?zoneUuid=<uuid>`
- Firewall: `GET /firewallrule/firewallRuleList?zoneUuid=<uuid>`

**Key Observations:**
- Most endpoints use GET (including some operations like stop/delete)
- Parameter naming: `vmUuid` for specific instance, `zoneUuid` for region filter
- Response wrapper: Most return `{count, list<Resource>Response: [...]}`

## 🔐 Security Considerations

### Credential Storage
- Config file: `~/.sannti/config.yaml` with `0600` permissions
- Never logs credentials
- Hidden input for secret key during `configure`

### API Communication
- HTTPS only (https://console.sannti.cloud)
- Credentials in headers, never in URL
- 30-second timeout for requests

## 🚀 Performance Optimizations

### Region Cache
- **Problem**: Every API call that needs zoneUuid would trigger `/zone/zonelist`
- **Solution**: In-memory cache with automatic refresh
- **Thread-Safety**: sync.RWMutex for concurrent access
```go
type ZoneCache struct {
    zones map[string]models.Zone
    mu    sync.RWMutex
}
```

### Single Binary
- No external dependencies at runtime
- Static linking with `go build`
- Cross-compilation for all platforms

## 🎨 Output Design

### Table Format (Default)
```
UUID                                  NAME              STATE    VCPU  MEMORY (MB)
ecbc1ae5-bdd6-4a89-bfd4-dce1b8288416  web-prod          RUNNING  2     2048
```

**Features:**
- Clean, aligned columns
- No borders for easy parsing
- Comprehensive information (vCPU, memory, network, private IP for `get`)

### JSON Format
```json
[
  {
    "uuid": "ecbc1ae5-bdd6-4a89-bfd4-dce1b8288416",
    "name": "web-prod",
    "state": "RUNNING",
    "cpuCore": "2",
    "memory": "2048",
    "instancePrivateIp": "10.0.3.71",
    "networkName": "default-network1"
  }
]
```

**Use Cases:**
- Scripting and automation
- Piping to `jq` for filtering
- Integration with other tools

### YAML Format
```yaml
- uuid: ecbc1ae5-bdd6-4a89-bfd4-dce1b8288416
  name: web-prod
  state: RUNNING
  cpuCore: "2"
  memory: "2048"
```

**Use Cases:**
- Configuration generation
- Human-readable structured output

## 🧪 Testing Strategy

### Production Validation (v0.1.0)
All commands tested in production environment:
- ✅ Region listing and caching
- ✅ Instance lifecycle (create/start/stop/delete)
- ✅ Resource discovery (images, sizes, networks, IPs, firewall)
- ✅ Multi-format output (table/json/yaml)
- ✅ Error handling and user messages

### Future Testing
```bash
# Unit tests
go test ./internal/...

# Integration tests  
./scripts/integration-test.sh
```

## 🔮 Future Enhancements

### v0.2.0
- Kubernetes cluster management (create/list/delete)
- Volume and snapshot operations
- SSH key management
- Shell completion (bash/zsh)

### v0.3.0
- Interactive mode with prompts
- Resource templates (YAML/JSON)
- Batch operations
- Progress indicators

### v1.0.0
- Full API coverage
- Comprehensive testing suite
- Official package repositories
- Auto-update mechanism

## 📊 Error Handling

### Graceful Degradation
```go
if err != nil {
    return fmt.Errorf("failed to list instances: %w", err)
}
```

### User-Friendly Messages
```
✗ Region 'invalid-region' not found.
  Available regions: [br-southeast-1]
  Run 'sannti region list' for details
```

### HTTP Error Handling
```go
if resp.StatusCode != 200 {
    return fmt.Errorf("API error (status %d): %s",
        resp.StatusCode, string(respBody))
}
```

## 🎯 Design Decisions

### Why Go?
- Single binary distribution
- Fast compilation
- Excellent HTTP client
- Strong CLI ecosystem (Cobra)

### Why Cobra + Viper?
- Industry standard for Go CLIs
- Rich flag/config management
- Auto-generated help
- Easy subcommand structure

### Why Not SDK Generation?
- More control over user experience
- Cleaner abstractions (region vs zone)
- Simpler codebase
- Easier to maintain

### Why Manual API Integration?
- API endpoint naming inconsistencies (`vmUuid` vs `uuid`, `instanceList` vs `templateList`)
- Need for custom region abstraction layer
- Better error messages tailored to user needs
- Flexibility to work around API quirks

## 📝 Code Conventions

- Package `cmd`: Thin command handlers
- Package `client`: API communication logic
- Package `models`: Data structures only
- Package `output`: Presentation logic
- Error messages: Lowercase, actionable
- Success messages: Include resource identifiers
- Field naming: Follow API response exactly (e.g., `cpuCore` not `CPUCore`)

## 🔍 API Response Mapping

The CLI carefully maps API responses to user-friendly structures:

**Instance Model:**
```go
type Instance struct {
    UUID         string `json:"uuid"`
    Name         string `json:"name"`
    State        string `json:"state"`
    CPUCore      string `json:"cpuCore"`        // vCPU count
    MemoryMB     string `json:"memory"`         // RAM in MB
    VolumeSize   string `json:"volumeSize"`     // Disk in bytes
    NetworkName  string `json:"networkName"`    // Network name
    PrivateIP    string `json:"instancePrivateIp"` // Private IP
    Status       string `json:"status"`         // Detailed status
}
```

**Key Learnings:**
- Some fields are strings despite being numeric (API design)
- Field names vary across endpoints (requires careful mapping)
- Wrapper objects required for most list operations

---

**Philosophy**: The CLI should feel native to cloud developers, not expose internal platform details.
