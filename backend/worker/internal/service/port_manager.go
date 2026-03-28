package service

import (
	"fmt"
	"sync"
)

// PortManager handles port allocation for game servers
// It ensures each server gets a unique port and tracks allocations
type PortManager struct {
	mu          sync.RWMutex
	basePort    int
	maxPorts    int
	allocations map[string]int // serverID -> port
	usedPorts   map[int]string // port -> serverID
}

// NewPortManager creates a new port manager
func NewPortManager(basePort, maxServers int) *PortManager {
	return &PortManager{
		basePort:    basePort,
		maxPorts:    maxServers,
		allocations: make(map[string]int),
		usedPorts:   make(map[int]string),
	}
}

// AllocatePort assigns a port to a server
// Returns the allocated port or an error if no ports are available
func (pm *PortManager) AllocatePort(serverID string) (int, error) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	// Check if server already has a port allocated
	if port, exists := pm.allocations[serverID]; exists {
		return port, nil
	}

	// Find the next available port
	for i := 0; i < pm.maxPorts; i++ {
		port := pm.basePort + i
		if _, used := pm.usedPorts[port]; !used {
			// Allocate this port
			pm.allocations[serverID] = port
			pm.usedPorts[port] = serverID
			return port, nil
		}
	}

	return 0, fmt.Errorf("no available ports (all %d ports in use)", pm.maxPorts)
}

// ReleasePort frees a port previously allocated to a server
func (pm *PortManager) ReleasePort(serverID string) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	if port, exists := pm.allocations[serverID]; exists {
		delete(pm.allocations, serverID)
		delete(pm.usedPorts, port)
	}
}

// GetPort returns the port allocated to a server, or 0 if not allocated
func (pm *PortManager) GetPort(serverID string) int {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	return pm.allocations[serverID]
}

// IsPortAllocated checks if a specific port is in use
func (pm *PortManager) IsPortAllocated(port int) bool {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	_, exists := pm.usedPorts[port]
	return exists
}

// GetAllocatedCount returns the number of currently allocated ports
func (pm *PortManager) GetAllocatedCount() int {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	return len(pm.allocations)
}

// GetAllAllocations returns a copy of all current allocations
func (pm *PortManager) GetAllAllocations() map[string]int {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	result := make(map[string]int, len(pm.allocations))
	for k, v := range pm.allocations {
		result[k] = v
	}
	return result
}
