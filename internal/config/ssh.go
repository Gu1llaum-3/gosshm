package config

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// SSHHost represents an SSH host configuration
type SSHHost struct {
	Name      string
	Hostname  string
	User      string
	Port      string
	Identity  string
	ProxyJump string
	Tags      []string
}

// configMutex protects SSH config file operations from race conditions
var configMutex sync.Mutex

// backupConfig creates a backup of the SSH config file
func backupConfig(configPath string) error {
	backupPath := configPath + ".backup"
	src, err := os.Open(configPath)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(backupPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	return err
}

// ParseSSHConfig parses the SSH config file and returns the list of hosts
func ParseSSHConfig() ([]SSHHost, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(homeDir, ".ssh", "config")
	return ParseSSHConfigFile(configPath)
}

// ParseSSHConfigFile parses a specific SSH config file and returns the list of hosts
func ParseSSHConfigFile(configPath string) ([]SSHHost, error) {
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var hosts []SSHHost
	var currentHost *SSHHost
	var pendingTags []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Ignore empty lines
		if line == "" {
			continue
		}

		// Check for tags comment
		if strings.HasPrefix(line, "# Tags:") {
			tagsStr := strings.TrimPrefix(line, "# Tags:")
			tagsStr = strings.TrimSpace(tagsStr)
			if tagsStr != "" {
				// Split tags by comma and trim whitespace
				for _, tag := range strings.Split(tagsStr, ",") {
					tag = strings.TrimSpace(tag)
					if tag != "" {
						pendingTags = append(pendingTags, tag)
					}
				}
			}
			continue
		}

		// Ignore other comments
		if strings.HasPrefix(line, "#") {
			continue
		}

		// Split line into words
		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}

		key := strings.ToLower(parts[0])
		value := strings.Join(parts[1:], " ")

		switch key {
		case "host":
			// New host, save previous one if it exists
			if currentHost != nil {
				hosts = append(hosts, *currentHost)
			}
			// Create new host
			currentHost = &SSHHost{
				Name: value,
				Port: "22",        // Default port
				Tags: pendingTags, // Assign pending tags to this host
			}
			// Clear pending tags for next host
			pendingTags = nil
		case "hostname":
			if currentHost != nil {
				currentHost.Hostname = value
			}
		case "user":
			if currentHost != nil {
				currentHost.User = value
			}
		case "port":
			if currentHost != nil {
				currentHost.Port = value
			}
		case "identityfile":
			if currentHost != nil {
				currentHost.Identity = value
			}
		case "proxyjump":
			if currentHost != nil {
				currentHost.ProxyJump = value
			}
		}
	}

	// Add the last host if it exists
	if currentHost != nil {
		hosts = append(hosts, *currentHost)
	}

	return hosts, scanner.Err()
}

// AddSSHHost adds a new SSH host to the config file
func AddSSHHost(host SSHHost) error {
	configMutex.Lock()
	defer configMutex.Unlock()

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configPath := filepath.Join(homeDir, ".ssh", "config")

	// Create backup before modification if file exists
	if _, err := os.Stat(configPath); err == nil {
		if err := backupConfig(configPath); err != nil {
			return fmt.Errorf("failed to create backup: %w", err)
		}
	}

	// Check if host already exists
	exists, err := HostExists(host.Name)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("host '%s' already exists", host.Name)
	}

	// Open file in append mode
	file, err := os.OpenFile(configPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the configuration
	_, err = file.WriteString("\n")
	if err != nil {
		return err
	}

	// Write tags if present
	if len(host.Tags) > 0 {
		_, err = file.WriteString("# Tags: " + strings.Join(host.Tags, ", ") + "\n")
		if err != nil {
			return err
		}
	}

	// Write host configuration
	_, err = file.WriteString(fmt.Sprintf("Host %s\n", host.Name))
	if err != nil {
		return err
	}

	_, err = file.WriteString(fmt.Sprintf("    HostName %s\n", host.Hostname))
	if err != nil {
		return err
	}

	if host.User != "" {
		_, err = file.WriteString(fmt.Sprintf("    User %s\n", host.User))
		if err != nil {
			return err
		}
	}

	if host.Port != "" && host.Port != "22" {
		_, err = file.WriteString(fmt.Sprintf("    Port %s\n", host.Port))
		if err != nil {
			return err
		}
	}

	if host.Identity != "" {
		_, err = file.WriteString(fmt.Sprintf("    IdentityFile %s\n", host.Identity))
		if err != nil {
			return err
		}
	}

	if host.ProxyJump != "" {
		_, err = file.WriteString(fmt.Sprintf("    ProxyJump %s\n", host.ProxyJump))
		if err != nil {
			return err
		}
	}

	return nil
}

// HostExists checks if a host already exists in the config
func HostExists(hostName string) (bool, error) {
	hosts, err := ParseSSHConfig()
	if err != nil {
		return false, err
	}

	for _, host := range hosts {
		if host.Name == hostName {
			return true, nil
		}
	}
	return false, nil
}

// GetSSHHost retrieves a specific host configuration by name
func GetSSHHost(hostName string) (*SSHHost, error) {
	hosts, err := ParseSSHConfig()
	if err != nil {
		return nil, err
	}

	for _, host := range hosts {
		if host.Name == hostName {
			return &host, nil
		}
	}
	return nil, fmt.Errorf("host '%s' not found", hostName)
}

