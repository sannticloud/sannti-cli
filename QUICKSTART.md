# Sannti CLI - Quickstart Guide

## 🚀 Getting Started

### 1. Installation
```bash
# Build from source
cd /srv/sannti-cli
make build

# Or install to system
make install
```

### 2. Configure Credentials
```bash
sannti configure
```

Enter your credentials:
- **Sannti Access Key**: your API key
- **Sannti Secret Key**: your secret key  
- **Default Region**: `br-southeast-1`

### 3. Verify Installation
```bash
# Check version
sannti version

# List regions
sannti region list

# List instances
sannti compute list
```

## 📝 Common Tasks

### Listing Resources
```bash
# List compute instances
sannti compute list

# Get detailed instance info
sannti compute get <uuid>

# List available images
sannti compute images

# List available sizes
sannti compute sizes

# List networks
sannti network list

# List IPs
sannti ip list

# List firewall rules
sannti firewall list
```

### Creating an Instance
```bash
# 1. Find an image
sannti compute images | grep "Ubuntu 24.04"

# 2. Find a size
sannti compute sizes | grep "s1.nano"

# 3. Get network UUID
sannti network list

# 4. Create instance
sannti compute create \
  --name my-server \
  --image <image-uuid> \
  --size <size-uuid> \
  --network <network-uuid>
```

### Managing Instances
```bash
# Stop instance
sannti compute stop <uuid>

# Start instance
sannti compute start <uuid>

# Delete instance
sannti compute delete <uuid>
```

## 🎨 Output Formats
```bash
# Table (default)
sannti compute list

# JSON
sannti compute list --output json

# YAML
sannti compute list --output yaml
```

## 🔧 Using Different Regions
```bash
# Override default region
sannti compute list --region br-southeast-1
```

## 🐛 Troubleshooting

### Check Configuration
```bash
cat ~/.sannti/config.yaml
```

### Use Environment Variables
```bash
export SANNTI_ACCESS_KEY="your-key"
export SANNTI_SECRET_KEY="your-secret"
export SANNTI_REGION="br-southeast-1"
```

### Verify API Access
```bash
curl -I https://console.sannti.cloud/restapi/zone/zonelist
```

## 📚 References

- [README](README.md) - Full documentation
- [ARCHITECTURE](ARCHITECTURE.md) - Technical details
- [API Docs](https://console.sannti.cloud/apidocs/swagger-ui/index.html)

---

**Tip**: Use `--output json` to debug API responses!
