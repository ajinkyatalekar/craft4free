package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"

	"host-lotus-worker/internal/domain"
	"host-lotus-worker/internal/repository"
)

type AssignmentProcessor struct {
	redisClient       *redis.Client
	serverRepo        *repository.ServerRepository
	dockerService     *DockerService
	portManager       *PortManager
	logStreamer       *LogStreamer
	cloudflareService *CloudflareService
	worker            *domain.Worker
	interval          time.Duration
	ctx               context.Context
	cancel            context.CancelFunc
}

func NewAssignmentProcessor(redisClient *redis.Client, worker *domain.Worker, serverRepo *repository.ServerRepository, dockerService *DockerService, portManager *PortManager, logStreamer *LogStreamer, cloudflareService *CloudflareService, interval time.Duration) *AssignmentProcessor {
	ctx, cancel := context.WithCancel(context.Background())
	return &AssignmentProcessor{
		redisClient:       redisClient,
		serverRepo:        serverRepo,
		dockerService:     dockerService,
		portManager:       portManager,
		logStreamer:       logStreamer,
		cloudflareService: cloudflareService,
		worker:            worker,
		interval:          interval,
		ctx:               ctx,
		cancel:            cancel,
	}
}

// Start begins the assignment processor background task
func (a *AssignmentProcessor) Start() {
	go a.run()
	log.Printf("Assignment processor started for worker %s with interval %v", a.worker.ID, a.interval)
}

// Stop gracefully stops the assignment processor background task
func (a *AssignmentProcessor) Stop() {
	a.cancel()
	log.Println("Assignment processor stopped")
}

func (a *AssignmentProcessor) run() {
	// Wait for heartbeat service to start and set worker state
	time.Sleep(a.interval)

	ticker := time.NewTicker(a.interval)
	defer ticker.Stop()

	if err := a.checkAssignments(); err != nil {
		log.Printf("Failed to check initial assignments: %v", err)
	}

	for {
		select {
		case <-ticker.C:
			if err := a.checkAssignments(); err != nil {
				log.Printf("Failed to check assignments: %v", err)
			}
		case <-a.ctx.Done():
			log.Println("Assignment processor shutting down")
			return
		}
	}
}

// worker/internal/service/assignment_processor.go
const processAssignmentScript = `
    local key = KEYS[1]
    local serverID = ARGV[1]
    
    local data = redis.call('GET', key)
    if not data then
        return nil
    end
    
    local worker = cjson.decode(data)
    local assigned_servers = worker.assigned_servers or {}
    local running_servers = worker.running_servers or {}
    
    -- Check if this server is in assignments
    if not assigned_servers[serverID] then
        return redis.error_reply("Server not in assignments")
    end
    
    -- Move from assigned to running
    assigned_servers[serverID] = nil
    running_servers[serverID] = serverID
    
    worker.assigned_servers = assigned_servers
    worker.running_servers = running_servers
    
    redis.call('SET', key, cjson.encode(worker))
    return 1
`

// Removes serverID from assigned_servers without moving it to running (aborted assignment).
const removeFromAssignedScript = `
    local key = KEYS[1]
    local serverID = ARGV[1]
    local data = redis.call('GET', key)
    if not data then
        return nil
    end
    local worker = cjson.decode(data)
    local assigned_servers = worker.assigned_servers or {}
    assigned_servers[serverID] = nil
    worker.assigned_servers = assigned_servers
    redis.call('SET', key, cjson.encode(worker))
    return 1
`

func (a *AssignmentProcessor) removeFromAssignedQueue(serverID string) error {
	script := redis.NewScript(removeFromAssignedScript)
	return script.Run(a.ctx, a.redisClient,
		[]string{fmt.Sprintf(workerKey, a.worker.ID)},
		serverID,
	).Err()
}

func (a *AssignmentProcessor) processAssignment(serverID string) error {
	log.Printf("Starting server %s", serverID)

	server, err := a.serverRepo.GetServerById(a.ctx, serverID)
	if err != nil {
		return fmt.Errorf("failed to get server: %w", err)
	}

	// Update status to starting
	server.Status = domain.StatusStarting
	err = a.serverRepo.UpdateServer(a.ctx, server)
	if err != nil {
		return fmt.Errorf("failed to update server status to starting: %w", err)
	}

	// Allocate a port for this server
	hostPort, err := a.portManager.AllocatePort(serverID)
	if err != nil {
		server.Status = domain.StatusStopped
		_ = a.serverRepo.UpdateServer(a.ctx, server)
		if rmErr := a.removeFromAssignedQueue(serverID); rmErr != nil {
			log.Printf("Failed to remove server %s from assignment queue after port allocation failure: %v", serverID, rmErr)
		}
		log.Printf("Assignment aborted for server %s: %v", serverID, err)
		return nil
	}

	// Start the Docker container
	if err := a.dockerService.StartServer(a.ctx, &server, hostPort); err != nil {
		a.portManager.ReleasePort(serverID)
		server.Status = domain.StatusStopped
		_ = a.serverRepo.UpdateServer(a.ctx, server)
		if rmErr := a.removeFromAssignedQueue(serverID); rmErr != nil {
			log.Printf("Failed to remove server %s from assignment queue after start failure: %v", serverID, rmErr)
		}
		log.Printf("Assignment aborted for server %s (not started): %v", serverID, err)
		return nil
	}

	// Start log streaming for this server
	if err := a.logStreamer.StartStreaming(serverID); err != nil {
		log.Printf("Warning: failed to start log streaming for server %s: %v", serverID, err)
	}

	// Create DNS SRV record for this server (if Cloudflare is configured)
	if a.cloudflareService != nil {
		if err := a.cloudflareService.CreateMinecraftSRVRecord(a.ctx, server.Name, a.worker.OciId, hostPort); err != nil {
			log.Printf("Warning: failed to create DNS record for server %s: %v", serverID, err)
			// Don't fail the server start if DNS creation fails
		}
	}

	// Update server with running status and IP
	server.Status = domain.StatusRunning
	server.IP = fmt.Sprintf("%s:%d", a.worker.PublicIp, hostPort)
	err = a.serverRepo.UpdateServer(a.ctx, server)
	if err != nil {
		return fmt.Errorf("failed to update server status to running: %w", err)
	}

	// Then atomically move from assigned to running in Redis
	a.worker.RunningServers[serverID] = serverID

	script := redis.NewScript(processAssignmentScript)
	return script.Run(a.ctx, a.redisClient,
		[]string{fmt.Sprintf(workerKey, a.worker.ID)},
		serverID,
	).Err()
}

func (a *AssignmentProcessor) checkAssignments() error {
	// Get current worker state from Redis
	key := fmt.Sprintf(workerKey, a.worker.ID)
	data, err := a.redisClient.Get(a.ctx, key).Result()
	if err != nil {
		return err
	}

	var worker domain.Worker
	json.Unmarshal([]byte(data), &worker)

	// Process any assigned servers
	for serverID := range worker.AssignedServers {
		// Start server and move to running
		err := a.processAssignment(serverID)
		if err != nil {
			log.Printf("Failed to process assignment: %v", err)
		}
	}

	return nil
}
