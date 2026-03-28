package service

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/redis/go-redis/v9"
)

const (
	// Redis Pub/Sub channel for server logs
	logChannelPrefix = "logs:server:"
	// Redis list key for storing recent logs
	logHistoryPrefix = "logs:history:"
	// Maximum number of log lines to store in history
	maxLogHistory = 200
)

// LogStreamer manages log streaming from Docker containers to Redis Pub/Sub
type LogStreamer struct {
	dockerClient *client.Client
	redisClient  *redis.Client
	activeStreams map[string]context.CancelFunc
	mu           sync.Mutex
}

// NewLogStreamer creates a new log streamer
func NewLogStreamer(dockerClient *client.Client, redisClient *redis.Client) *LogStreamer {
	return &LogStreamer{
		dockerClient:  dockerClient,
		redisClient:   redisClient,
		activeStreams: make(map[string]context.CancelFunc),
	}
}

// StartStreaming starts streaming logs for a server to Redis Pub/Sub
func (l *LogStreamer) StartStreaming(serverID string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Check if already streaming
	if _, exists := l.activeStreams[serverID]; exists {
		log.Printf("Log streaming already active for server %s", serverID)
		return nil
	}

	containerName := containerPrefix + serverID
	channel := logChannelPrefix + serverID

	// Create cancellable context for this stream
	ctx, cancel := context.WithCancel(context.Background())
	l.activeStreams[serverID] = cancel

	go l.streamLogs(ctx, containerName, channel, serverID)

	log.Printf("Started log streaming for server %s to channel %s", serverID, channel)
	return nil
}

// StopStreaming stops streaming logs for a server
func (l *LogStreamer) StopStreaming(serverID string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if cancel, exists := l.activeStreams[serverID]; exists {
		cancel()
		delete(l.activeStreams, serverID)
		log.Printf("Stopped log streaming for server %s", serverID)
	}
}

// StopAll stops all active log streams
func (l *LogStreamer) StopAll() {
	l.mu.Lock()
	defer l.mu.Unlock()

	for serverID, cancel := range l.activeStreams {
		cancel()
		log.Printf("Stopped log streaming for server %s", serverID)
	}
	l.activeStreams = make(map[string]context.CancelFunc)
}

func (l *LogStreamer) streamLogs(ctx context.Context, containerName, channel, serverID string) {
	historyKey := logHistoryPrefix + serverID

	// Clear old history when starting fresh
	l.redisClient.Del(ctx, historyKey)

	// Find container ID
	containers, err := l.dockerClient.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		log.Printf("Failed to list containers for log streaming: %v", err)
		l.publishError(channel, "Failed to access container")
		return
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
		log.Printf("Container %s not found for log streaming", containerName)
		l.publishError(channel, "[SERVER_NOT_RUNNING]")
		return
	}

	// Get container logs with follow
	options := container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
		Tail:       "100", // Get last 100 lines initially
		Timestamps: false,
	}

	reader, err := l.dockerClient.ContainerLogs(ctx, containerID, options)
	if err != nil {
		log.Printf("Failed to get container logs: %v", err)
		l.publishError(channel, "Failed to read container logs")
		return
	}
	defer reader.Close()

	// Read and publish logs line by line
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return
		default:
			line := scanner.Text()
			// Docker multiplexed stream has 8-byte header, skip it if present
			if len(line) > 8 {
				// Check if first byte is 1 (stdout) or 2 (stderr)
				if line[0] == 1 || line[0] == 2 {
					line = line[8:]
				}
			}

			if line != "" {
				// Store in history list (RPUSH to append, LTRIM to keep only recent)
				pipe := l.redisClient.Pipeline()
				pipe.RPush(ctx, historyKey, line)
				pipe.LTrim(ctx, historyKey, -maxLogHistory, -1) // Keep only last N lines
				_, err := pipe.Exec(ctx)
				if err != nil {
					log.Printf("Failed to store log in history: %v", err)
				}

				// Publish to Pub/Sub for real-time streaming
				err = l.redisClient.Publish(ctx, channel, line).Err()
				if err != nil {
					log.Printf("Failed to publish log to Redis: %v", err)
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		if ctx.Err() == nil {
			log.Printf("Error reading container logs: %v", err)
		}
	}

	// Notify that container stopped
	l.redisClient.Publish(context.Background(), channel, "[SERVER_STOPPED]")
}

func (l *LogStreamer) publishError(channel, message string) {
	ctx := context.Background()
	l.redisClient.Publish(ctx, channel, message)
}

// GetLogChannel returns the Redis Pub/Sub channel name for a server
func GetLogChannel(serverID string) string {
	return fmt.Sprintf("%s%s", logChannelPrefix, serverID)
}

