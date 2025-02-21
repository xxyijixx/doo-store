package compose

import (
	"bufio"
	"encoding/json"
	"fmt"
	"strings"
)

// DockerContainer represents a container from docker compose ps --format json output
type DockerContainer struct {
	Command      string `json:"Command"`
	CreatedAt    string `json:"CreatedAt"`
	ExitCode     int    `json:"ExitCode"`
	Health       string `json:"Health"`
	ID           string `json:"ID"`
	Image        string `json:"Image"`
	Labels       string `json:"Labels"`
	LocalVolumes string `json:"LocalVolumes"`
	Mounts       string `json:"Mounts"`
	Name         string `json:"Name"`
	Names        string `json:"Names"`
	Networks     string `json:"Networks"`
	Ports        string `json:"Ports"`
	Project      string `json:"Project"`
	Publishers   []struct {
		URL           string `json:"URL"`
		TargetPort    int    `json:"TargetPort"`
		PublishedPort int    `json:"PublishedPort"`
		Protocol      string `json:"Protocol"`
	} `json:"Publishers"`
	RunningFor string `json:"RunningFor"`
	Service    string `json:"Service"`
	Size       string `json:"Size"`
	State      string `json:"State"`
	Status     string `json:"Status"`
}

func ParseDockerComposePsOutput(filePath string) ([]DockerContainer, error) {
	output, err := Operate(filePath, "ps --format json")
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(strings.NewReader(output))
	containers := []DockerContainer{}

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		var container DockerContainer
		err := json.Unmarshal([]byte(line), &container)
		if err != nil {
			fmt.Printf("Error parsing JSON: %v\n", err)
			continue
		}

		containers = append(containers, container)
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading input: %v\n", err)
		return nil, err
	}

	return containers, nil
}
