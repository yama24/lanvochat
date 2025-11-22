package main

import (
	"context"
	"fmt"
	"lanvochat/database"
	"lanvochat/network"
	"log"
	"math/rand"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx            context.Context
	db             *database.Database
	networkManager *network.NetworkManager
	localPeerID    string
	localName      string
}

// NewApp creates a new App application struct
func NewApp() *App {
	// Generate unique peer ID
	rand.Seed(time.Now().UnixNano())
	peerID := fmt.Sprintf("peer_%d", rand.Int63())

	return &App{
		localPeerID: peerID,
		localName:   "LanvoChat User", // Default name, can be changed later
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Initialize database
	db, err := database.NewDatabase("lanvochat.db")
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	a.db = db

	// Get local IP address
	localIP, err := a.getLocalIP()
	if err != nil {
		log.Printf("Warning: Could not determine local IP: %v", err)
		localIP = "127.0.0.1"
	}

	// Initialize network manager
	a.networkManager = network.NewNetworkManager(a.localPeerID, a.localName, localIP)
	a.networkManager.SetContext(ctx)
	a.networkManager.SetDatabase(a.db)

	// Start network operations
	if err := a.networkManager.Start(); err != nil {
		log.Printf("Warning: Failed to start network manager: %v", err)
	}

	fmt.Println("Database initialized successfully")
	fmt.Printf("Network manager started for peer: %s (%s)\n", a.localName, a.localPeerID)
}

// shutdown is called at application termination
func (a *App) shutdown(ctx context.Context) {
	if a.networkManager != nil {
		a.networkManager.Stop()
	}
	if a.db != nil {
		a.db.Close()
	}
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// SaveMessage saves a message to the database
func (a *App) SaveMessage(peerID, senderID, content string) error {
	return a.db.SaveMessage(peerID, senderID, content)
}

// GetMessageHistory retrieves message history for a peer
func (a *App) GetMessageHistory(peerID string, limit int) ([]database.Message, error) {
	return a.db.GetMessageHistory(peerID, limit)
}

// SavePeer saves a peer to the database
func (a *App) SavePeer(peerID, name, ipAddress string) error {
	return a.db.SavePeer(peerID, name, ipAddress)
}

// GetPeers retrieves all peers from the database
func (a *App) GetPeers() ([]database.Peer, error) {
	return a.db.GetPeers()
}

// ShowNotification displays a system notification
func (a *App) ShowNotification(title, message string) {
	runtime.EventsEmit(a.ctx, "notification", map[string]string{
		"title":   title,
		"message": message,
	})
}

// SendMessage sends a message to a specific peer
func (a *App) SendMessage(peerID, content string) error {
	if a.networkManager == nil {
		return fmt.Errorf("network manager not initialized")
	}
	return a.networkManager.SendMessageToPeer(peerID, content)
}

// BroadcastMessage sends a message to all peers
func (a *App) BroadcastMessage(content string) error {
	if a.networkManager == nil {
		return fmt.Errorf("network manager not initialized")
	}
	return a.networkManager.BroadcastMessage(content)
}

// GetActivePeers returns currently active peers
func (a *App) GetActivePeers() map[string]interface{} {
	if a.networkManager == nil {
		return make(map[string]interface{})
	}

	peers := a.networkManager.GetActivePeers()
	result := make(map[string]interface{})

	for k, v := range peers {
		result[k] = map[string]interface{}{
			"peer_id":   v.PeerID,
			"name":      v.Name,
			"ip":        v.IP,
			"port":      v.Port,
			"last_seen": v.LastSeen,
			"is_online": v.IsOnline,
		}
	}

	return result
}

// SetLocalName sets the local peer name
func (a *App) SetLocalName(name string) {
	a.localName = name
	if a.networkManager != nil {
		// Note: In a real implementation, you might want to restart
		// the network manager or send an update message
		log.Printf("Local name updated to: %s", name)
	}
}

// GetLocalPeerInfo returns local peer information
func (a *App) GetLocalPeerInfo() map[string]string {
	return map[string]string{
		"peer_id": a.localPeerID,
		"name":    a.localName,
	}
}

// getLocalIP attempts to get the local IP address
func (a *App) getLocalIP() (string, error) {
	// This is a simplified version. In production, you might want
	// to enumerate all interfaces and choose the most appropriate one
	return "192.168.1.100", nil // Placeholder - would need proper implementation
}
