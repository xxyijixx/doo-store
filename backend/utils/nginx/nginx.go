package nginx

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"regexp"
	"text/template"

	"doo-store/backend/config"
	"doo-store/backend/constant"
	"doo-store/backend/utils/docker"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/pkg/stdcopy"
	log "github.com/sirupsen/logrus"
)

// NginxManager handles all Nginx-related operations including configuration management,
// location blocks handling, and container operations
type NginxManager struct {
	dockerClient *docker.Client // Docker client for container operations
	containerID  string         // ID of the Nginx container
	nginxConfig  *NginxConfig   // Nginx configuration settings
}

// NginxConfig contains Nginx configuration settings
type NginxConfig struct {
	DefaultTemplate string // Default template for location blocks
	ConfigDir       string // Directory containing Nginx configuration files
	ContainerName   string // Name of the Nginx container
}

// NewNginxManager creates a new NginxManager instance with initialized Docker client
// and container information
func NewNginxManager() (*NginxManager, error) {
	log.Info("Initializing Nginx Manager")
	dockerClient, err := docker.NewClient()
	if err != nil {
		log.Errorf("Failed to create docker client: %v", err)
		return nil, fmt.Errorf("failed to create docker client: %w", err)
	}

	container, err := getContainer(&dockerClient, config.EnvConfig.GetNginxContainerName())
	if err != nil {
		log.Errorf("Failed to get nginx container: %v", err)
		return nil, fmt.Errorf("failed to get nginx container: %w", err)
	}
	log.Infof("Successfully initialized Nginx Manager with container ID: %s", container.ID)

	return &NginxManager{
		dockerClient: &dockerClient,
		containerID:  container.ID,
		nginxConfig: &NginxConfig{
			ConfigDir:     constant.NginxDir,
			ContainerName: config.EnvConfig.GetNginxContainerName(),
		},
	}, nil
}

// AddLocation adds a new location block to Nginx configuration
// It handles the entire process including file generation, container updates,
// and configuration testing
func (nm *NginxManager) AddLocation(locationConfig *LocationConfig) error {
	log.Infof("Adding new location block for: %s", locationConfig.Name)
	locationPath := nm.getLocationPath(locationConfig.Name)

	// Generate configuration content
	content, err := nm.generateLocationContent(locationConfig)
	if err != nil {
		log.Errorf("Failed to generate location content: %v", err)
		return fmt.Errorf("failed to generate location content: %w", err)
	}

	// Write configuration to local file
	if err := os.WriteFile(locationPath, []byte(content), 0644); err != nil {
		log.Errorf("Failed to write configuration file: %v", err)
		return fmt.Errorf("failed to write configuration file: %w", err)
	}
	log.Debugf("Successfully wrote configuration to: %s", locationPath)

	// Handle default configuration backup
	if err := nm.handleDefaultConfig(locationConfig.Name); err != nil {
		log.Errorf("Failed to handle default config: %v", err)
		return fmt.Errorf("failed to handle default config: %w", err)
	}

	// Copy configuration to container
	containerPath := fmt.Sprintf("/etc/nginx/conf.d/apps/%s.conf", locationConfig.Name)
	if err := nm.dockerClient.CopyFileToContainer(nm.containerID, locationPath, containerPath); err != nil {
		log.Errorf("Failed to copy config to container: %v", err)
		nm.rollbackChanges(locationConfig.Name)
		return fmt.Errorf("failed to copy config to container: %w", err)
	}
	log.Debugf("Successfully copied configuration to container path: %s", containerPath)

	// Test and reload configuration
	if err := nm.testAndReload(); err != nil {
		log.Errorf("Failed to test and reload nginx: %v", err)
		nm.rollbackChanges(locationConfig.Name)
		return fmt.Errorf("failed to test and reload nginx: %w", err)
	}

	log.Infof("Successfully added location block for: %s", locationConfig.Name)
	return nil
}

// RemoveLocation removes a location block from Nginx configuration
// It handles cleanup of both container and local files
func (nm *NginxManager) RemoveLocation(locationName string) error {
	log.Infof("Removing location block for: %s", locationName)

	// Remove configuration file from container
	containerPath := fmt.Sprintf("/etc/nginx/conf.d/apps/%s.conf", locationName)
	if err := nm.dockerClient.RemoveFileFormContainer(nm.containerID, containerPath); err != nil {
		log.Errorf("Failed to remove config from container: %v", err)
		return fmt.Errorf("failed to remove config from container: %w", err)
	}

	// Restore default configuration if exists
	if err := nm.restoreDefaultConfig(locationName); err != nil {
		log.Errorf("Failed to restore default config: %v", err)
		return fmt.Errorf("failed to restore default config: %w", err)
	}

	// Remove local configuration file
	locationPath := nm.getLocationPath(locationName)
	if err := os.Remove(locationPath); err != nil {
		log.Errorf("Failed to remove local config file: %v", err)
		return fmt.Errorf("failed to remove local config file: %w", err)
	}

	log.Infof("Successfully removed location block for: %s", locationName)
	return nm.testAndReload()
}

