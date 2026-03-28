package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"host-lotus-worker/internal/config"
	"host-lotus-worker/internal/domain"
	"host-lotus-worker/internal/repository"
	"host-lotus-worker/internal/service"
)

func main() {
	// Load .env
	_ = godotenv.Load()

	// Load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize Supabase client
	supabaseClient, err := service.NewSupabaseClient(cfg.SupabaseURL, cfg.SupabaseServiceKey)
	if err != nil {
		log.Fatalf("Failed to initialize Supabase client: %v", err)
	}
	log.Println("Supabase client initialized successfully")

	// Initialize Redis client
	redisClient, err := service.NewRedisClient(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB)
	if err != nil {
		log.Fatalf("Failed to initialize Redis client: %v", err)
	}
	defer redisClient.Close()
	log.Println("Redis client initialized successfully")

	serverRepo := repository.NewServerRepository(redisClient, supabaseClient)

	// Initialize Docker service
	dockerService, err := service.NewDockerService(cfg.DockerDataDir, cfg.DockerBasePort)
	if err != nil {
		log.Fatalf("Failed to initialize Docker service: %v", err)
	}
	defer dockerService.Close()
	log.Println("Docker service initialized successfully")

	// Initialize port manager for tracking server port allocations
	portManager := service.NewPortManager(cfg.DockerBasePort, cfg.MaxServers)
	log.Printf("Port manager initialized (base port: %d, max servers: %d)", cfg.DockerBasePort, cfg.MaxServers)

	// Initialize log streamer for streaming Docker logs to Redis Pub/Sub
	logStreamer := service.NewLogStreamer(dockerService.GetClient(), redisClient)
	log.Println("Log streamer initialized successfully")

	// Initialize Cloudflare service (optional - for DNS record management)
	var cloudflareService *service.CloudflareService
	if cfg.CloudflareAPIToken != "" && cfg.CloudflareZoneID != "" {
		cloudflareService = service.NewCloudflareService(cfg.CloudflareAPIToken, cfg.CloudflareZoneID, cfg.CloudflareZoneName)
		log.Printf("Cloudflare DNS service initialized for zone %s", cfg.CloudflareZoneName)
	} else {
		log.Println("Cloudflare DNS service not configured (CLOUDFLARE_API_TOKEN or CLOUDFLARE_ZONE_ID missing)")
	}

	worker := &domain.Worker{
		ID:              cfg.ID,
		OciId:           cfg.OciId,
		CreatedAt:       cfg.CreatedAt,
		PublicIp:        cfg.PublicIp,
		PrivateIp:       cfg.PrivateIp,
		MaxServers:      cfg.MaxServers,
		RunningServers:  map[string]string{},
		AssignedServers: map[string]string{},
	}

	// Initialize and start heartbeat service
	heartbeatInterval := 1 * time.Second
	heartbeatService := service.NewHeartbeatService(redisClient, worker, heartbeatInterval)
	heartbeatService.Start()

	assignmentInterval := 500 * time.Millisecond
	assignmentProcessor := service.NewAssignmentProcessor(redisClient, worker, serverRepo, dockerService, portManager, logStreamer, cloudflareService, assignmentInterval)
	assignmentProcessor.Start()

	// Add command processor
	commandProcessor := service.NewCommandProcessor(redisClient, worker, serverRepo, dockerService, portManager, logStreamer, cloudflareService)
	commandProcessor.Start()

	log.Printf("Worker %s started successfully", worker.ID)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down worker...")
	heartbeatService.Stop()
	assignmentProcessor.Stop()
	commandProcessor.Stop()
	logStreamer.StopAll() // Stop all log streams
	_ = supabaseClient
	log.Println("Worker shut down successfully")
}
