package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"

	"host-lotus-controller/internal/repository"
)

const (
	// Redis Pub/Sub channel prefix for server logs
	logChannelPrefix = "logs:server:"
	// Redis list key prefix for stored log history
	logHistoryPrefix = "logs:history:"

	// WebSocket configuration
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins for development
		// TODO: Configure allowed origins for production
		return true
	},
}

type ConsoleHandler struct {
	redisClient *redis.Client
	serverRepo  *repository.ServerRepository
}

func NewConsoleHandler(redisClient *redis.Client, serverRepo *repository.ServerRepository) *ConsoleHandler {
	return &ConsoleHandler{
		redisClient: redisClient,
		serverRepo:  serverRepo,
	}
}

// HandleConsoleWebSocket handles WebSocket connections for server console logs
func (h *ConsoleHandler) HandleConsoleWebSocket(c *gin.Context) {
	serverID := c.Param("server_id")

	// Upgrade to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade WebSocket connection: %v", err)
		return
	}
	defer conn.Close()

	// Set up WebSocket parameters
	conn.SetReadLimit(maxMessageSize)
	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	// Create context for this connection
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Subscribe to Redis Pub/Sub channel for this server's logs
	channel := logChannelPrefix + serverID
	pubsub := h.redisClient.Subscribe(ctx, channel)
	defer pubsub.Close()

	// Wait for subscription to be ready
	_, err = pubsub.Receive(ctx)
	if err != nil {
		log.Printf("Failed to subscribe to Redis channel %s: %v", channel, err)
		conn.WriteMessage(websocket.TextMessage, []byte("[CONNECTION_ERROR]"))
		return
	}

	log.Printf("Client connected to console for server %s", serverID)

	// Fetch and send historical logs first
	historyKey := logHistoryPrefix + serverID
	historicalLogs, err := h.redisClient.LRange(ctx, historyKey, 0, -1).Result()
	if err != nil {
		log.Printf("Failed to fetch log history: %v", err)
	} else if len(historicalLogs) > 0 {
		log.Printf("Sending %d historical log lines to client", len(historicalLogs))
		for _, logLine := range historicalLogs {
			conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := conn.WriteMessage(websocket.TextMessage, []byte(logLine)); err != nil {
				log.Printf("Failed to send historical log: %v", err)
				return
			}
		}
	}

	// Channel to receive messages from Redis
	msgChan := pubsub.Channel()

	// Channel to signal when WebSocket is closed
	done := make(chan struct{})

	// Goroutine to read from WebSocket (handles client disconnect)
	go func() {
		defer close(done)
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("WebSocket error: %v", err)
				}
				return
			}
		}
	}()

	// Ping ticker to keep connection alive
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	// Main loop: forward Redis messages to WebSocket
	for {
		select {
		case msg, ok := <-msgChan:
			if !ok {
				// Channel closed
				conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := conn.WriteMessage(websocket.TextMessage, []byte(msg.Payload)); err != nil {
				log.Printf("Failed to write to WebSocket: %v", err)
				return
			}

		case <-ticker.C:
			conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}

		case <-done:
			// Client disconnected
			log.Printf("Client disconnected from console for server %s", serverID)
			return

		case <-ctx.Done():
			return
		}
	}
}