// handleDefaultConfig handles the backup of default configuration
// It creates a backup of existing default configuration if it exists
func (nm *NginxManager) handleDefaultConfig(locationName string) error {
	log.Debugf("Handling default config for: %s", locationName)
	defaultConfPath := fmt.Sprintf("/etc/nginx/conf.d/apps/%s-default.conf", locationName)
	exists, err := nm.dockerClient.FileExistsInContainer(nm.containerID, defaultConfPath)
	if err != nil {
		log.Errorf("Failed to check default config existence: %v", err)
		return fmt.Errorf("failed to check default config existence: %w", err)
	}

	if exists {
		backupPath := defaultConfPath + ".bak"
		if err := nm.dockerClient.MoveFileWithCheck(nm.containerID, defaultConfPath, backupPath); err != nil {
			log.Errorf("Failed to backup default config: %v", err)
			return fmt.Errorf("failed to backup default config: %w", err)
		}
		log.Debugf("Successfully backed up default config to: %s", backupPath)
	}

	return nil
}

// restoreDefaultConfig restores the default configuration if it exists
// It moves the backup file back to its original location
func (nm *NginxManager) restoreDefaultConfig(locationName string) error {
	log.Debugf("Restoring default config for: %s", locationName)
	backupPath := fmt.Sprintf("/etc/nginx/conf.d/apps/%s-default.conf.bak", locationName)
	defaultPath := fmt.Sprintf("/etc/nginx/conf.d/apps/%s-default.conf", locationName)

	exists, err := nm.dockerClient.FileExistsInContainer(nm.containerID, backupPath)
	if err != nil {
		log.Errorf("Failed to check backup config existence: %v", err)
		return fmt.Errorf("failed to check backup config existence: %w", err)
	}

	if exists {
		if err := nm.dockerClient.MoveFileInContainer(nm.containerID, backupPath, defaultPath); err != nil {
			log.Errorf("Failed to restore default config: %v", err)
			return fmt.Errorf("failed to restore default config: %w", err)
		}
		log.Debugf("Successfully restored default config from backup")
	}

	return nil
}

// rollbackChanges rolls back any changes made during the configuration process
// It attempts to restore the system to its previous state in case of failure
func (nm *NginxManager) rollbackChanges(locationName string) error {
	log.Infof("Rolling back changes for: %s", locationName)
	containerPath := fmt.Sprintf("/etc/nginx/conf.d/apps/%s.conf", locationName)

	// Remove the new configuration file if it exists
	if err := nm.dockerClient.RemoveFileFormContainer(nm.containerID, containerPath); err != nil {
		log.Warnf("Failed to remove new config during rollback: %v", err)
	}

	// Restore the default configuration if it was backed up
	if err := nm.restoreDefaultConfig(locationName); err != nil {
		log.Warnf("Failed to restore default config during rollback: %v", err)
	}

	// Remove the local configuration file
	locationPath := nm.getLocationPath(locationName)
	if err := os.Remove(locationPath); err != nil {
		log.Warnf("Failed to remove local config file during rollback: %v", err)
	}

	log.Info("Rollback completed")
	return nil
}

// ExtractLocations extracts all location blocks from Nginx configuration
// It uses regex to find and return all location paths
func (nm *NginxManager) ExtractLocations(nginxConfig string) []string {
	log.Debug("Extracting locations from nginx config")
	re := regexp.MustCompile(`location\s+(/[^/]+(?:/[^/]+)*/*)\s+{`)
	matches := re.FindAllStringSubmatch(nginxConfig, -1)

	locations := make([]string, 0, len(matches))
	for _, match := range matches {
		if len(match) > 1 {
			locations = append(locations, match[1])
		}
	}
	log.Debugf("Found %d locations", len(locations))
	return locations
}


