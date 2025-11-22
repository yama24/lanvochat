package main

import (
	"context"
	"fmt"
	"lanvochat/database"
	"log"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
	db  *database.Database
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
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

	fmt.Println("Database initialized successfully")
}

// shutdown is called at application termination
func (a *App) shutdown(ctx context.Context) {
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
