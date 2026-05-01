# action1-cli

CLI for the Action1 API — manage your endpoints, automations, patches, and more from the terminal.

[Documentation en francais](README.fr.md)

## Features

- **~140 API endpoints** covered: endpoints, automations, reports, software deployment, vulnerabilities, RBAC, and more
- **Secure credential storage**: client secrets stored in OS keychain (macOS Keychain, Windows Credential Manager, Linux Secret Service)
- **Multi-profile**: manage multiple Action1 accounts and organizations
- **Multi-region**: North America, Europe, Australia
- **Flexible output**: table, JSON, CSV, YAML
- **Cross-platform**: macOS, Linux, Windows (amd64 and arm64)
- **Auto-pagination**: fetch all results with `--all`

## Prerequisites

- An [Action1](https://www.action1.com) account with API access
- **API credentials** (Client ID and Client Secret) generated from the Action1 console

## Installation

### Method 1: Download the binary (recommended)

Go to the [Releases](https://github.com/dimer47/action1-cli/releases/latest) page and download the archive for your platform.

Or with a single command:

**macOS (Apple Silicon — M1/M2/M3/M4):**

```bash
curl -sL https://github.com/dimer47/action1-cli/releases/latest/download/action1-cli_darwin_arm64.tar.gz | tar xz
sudo mv action1 /usr/local/bin/
```

**macOS (Intel):**

```bash
curl -sL https://github.com/dimer47/action1-cli/releases/latest/download/action1-cli_darwin_amd64.tar.gz | tar xz
sudo mv action1 /usr/local/bin/
```

**Linux (amd64):**

```bash
curl -sL https://github.com/dimer47/action1-cli/releases/latest/download/action1-cli_linux_amd64.tar.gz | tar xz
sudo mv action1 /usr/local/bin/
```

**Linux (arm64 — Raspberry Pi, etc.):**

```bash
curl -sL https://github.com/dimer47/action1-cli/releases/latest/download/action1-cli_linux_arm64.tar.gz | tar xz
sudo mv action1 /usr/local/bin/
```

**Windows:**

Download `action1-cli_windows_amd64.zip` from the [Releases](https://github.com/dimer47/action1-cli/releases/latest) page, extract and add the folder to your `PATH`.

### Method 2: From source (requires Go 1.21+)

```bash
go install github.com/dimer47/action1-cli@latest
```

The binary will be installed in `$GOPATH/bin/` (usually `~/go/bin/`). Make sure this directory is in your `PATH`.

### Method 3: Build locally

```bash
git clone https://github.com/dimer47/action1-cli.git
cd action1-cli
go build -o action1 .
./action1 version
```

### Verify installation

```bash
action1 version
# action1 version 0.1.0
```

## Quick Start

### 1. Generate API credentials

1. Log in to the [Action1 console](https://app.action1.com)
2. Go to **Configuration > API Credentials**
3. Click **Create API Credentials**
4. Copy the **Client ID** and **Client Secret**
5. Note your region (North America, Europe, or Australia)

### 2. Authenticate

```bash
action1 auth login --region eu
```

Answer the prompts:
```
Region: eu (https://app.eu.action1.com/api/3.0)
Client ID: api-key-xxxxx@action1.com     # Paste your Client ID
Client Secret:                             # Paste your secret (hidden)
Successfully authenticated.
Token expires in 3600 seconds.
```

The client credentials are stored in the **OS keychain** (encrypted). Tokens are stored locally with restricted permissions (0600).

### 3. Set your default organization

```bash
# List organizations to find your org ID
action1 org list

# Set the default
action1 config set org <orgId>
```

### 4. Start using it

```bash
# List your endpoints
action1 endpoint list

# Check endpoint status
action1 endpoint status

# List vulnerabilities
action1 vulnerability list

# Run a report
action1 report data <reportId>
```

## Configuration

### Config file

Located at:
- **macOS**: `~/Library/Application Support/action1/config.yaml`
- **Linux**: `~/.config/action1/config.yaml`
- **Windows**: `%APPDATA%\action1\config.yaml`

### Multi-profile

```bash
# Create a "production" profile
action1 config use-profile production
action1 auth login --region eu

# Create a "staging" profile
action1 config use-profile staging
action1 auth login --region na

# List profiles (* = active)
action1 config list-profiles
#   default
# * production
#   staging

# Switch profile
action1 config use-profile production

# Use a profile for a single command
action1 endpoint list --profile staging
```

### Credential storage

| OS | Backend | What is stored |
|----|---------|----------------|
| macOS | Keychain | Client ID + Secret |
| Windows | Credential Manager | Client ID + Secret |
| Linux | Secret Service (GNOME Keyring / KWallet) | Client ID + Secret |
| Fallback | Config file (0600) | Everything |

Tokens (access + refresh) are always stored in the config file with restricted permissions, as they are too large for some keychain implementations.

Use `--no-keychain` to force file-based storage for everything.

## Usage

### Endpoints

```bash
action1 endpoint list                              # List all endpoints
action1 endpoint list --all                        # Fetch all (auto-paginate)
action1 endpoint list --filter "name eq 'SRV01'"   # OData filter
action1 endpoint get <id>                          # Endpoint details
action1 endpoint status                            # Online/offline counts
action1 endpoint update <id> --name "new-name"     # Rename
action1 endpoint delete <id>                       # Delete (asks confirmation)
action1 endpoint move <id> --to-org <orgId>        # Move to another org
action1 endpoint missing-updates <id>              # Missing patches
action1 endpoint install-url windowsEXE            # Agent install URL
```

### Endpoint Groups

```bash
action1 endpoint-group list                        # List groups
action1 endpoint-group create --name "Servers"     # Create
action1 endpoint-group members <id>                # List members
action1 endpoint-group add <id> --endpoints a,b,c  # Add endpoints
action1 endpoint-group remove <id> --endpoints a,b # Remove endpoints
```

### Automations

```bash
# Schedules
action1 automation schedule list
action1 automation schedule create --data @schedule.json
action1 automation schedule get <id>
action1 automation schedule delete <id>

# Run immediately
action1 automation instance run --data @automation.json
action1 automation instance results <instanceId>
action1 automation instance stop <instanceId>

# Action templates
action1 automation template list
action1 automation template get <templateId>
```

### Reports

```bash
action1 report list                                # List reports
action1 report data <reportId>                     # Get report rows
action1 report data <reportId> --all               # All rows
action1 report export <reportId> --output-file r.csv
action1 report requery <reportId>                  # Re-run
action1 report drilldown <reportId> <rowId>        # Row details
```

### Software Repository

```bash
action1 software list                              # List packages
action1 software get <id>                          # Package details
action1 software clone <id>                        # Clone a package
action1 software version create <pkgId> --data @v.json
action1 software upload <pkgId> <verId> ./installer.exe
```

### Updates (Patches)

```bash
action1 update list                                # All missing updates
action1 update get <packageId>                     # Updates for a package
action1 update endpoints <pkgId> <verId>           # Endpoints missing update
```

### Installed Software Inventory

```bash
action1 installed-software list                    # All installed apps
action1 installed-software get <endpointId>        # Apps on one endpoint
action1 installed-software requery                 # Refresh data
```

### Vulnerabilities

```bash
action1 vulnerability list                         # Vulnerable software
action1 vulnerability get <cveId>                  # CVE details (org)
action1 vulnerability cve <cveId>                  # CVE details (global)
action1 vulnerability endpoints <cveId>            # Affected endpoints
action1 vulnerability remediation list <cveId>     # Past remediations
action1 vulnerability remediation create <cveId> --data @control.json
```

### Scripts

```bash
action1 script list                                # List scripts
action1 script create --name "Cleanup" --type powershell --file cleanup.ps1
action1 script get <id>
action1 script update <id> --file cleanup_v2.ps1
action1 script delete <id>
```

### Users & RBAC

```bash
action1 user me                                    # Current user
action1 user list                                  # All users
action1 user create --email admin@co.com --name "Admin"
action1 role list                                  # All roles
action1 role assign <roleId> <userId>              # Assign user to role
action1 role permissions                           # Permission templates
```

### Organizations

```bash
action1 org list                                   # List orgs
action1 org create --name "Production"             # Create
action1 org update <id> --data '{"name":"Prod"}'   # Update
action1 org delete <id>                            # Delete
```

### Audit Trail

```bash
action1 audit list                                 # Audit events
action1 audit list --from 2026-01-01 --to 2026-01-31
action1 audit get <id>                             # Event details
action1 audit export --output-file audit.json      # Export
```

### Other Commands

```bash
action1 search "query"                             # Quick search
action1 log get                                    # Diagnostic logs
action1 enterprise get                             # Enterprise settings
action1 subscription info                          # License info
action1 subscription usage                         # Usage stats
```

## Global Flags

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--org` | `-o` | Organization ID | from config |
| `--region` | `-r` | Server region: `na`, `eu`, `au` | from config |
| `--output` | `-O` | Output format: `table`, `json`, `csv`, `yaml` | `table` |
| `--profile` | `-p` | Config profile | from config |
| `--quiet` | `-q` | Suppress headers and decorations | `false` |
| `--verbose` | `-v` | Show HTTP requests | `false` |
| `--no-color` | | Disable colored output | `false` |
| `--no-keychain` | | Use file-based credential storage | `false` |
| `--config` | | Config file path | auto-detected |

## Shell Completion

```bash
# Bash
action1 completion bash > /etc/bash_completion.d/action1

# Zsh (add to your .zshrc)
action1 completion zsh > "${fpath[1]}/_action1"

# Fish
action1 completion fish > ~/.config/fish/completions/action1.fish

# PowerShell
action1 completion powershell > action1.ps1
```

## Command Aliases

For faster typing:

| Command | Alias |
|---------|-------|
| `endpoint` | `ep` |
| `endpoint-group` | `epg` |
| `automation` | `auto` |
| `vulnerability` | `vuln` |
| `software` | `sw` |
| `installed-software` | `isw` |
| `data-source` | `ds` |
| `subscription` | `sub` |
| `report-subscription` | `report-sub` |

## Development

```bash
# Clone
git clone https://github.com/dimer47/action1-cli.git
cd action1-cli

# Build
go build -o action1 .

# Run tests
go test ./...

# Lint
go vet ./...
```

### Creating a new release

```bash
git tag v0.1.0
git push origin v0.1.0
# GitHub Actions builds and publishes automatically
```

## License

MIT