// generateLocationContent generates the content for a location block
// It either uses a custom template or generates a default one
func (nm *NginxManager) generateLocationContent(locationConfig *LocationConfig) (string, error) {
	log.Debugf("Generating location content for: %s", locationConfig.Name)
	if locationConfig.Template == "" {
		return nm.generateDefaultTemplate(locationConfig), nil
	}

	t, err := template.New("nginx").Parse(locationConfig.Template)
	if err != nil {
		log.Errorf("Failed to parse template: %v", err)
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	data := map[string]interface{}{
		"Key":           locationConfig.Name,
		"ContainerName": locationConfig.ProxyServerName,
		"Port":          locationConfig.Port,
	}

	if err := t.Execute(&buf, data); err != nil {
		log.Errorf("Failed to execute template: %v", err)
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

// generateDefaultTemplate generates a default Nginx location block template
func (nm *NginxManager) generateDefaultTemplate(locationConfig *LocationConfig) string {
	log.Debugf("Generating default template for: %s", locationConfig.Name)
	proxyPass := fmt.Sprintf("http://%s/", locationConfig.ProxyServerName)
	if locationConfig.Port != 0 {
		proxyPass = fmt.Sprintf("http://%s:%d/", locationConfig.ProxyServerName, locationConfig.Port)
	}

	return fmt.Sprintf(`location /plugin/%s/ {
	proxy_http_version 1.1;
	proxy_set_header X-Real-IP $remote_addr;
	proxy_set_header X-Real-PORT $remote_port;
	proxy_set_header X-Forwarded-Host $the_host;
	proxy_set_header X-Forwarded-Proto $the_scheme;
	proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
	proxy_set_header Host $http_host;
	proxy_set_header Scheme $scheme;
	proxy_set_header Server-Protocol $server_protocol;
	proxy_set_header Server-Name $server_name;
	proxy_set_header Server-Addr $server_addr;
	proxy_set_header Server-Port $server_port;
	proxy_set_header Upgrade $http_upgrade;
	proxy_set_header Connection $connection_upgrade;
	proxy_read_timeout 3600s;
	proxy_send_timeout 3600s;
	proxy_connect_timeout 3600s;
	proxy_pass %s;
}`, locationConfig.Name, proxyPass)
}

// ExtractLocationsByKey extracts locations from a specific configuration file
func (nm *NginxManager) ExtractLocationsByKey(key string) ([]string, error) {
	log.Debugf("Extracting locations for key: %s", key)
	locationPath := nm.getLocationPath(key)
	content, err := os.ReadFile(locationPath)
	if err != nil {
		log.Errorf("Failed to read Nginx config file: %v", err)
		return []string{}, err
	}
	locations := nm.ExtractLocations(string(content))
	return locations, nil
}

// testAndReload tests the Nginx configuration and reloads if valid
func (nm *NginxManager) testAndReload() error {
	log.Debug("Testing and reloading Nginx configuration")
	if err := nm.testConfig(); err != nil {
		log.Errorf("Nginx configuration test failed: %v", err)
		return fmt.Errorf("nginx configuration test failed: %w", err)
	}
	return nm.reload()
}

// testConfig tests the Nginx configuration
func (nm *NginxManager) testConfig() error {
	log.Debug("Testing Nginx configuration")
	return nm.executeCommand([]string{"nginx", "-t"})
}

// reload reloads the Nginx configuration
func (nm *NginxManager) reload() error {
	log.Debug("Reloading Nginx configuration")
	return nm.executeCommand([]string{"nginx", "-s", "reload"})
}

// executeCommand executes a command in the Nginx container
func (nm *NginxManager) executeCommand(cmd []string) error {
	log.Debugf("Executing command in container: %v", cmd)
	execConfig := container.ExecOptions{
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          cmd,
	}

	execIDResp, err := nm.dockerClient.GetClient().ContainerExecCreate(context.Background(), nm.containerID, execConfig)
	if err != nil {
		log.Errorf("Failed to create exec: %v", err)
		return fmt.Errorf("failed to create exec: %w", err)
	}

	execAttachResp, err := nm.dockerClient.GetClient().ContainerExecAttach(context.Background(), execIDResp.ID, container.ExecStartOptions{})
	if err != nil {
		log.Errorf("Failed to attach exec: %v", err)
		return fmt.Errorf("failed to attach exec: %w", err)
	}
	defer execAttachResp.Close()

	outputDone := make(chan error)
	go func() {
		_, err := stdcopy.StdCopy(os.Stdout, os.Stderr, execAttachResp.Reader)
		outputDone <- err
	}()

	if err := <-outputDone; err != nil && err != io.EOF {
		log.Errorf("Command execution failed: %v", err)
		return fmt.Errorf("command execution failed: %w", err)
	}

	return nil
}

// getLocationPath returns the full path for a location configuration file
func (nm *NginxManager) getLocationPath(key string) string {
	return fmt.Sprintf("%s/%s.conf", nm.nginxConfig.ConfigDir, key)
}

// getContainer retrieves the Nginx container by name
func getContainer(client *docker.Client, name string) (types.Container, error) {
	log.Debugf("Getting container with name: %s", name)
	containers, err := client.ListContainersByName([]string{name})
	if err != nil {
		log.Errorf("Failed to list containers: %v", err)
		return types.Container{}, err
	}
	if len(containers) == 0 {
		log.Errorf("Nginx container not found: %s", name)
		return types.Container{}, fmt.Errorf("nginx container not found: %s", name)
	}
	return containers[0], nil
}
