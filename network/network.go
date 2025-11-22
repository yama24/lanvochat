package network

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

// Message represents a chat message structure
type Message struct {
	Type      string    `json:"type"`
	PeerID    string    `json:"peer_id"`
	SenderID  string    `json:"sender_id"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

// DiscoveryMessage represents a peer discovery message
type DiscoveryMessage struct {
	Type   string `json:"type"`
	PeerID string `json:"peer_id"`
	Name   string `json:"name"`
	IP     string `json:"ip"`
	Port   int    `json:"port"`
	Status string `json:"status"`
}

// NetworkManager handles all network operations
type NetworkManager struct {
	ctx           context.Context
	db            interface{}
	localPeerID   string
	localName     string
	localIP       string
	multicastAddr string
	tcpPort       int
	udpPort       int

	multicastConn *net.UDPConn
	udpConn       *net.UDPConn
	tcpListener   *net.TCPListener
	tcpAddr       *net.TCPAddr

	stopChan chan bool
	wg       sync.WaitGroup

	activePeers map[string]*PeerInfo
	peersMutex  sync.RWMutex
}

// PeerInfo holds information about discovered peers
type PeerInfo struct {
	PeerID   string
	Name     string
	IP       string
	Port     int
	LastSeen time.Time
	IsOnline bool
}

// NewNetworkManager creates a new network manager
func NewNetworkManager(peerID, name, localIP string) *NetworkManager {
	return &NetworkManager{
		localPeerID:   peerID,
		localName:     name,
		localIP:       localIP,
		multicastAddr: "239.255.255.250:1900",
		tcpPort:       8080,
		udpPort:       8081,
		stopChan:      make(chan bool),
		activePeers:   make(map[string]*PeerInfo),
	}
}

// SetContext sets the Wails context for event emission
func (nm *NetworkManager) SetContext(ctx context.Context) {
	nm.ctx = ctx
}

// SetDatabase sets the database interface
func (nm *NetworkManager) SetDatabase(db interface{}) {
	nm.db = db
}

// Start begins all network operations
func (nm *NetworkManager) Start() error {
	log.Println("Starting network manager...")

	if err := nm.startMulticastDiscovery(); err != nil {
		return fmt.Errorf("failed to start multicast discovery: %w", err)
	}

	if err := nm.startTCPListener(); err != nil {
		return fmt.Errorf("failed to start TCP listener: %w", err)
	}

	go nm.peerCleanupRoutine()

	log.Println("Network manager started successfully")
	return nil
}

// Stop gracefully stops all network operations
func (nm *NetworkManager) Stop() {
	log.Println("Stopping network manager...")
	close(nm.stopChan)
	nm.wg.Wait()

	if nm.multicastConn != nil {
		nm.multicastConn.Close()
	}
	if nm.tcpListener != nil {
		nm.tcpListener.Close()
	}
	if nm.udpConn != nil {
		nm.udpConn.Close()
	}

	log.Println("Network manager stopped")
}
