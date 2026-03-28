package handlers

import (
	"host-lotus-controller/internal/domain"
	"host-lotus-controller/internal/repository"
	"host-lotus-controller/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ServerHandler struct {
	serverRepo    *repository.ServerRepository
	workerRepo    *repository.WorkerRepository
	serverService *service.ServerService
}

func NewServerHandler(serverRepo *repository.ServerRepository, workerRepo *repository.WorkerRepository, serverService *service.ServerService) *ServerHandler {
	return &ServerHandler{
		serverRepo:    serverRepo,
		workerRepo:    workerRepo,
		serverService: serverService,
	}
}

// Add this struct for JSON binding
type CreateServerRequest struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Version string `json:"version"`
}

// CreateServer creates a new server in the persistent database.
// This does not start or send any RPCs to workers.
// TODO: Input validation
func (h *ServerHandler) CreateServer(c *gin.Context) {
	user_id, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	var req CreateServerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"success": false, "error": "Invalid request body: " + err.Error()})
		return
	}

	server := domain.ServerPersistent{
		ID:      uuid.New().String(),
		UserId:  user_id.(string),
		Name:    req.Name,
		Type:    req.Type,
		Version: req.Version,
	}

	err := h.serverRepo.CreateServer(c.Request.Context(), server)
	if err != nil {
		c.JSON(500, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"success": true, "data": server})
}

// GetAllUserServers returns all servers for a user.
// If dynamic row exists, sends dynamic data. Otherwise, sends persistent data.
func (h *ServerHandler) GetAllUserServers(c *gin.Context) {
	user_id, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	servers, err := h.serverRepo.GetAllServers(c.Request.Context(), user_id.(string))

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"success": true, "data": servers})
}

// GetServerById returns a server by ID.
// If dynamic row exists, sends dynamic data. Otherwise, sends persistent data.
func (h *ServerHandler) GetServerById(c *gin.Context) {
	user_id, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	server_id := c.Param("server_id")

	server, err := h.serverRepo.GetServerById(c.Request.Context(), user_id.(string), server_id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"success": true, "data": server[0]})
}

// StartServer starts a server by ID.
func (h *ServerHandler) StartServer(c *gin.Context) {
	userID, _ := c.Get("user_id")
	serverID := c.Param("server_id")
	err := h.serverService.StartServer(c.Request.Context(), userID.(string), serverID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"success": true})
}

func (h *ServerHandler) StopServer(c *gin.Context) {
	userID, _ := c.Get("user_id")
	serverID := c.Param("server_id")
	err := h.serverService.StopServer(c.Request.Context(), userID.(string), serverID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"success": true})
}
