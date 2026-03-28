package service

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"

	"host-lotus-worker/internal/domain"
)

const (
	minecraftImage  = "itzg/minecraft-server"
	containerPrefix = "mc-server-"
)

// DockerService manages Docker containers for game servers
type DockerService struct {
	client     *client.Client
	dataDir    string
	portOffset int // Base port offset for server allocation
}

// NewDockerService creates a new Docker service
func NewDockerService(dataDir string, portOffset int) (*DockerService, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("failed to create docker client: %w", err)
	}

	// Verify Docker is accessible
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = cli.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to docker: %w", err)
	}

	return &DockerService{
		client:     cli,
		dataDir:    dataDir,
		portOffset: portOffset,
	}, nil
}

// Close closes the Docker client connection
func (d *DockerService) Close() error {
	return d.client.Close()
}

// GetClient returns the underlying Docker client for use by other services
func (d *DockerService) GetClient() *client.Client {
	return d.client
}

// StartServer starts a Minecraft server container for the given server configuration
func (d *DockerService) StartServer(ctx context.Context, server *domain.ServerDynamic, hostPort int) error {
	containerName := containerPrefix + server.ID

	// Check if container already exists
	containers, err := d.client.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		return fmt.Errorf("failed to list containers: %w", err)
	}

	for _, c := range containers {
		for _, name := range c.Names {
			if name == "/"+containerName {
				// Container exists, check if running
				if c.State == "running" {
					log.Printf("Container %s is already running", containerName)
					return nil
				}
				// Start existing container
				log.Printf("Starting existing container %s", containerName)
				return d.client.ContainerStart(ctx, c.ID, container.StartOptions{})
			}
		}
	}

	// Pull the image if not present
	log.Printf("Pulling image %s...", minecraftImage)
	reader, err := d.client.ImagePull(ctx, minecraftImage, image.PullOptions{})
	if err != nil {
		return fmt.Errorf("failed to pull image: %w", err)
	}
	defer reader.Close()
	// Drain the reader to complete the pull
	buf := make([]byte, 1024)
	for {
		_, err := reader.Read(buf)
		if err != nil {
			break
		}
	}
	log.Printf("Image %s pulled successfully", minecraftImage)

	// Prepare environment variables
	env := []string{
		"EULA=TRUE",
		fmt.Sprintf("VERSION=%s", server.Version),
		fmt.Sprintf("TYPE=%s", getServerType(server.Type)),
		fmt.Sprintf("SERVER_NAME=%s", server.Name),
		"MEMORY=2G",
		"ENABLE_RCON=true",
		"RCON_PASSWORD=minecraft",
	}

	// Prepare port bindings
	containerPort := nat.Port("25565/tcp")
	hostBinding := nat.PortBinding{
		HostIP:   "0.0.0.0",
		HostPort: strconv.Itoa(hostPort),
	}
	portBindings := nat.PortMap{
		containerPort: []nat.PortBinding{hostBinding},
	}

	// Prepare volume mounts
	dataPath := fmt.Sprintf("%s/%s", d.dataDir, server.ID)
	binds := []string{
		fmt.Sprintf("%s:/data", dataPath),
	}

	// Create container configuration
	config := &container.Config{
		Image: minecraftImage,
		Env:   env,
		ExposedPorts: nat.PortSet{
			containerPort: struct{}{},
		},
		Labels: map[string]string{
			"lotus.server.id":   server.ID,
			"lotus.server.name": server.Name,
			"lotus.server.type": server.Type,
		},
	}

	hostConfig := &container.HostConfig{
		PortBindings: portBindings,
		Binds:        binds,
		RestartPolicy: container.RestartPolicy{
			Name: container.RestartPolicyUnlessStopped,
		},
		Resources: container.Resources{
			Memory:   4 * 1024 * 1024 * 1024, // 4GB
			NanoCPUs: 1 * 1e9,                // 1 CPU core
		},
	}

	networkConfig := &network.NetworkingConfig{}

	// Create the container
	log.Printf("Creating container %s for server %s", containerName, server.ID)
	resp, err := d.client.ContainerCreate(ctx, config, hostConfig, networkConfig, nil, containerName)
	if err != nil {
		return fmt.Errorf("failed to create container: %w", err)
	}

	// Start the container
	log.Printf("Starting container %s (ID: %s)", containerName, resp.ID[:12])
	if err := d.client.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return fmt.Errorf("failed to start container: %w", err)
	}

	log.Printf("Server %s started successfully on port %d", server.ID, hostPort)
	return nil
}

// StopServer stops and optionally removes a Minecraft server container
func (d *DockerService) StopServer(ctx context.Context, serverID string, remove bool) error {
	containerName := containerPrefix + serverID

	// Find the container
	containers, err := d.client.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		return fmt.Errorf("failed to list containers: %w", err)
	}

	var containerID string
	for _, c := range containers {
		for _, name := range c.Names {
			if name == "/"+containerName {
				containerID = c.ID
				break
			}
		}
	}

	if containerID == "" {
		log.Printf("Container %s not found, nothing to stop", containerName)
		return nil
	}

	// Stop the container with a timeout
	log.Printf("Stopping container %s", containerName)
	timeout := 30 // seconds
	if err := d.client.ContainerStop(ctx, containerID, container.StopOptions{Timeout: &timeout}); err != nil {
		return fmt.Errorf("failed to stop container: %w", err)
	}

	log.Printf("Container %s stopped", containerName)

	// Optionally remove the container
	if remove {
		log.Printf("Removing container %s", containerName)
		if err := d.client.ContainerRemove(ctx, containerID, container.RemoveOptions{}); err != nil {
			return fmt.Errorf("failed to remove container: %w", err)
		}
		log.Printf("Container %s removed", containerName)
	}

	return nil
}

// RestartServer restarts a Minecraft server container
func (d *DockerService) RestartServer(ctx context.Context, serverID string) error {
	containerName := containerPrefix + serverID

	// Find the container
	containers, err := d.client.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		return fmt.Errorf("failed to list containers: %w", err)
	}

	var containerID string
	for _, c := range containers {
		for _, name := range c.Names {
			if name == "/"+containerName {
				containerID = c.ID
				break
			}
		}
	}

	if containerID == "" {
		return fmt.Errorf("container %s not found", containerName)
	}

	log.Printf("Restarting container %s", containerName)
	timeout := 30 // seconds
	if err := d.client.ContainerRestart(ctx, containerID, container.StopOptions{Timeout: &timeout}); err != nil {
		return fmt.Errorf("failed to restart container: %w", err)
	}

	log.Printf("Container %s restarted", containerName)
	return nil
}

// GetContainerStatus returns the status of a server container
func (d *DockerService) GetContainerStatus(ctx context.Context, serverID string) (string, error) {
	containerName := containerPrefix + serverID

	containers, err := d.client.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		return "", fmt.Errorf("failed to list containers: %w", err)
	}

	for _, c := range containers {
		for _, name := range c.Names {
			if name == "/"+containerName {
				return c.State, nil
			}
		}
	}

	return "not_found", nil
}

// getServerType maps our server type to itzg/minecraft-server TYPE env var
func getServerType(serverType string) string {
	switch serverType {
	case "vanilla":
		return "VANILLA"
	case "paper":
		return "PAPER"
	case "spigot":
		return "SPIGOT"
	case "fabric":
		return "FABRIC"
	case "forge":
		return "FORGE"
	case "bukkit":
		return "BUKKIT"
	default:
		return "VANILLA"
	}
}