// UpdateSSHHost updates an existing SSH host configuration
func UpdateSSHHost(oldName string, newHost SSHHost) error {
	configMutex.Lock()
	defer configMutex.Unlock()

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configPath := filepath.Join(homeDir, ".ssh", "config")

	// Create backup before modification
	if err := backupConfig(configPath); err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}

	// Read the current config
	content, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	var newLines []string
	i := 0
	hostFound := false

	for i < len(lines) {
		line := strings.TrimSpace(lines[i])

		// Check for tags comment followed by Host
		if strings.HasPrefix(line, "# Tags:") && i+1 < len(lines) {
			nextLine := strings.TrimSpace(lines[i+1])
			if nextLine == "Host "+oldName {
				// Found the host to update, skip the old configuration
				hostFound = true

				// Skip until we find the end of this host block (empty line or next Host)
				i += 2 // Skip tags and Host line
				for i < len(lines) && strings.TrimSpace(lines[i]) != "" && !strings.HasPrefix(strings.TrimSpace(lines[i]), "Host ") {
					i++
				}

				// Insert new configuration at this position
				newLines = append(newLines, "")
				if len(newHost.Tags) > 0 {
					newLines = append(newLines, "# Tags: "+strings.Join(newHost.Tags, ", "))
				}
				newLines = append(newLines, "Host "+newHost.Name)
				newLines = append(newLines, "    HostName "+newHost.Hostname)
				if newHost.User != "" {
					newLines = append(newLines, "    User "+newHost.User)
				}
				if newHost.Port != "" && newHost.Port != "22" {
					newLines = append(newLines, "    Port "+newHost.Port)
				}
				if newHost.Identity != "" {
					newLines = append(newLines, "    IdentityFile "+newHost.Identity)
				}

				continue
			}
		}

		// Check for Host line without tags
		if strings.HasPrefix(line, "Host ") && strings.Fields(line)[1] == oldName {
			hostFound = true

			// Skip until we find the end of this host block
			i++ // Skip Host line
			for i < len(lines) && strings.TrimSpace(lines[i]) != "" && !strings.HasPrefix(strings.TrimSpace(lines[i]), "Host ") {
				i++
			}

			// Insert new configuration
			newLines = append(newLines, "")
			if len(newHost.Tags) > 0 {
				newLines = append(newLines, "# Tags: "+strings.Join(newHost.Tags, ", "))
			}
			newLines = append(newLines, "Host "+newHost.Name)
			newLines = append(newLines, "    HostName "+newHost.Hostname)
			if newHost.User != "" {
				newLines = append(newLines, "    User "+newHost.User)
			}
			if newHost.Port != "" && newHost.Port != "22" {
				newLines = append(newLines, "    Port "+newHost.Port)
			}
			if newHost.Identity != "" {
				newLines = append(newLines, "    IdentityFile "+newHost.Identity)
			}

			continue
		}

		// Keep other lines as-is
		newLines = append(newLines, lines[i])
		i++
	}

	if !hostFound {
		return fmt.Errorf("host '%s' not found", oldName)
	}

	// Write back to file
	newContent := strings.Join(newLines, "\n")
	return os.WriteFile(configPath, []byte(newContent), 0600)
}

// DeleteSSHHost removes an SSH host configuration from the config file
func DeleteSSHHost(hostName string) error {
	configMutex.Lock()
	defer configMutex.Unlock()

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configPath := filepath.Join(homeDir, ".ssh", "config")

	// Create backup before modification
	if err := backupConfig(configPath); err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}

	// Read the current config
	content, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	var newLines []string
	i := 0
	hostFound := false

	for i < len(lines) {
		line := strings.TrimSpace(lines[i])

		// Check for tags comment followed by Host
		if strings.HasPrefix(line, "# Tags:") && i+1 < len(lines) {
			nextLine := strings.TrimSpace(lines[i+1])
			if nextLine == "Host "+hostName {
				// Found the host to delete, skip the configuration
				hostFound = true

				// Skip tags comment and Host line
				i += 2

				// Skip until we find the end of this host block (empty line or next Host)
				for i < len(lines) && strings.TrimSpace(lines[i]) != "" && !strings.HasPrefix(strings.TrimSpace(lines[i]), "Host ") {
					i++
				}

				// Skip the empty line after the host block if it exists
				if i < len(lines) && strings.TrimSpace(lines[i]) == "" {
					i++
				}

				continue
			}
		}

		// Check for Host line without tags
		if strings.HasPrefix(line, "Host ") && strings.Fields(line)[1] == hostName {
			hostFound = true

			// Skip Host line
			i++

			// Skip until we find the end of this host block
			for i < len(lines) && strings.TrimSpace(lines[i]) != "" && !strings.HasPrefix(strings.TrimSpace(lines[i]), "Host ") {
				i++
			}

			// Skip the empty line after the host block if it exists
			if i < len(lines) && strings.TrimSpace(lines[i]) == "" {
				i++
			}

			continue
		}

		// Keep other lines as-is
		newLines = append(newLines, lines[i])
		i++
	}

	if !hostFound {
		return fmt.Errorf("host '%s' not found", hostName)
	}

	// Write back to file
	newContent := strings.Join(newLines, "\n")
	return os.WriteFile(configPath, []byte(newContent), 0600)
}
