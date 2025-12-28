# Sannti CLI

> **Beyond Cloud. Without Barriers.**

Official command-line interface for [Sannti Cloud](https://sannti.com) - Brazilian cloud infrastructure with a focus on digital sovereignty.

[![Version](https://img.shields.io/badge/version-v0.1.0-blue)](https://github.com/sannticloud/sannti-cli)
[![License](https://img.shields.io/badge/license-Apache%202.0-green)](LICENSE)

## 🚀 Quick Start

### Installation

Choose your preferred installation method:

#### Option 1: Install Script (Recommended)
```bash
curl -fsSL https://get.sannti.cloud/install.sh | bash
```

> **Coming soon** - Script will be available at https://get.sannti.cloud/install.sh

#### Option 2: Homebrew (macOS/Linux)
```bash
brew install sannticloud/tap/sannti
```

> **Coming soon** - Homebrew formula will be available in v0.2.0

#### Option 3: NPM
```bash
npm install -g @sannticloud/sannti-cli
```

> **Coming soon** - NPM package will be available in v0.2.0

#### Option 4: Build from Source

**Requirements:** Go 1.21+
```bash
# Clone the repository
git clone https://github.com/sannticloud/sannti-cli.git
cd sannti-cli

# Build and install
make install

# Verify installation
sannti version
```

> ✅ **Available now** - Best option for v0.1.0

#### Option 5: Download Pre-built Binary

Download the latest release for your platform:

**Linux (x86_64)**:
```bash
curl -L https://github.com/sannticloud/sannti-cli/releases/download/v0.1.0/sannti-v0.1.0-linux-amd64 -o sannti
chmod +x sannti
sudo mv sannti /usr/local/bin/
```

**macOS (Apple Silicon)**:
```bash
curl -L https://github.com/sannticloud/sannti-cli/releases/download/v0.1.0/sannti-v0.1.0-darwin-arm64 -o sannti
chmod +x sannti
sudo mv sannti /usr/local/bin/
```

**macOS (Intel)**:
```bash
curl -L https://github.com/sannticloud/sannti-cli/releases/download/v0.1.0/sannti-v0.1.0-darwin-amd64 -o sannti
chmod +x sannti
sudo mv sannti /usr/local/bin/
```

> **Coming soon** - Binaries will be available in GitHub Releases

### Configuration

Configure your Sannti Cloud credentials:
```bash
sannti configure
```

You'll be prompted for:
- **Sannti Access Key**: Your API access key
- **Sannti Secret Key**: Your API secret key
- **Default Region**: e.g., `br-southeast-1`

Configuration is saved to `~/.sannti/config.yaml`.

### Alternative: Environment Variables
```bash
export SANNTI_ACCESS_KEY="your-access-key"
export SANNTI_SECRET_KEY="your-secret-key"
export SANNTI_REGION="br-southeast-1"
```

## 📚 Usage

### Core Commands
```bash
# Show version
sannti version

# View help
sannti --help
sannti compute --help
```

### Regions
```bash
# List all available regions
sannti region list

# Output as JSON
sannti region list --output json
```

### Compute Instances
```bash
# List all compute instances
sannti compute list

# Get detailed instance information
sannti compute get <instance-uuid>

# List available images (OS templates)
sannti compute images

# List available sizes (compute offerings)
sannti compute sizes

# Create a new instance
sannti compute create \
  --name web-server-01 \
  --region br-southeast-1 \
  --image <template-uuid> \
  --size <offering-uuid> \
  --network <network-uuid>

# Start a stopped instance
sannti compute start <instance-uuid>

# Stop a running instance
sannti compute stop <instance-uuid>

# Delete an instance
sannti compute delete <instance-uuid>
```

### Networking
```bash
# List networks
sannti network list

# List IP addresses
sannti ip list

# List firewall rules
sannti firewall list
```

### Kubernetes
```bash
# List available Kubernetes versions
sannti k8s versions
```

> **Note**: Kubernetes cluster management (create/list/delete) planned for future release.

## 🎨 Output Formats

The CLI supports multiple output formats:
```bash
# Table (default) - human-readable
sannti compute list

# JSON - for scripting/automation
sannti compute list --output json

# YAML - for configuration
sannti compute list --output yaml
```

## 🌍 Regions

Sannti Cloud operates in Brazilian regions:

- `br-southeast-1` - São Paulo region

View all available regions:
```bash
sannti region list
```

## 🔧 Advanced Usage

### Using Different Regions

Override the default region for any command:
```bash
sannti compute list --region br-southeast-1
```

### Scripting Examples

**Export all instances to JSON:**
```bash
sannti compute list --output json > instances.json
```

**Check if region is available:**
```bash
sannti region list --output json | jq '.[] | select(.name=="br-southeast-1")'
```

**List only running instances:**
```bash
sannti compute list --output json | jq '.[] | select(.state=="RUNNING")'
```

**Get instance details with full information:**
```bash
sannti compute get <uuid> --output json | jq '{name, vcpu: .cpuCore, memory, network: .networkName, ip: .instancePrivateIp}'
```

## 🏗️ Development

### Building from Source

Requirements:
- Go 1.21 or higher
```bash
# Clone the repository
git clone https://github.com/sannticloud/sannti-cli.git
cd sannti-cli

# Download dependencies
go mod download

# Build
make build

# Install locally
make install

# Build for all platforms
make release
```

### Project Structure
```
sannti-cli/
├── cmd/               # CLI commands
│   ├── root.go       # Base command + global flags
│   ├── compute.go    # Instance management
│   ├── network.go    # Network operations
│   └── ...
├── internal/
│   ├── client/       # API client
│   ├── config/       # Configuration management
│   ├── output/       # Output formatters
│   └── models/       # Data models
├── scripts/
│   └── install.sh    # Installation script
├── main.go
└── Makefile
```

## 📖 Documentation

- [Quickstart Guide](QUICKSTART.md)
- [Architecture Documentation](ARCHITECTURE.md)
- [Sannti Cloud Documentation](https://docs.sannti.cloud)
- [API Reference](https://console.sannti.cloud/apidocs/swagger-ui/index.html)

## 🐛 Issues & Feedback

Found a bug or have a feature request?

- [Report an issue](https://github.com/sannticloud/sannti-cli/issues)
- Email: support@sannti.com

## 📝 License

Apache License 2.0 - see [LICENSE](LICENSE) file for details.

## 🎯 Roadmap

**v0.2.0** (Next Release):
- [ ] Pre-built binaries (Linux/macOS/Windows)
- [ ] NPM package (`@sannticloud/sannti-cli`)
- [ ] Homebrew formula (`brew install sannti`)
- [ ] Install script at get.sannti.cloud
- [ ] Kubernetes cluster management
- [ ] Volume management
- [ ] SSH key management

**v0.3.0:**
- [ ] Snapshot operations
- [ ] VPC configuration
- [ ] Interactive mode
- [ ] Shell completion (bash/zsh)
- [ ] Progress indicators

**v1.0.0:**
- [ ] Full API coverage
- [ ] Comprehensive test suite
- [ ] Auto-update mechanism
- [ ] Cost estimation

## 🤝 Contributing

Contributions are welcome! This is a v0.1.0 experimental release.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## ⚡ About Sannti Cloud

Sannti Cloud is a Brazilian cloud infrastructure provider focused on digital sovereignty, with data centers in Brazil and full LGPD compliance.

**Key Features:**
- 🇧🇷 100% Brazilian infrastructure
- 🔒 LGPD compliant by design
- 🚀 Low latency for Brazilian users
- 💰 Competitive pricing
- 🎯 Developer-friendly APIs

---

**Made with ❤️ by the Sannti Cloud team**
