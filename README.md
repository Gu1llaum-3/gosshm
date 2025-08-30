# 🚀 SSHM - SSH Connection Manager

[![Go](https://img.shields.io/badge/Go-1.23+-00ADD8?style=for-the-badge&logo=go)](https://golang.org/)
[![Release](https://img.shields.io/github/v/release/Gu1llaum-3/gosshm?style=for-the-badge)](https://github.com/Gu1llaum-3/gosshm/releases)
[![License](https://img.shields.io/github/license/Gu1llaum-3/gosshm?style=for-the-badge)](LICENSE)
[![Platform](https://img.shields.io/badge/platform-Linux%20%7C%20macOS-lightgrey?style=for-the-badge)](https://github.com/Gu1llaum-3/gosshm/releases)

> **A modern, interactive SSH connection manager for your terminal** 🔥

SSHM is a beautiful command-line tool that transforms how you manage and connect to your SSH hosts. Built with Go and featuring an intuitive TUI interface, it makes SSH connection management effortless and enjoyable.

## ✨ Features

### 🎯 **Core Features**
- **🎨 Beautiful TUI Interface** - Navigate your SSH hosts with an elegant, interactive terminal UI
- **⚡ Quick Connect** - Connect to any host instantly
- **📝 Easy Management** - Add, edit, and manage SSH configurations seamlessly
- **🏷️ Tag Support** - Organize your hosts with custom tags for better categorization
- **🔍 Smart Search** - Find hosts quickly with built-in filtering and search
- **🔒 Secure** - Works directly with your existing `~/.ssh/config` file

### 🛠️ **Management Operations**
- **Add new SSH hosts** with interactive forms
- **Edit existing configurations** in-place
- **Delete hosts** with confirmation prompts
- **Backup configurations** automatically before changes
- **Validate settings** to prevent configuration errors
- **ProxyJump support** for secure connection tunneling through bastion hosts

### 🎮 **User Experience**
- **Zero configuration** - Works out of the box with your existing SSH setup
- **Keyboard shortcuts** for power users
- **Cross-platform** - Supports Linux and macOS (Intel & Apple Silicon)
- **Lightweight** - Single binary with no dependencies

## 🚀 Quick Start

### Installation

**One-line install (Recommended):**
```bash
curl -sSL https://raw.githubusercontent.com/Gu1llaum-3/gosshm/main/install/unix.sh | bash
```

**Alternative methods:**
```bash
# Download specific release
wget https://github.com/Gu1llaum-3/gosshm/releases/latest/download/sshm-linux-amd64.tar.gz

# Extract and install
tar -xzf sshm-linux-amd64.tar.gz
sudo mv sshm-linux-amd64 /usr/local/bin/sshm
```

## 📖 Usage

### Interactive Mode

Launch SSHM without arguments to enter the beautiful TUI interface:

```bash
sshm
```

**Navigation:**
- `↑/↓` or `j/k` - Navigate hosts
- `Enter` - Connect to selected host
- `a` - Add new host
- `e` - Edit selected host
- `d` - Delete selected host
- `q` - Quit
- `/` - Search/filter hosts

The interactive forms will guide you through configuration:
- **Hostname/IP** - Server address
- **Username** - SSH user
- **Port** - SSH port (default: 22)
- **Identity File** - Private key path
- **ProxyJump** - Jump server for connection tunneling
- **Tags** - Comma-separated tags for organization

### CLI Usage

SSHM provides both command-line operations and an interactive TUI interface:

```bash
# Launch interactive TUI mode for browsing and connecting to hosts
sshm

# Add a new host using interactive form
sshm add

# Add a new host with pre-filled hostname
sshm add hostname

# Edit an existing host configuration
sshm edit my-server

# Show version information
sshm --version

# Show help and available commands
sshm --help
```

## 🏗️ Configuration

SSHM works directly with your standard SSH configuration file (`~/.ssh/config`). It adds special comment tags for enhanced functionality while maintaining full compatibility with standard SSH tools.

**Example configuration:**
```ssh
# Tags: production, web, frontend
Host web-prod-01
    HostName 192.168.1.10
    User deploy
    Port 22
    IdentityFile ~/.ssh/production_key

# Tags: development, database
Host db-dev
    HostName dev-db.company.com
    User admin
    Port 2222
    IdentityFile ~/.ssh/dev_key

# Tags: production, backend
Host backend-prod
    HostName 10.0.1.50
    User app
    Port 22
    ProxyJump bastion.company.com
    IdentityFile ~/.ssh/production_key
```

### Supported SSH Options

- `HostName` - Server hostname or IP address
- `User` - Username for SSH connection
- `Port` - SSH port number
- `IdentityFile` - Path to private key file
- `ProxyJump` - Jump server for connection tunneling (e.g., `user@jumphost:port`)
- `Tags` - Custom tags (SSHM extension)

## 🛠️ Development

### Prerequisites

- Go 1.23+ 
- Git

### Build from Source

```bash
# Clone the repository
git clone https://github.com/Gu1llaum-3/gosshm.git
cd gosshm

# Build the binary
go build -o sshm .

# Run
./sshm
```

### Project Structure

```
sshm/
├── cmd/                 # CLI commands (Cobra)
│   ├── root.go         # Root command and interactive mode
│   ├── add.go          # Add host command
│   └── edit.go         # Edit host command
├── internal/
│   ├── config/         # SSH configuration management
│   │   └── ssh.go      # Config parsing and manipulation
│   ├── ui/             # Terminal UI components
│   │   ├── tui.go      # Main TUI interface
│   │   ├── add_form.go # Add host form
│   │   └── edit_form.go# Edit host form
│   └── validation/     # Input validation
│       └── ssh.go      # SSH config validation
├── install/            # Installation scripts
│   ├── unix.sh         # Unix/Linux/macOS installer
│   └── README.md       # Installation guide
└── .github/workflows/  # CI/CD pipelines
    └── build.yml       # Multi-platform builds
```

### Dependencies

- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Bubbles](https://github.com/charmbracelet/bubbles) - TUI components
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Styling

## 📦 Releases

Automated releases are built for multiple platforms:

| Platform | Architecture | Download |
|----------|-------------|----------|
| Linux | AMD64 | [sshm-linux-amd64.tar.gz](https://github.com/Gu1llaum-3/gosshm/releases/latest/download/sshm-linux-amd64.tar.gz) |
| Linux | ARM64 | [sshm-linux-arm64.tar.gz](https://github.com/Gu1llaum-3/gosshm/releases/latest/download/sshm-linux-arm64.tar.gz) |
| macOS | Intel | [sshm-darwin-amd64.tar.gz](https://github.com/Gu1llaum-3/gosshm/releases/latest/download/sshm-darwin-amd64.tar.gz) |
| macOS | Apple Silicon | [sshm-darwin-arm64.tar.gz](https://github.com/Gu1llaum-3/gosshm/releases/latest/download/sshm-darwin-arm64.tar.gz) |

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

### Development Workflow

1. **Fork** the repository
2. **Create** a feature branch (`git checkout -b feature/amazing-feature`)
3. **Commit** your changes (`git commit -m 'Add amazing feature'`)
4. **Push** to the branch (`git push origin feature/amazing-feature`)
5. **Open** a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Charm](https://charm.sh/) for the amazing TUI libraries
- [Cobra](https://cobra.dev/) for the excellent CLI framework
- The Go community for building such fantastic tools

---

<div align="center">

**Made with ❤️ by [Guillaume](https://github.com/Gu1llaum-3)**

⭐ **Star this repo if you found it useful!** ⭐

</div>
