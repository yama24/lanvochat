package network

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// startTCPListener starts the TCP listener for incoming messages
func (nm *NetworkManager) startTCPListener() error {
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", nm.tcpPort))
	if err != nil {
		return fmt.Errorf("failed to resolve TCP address: %w", err)
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to start TCP listener: %w", err)
	}

	nm.tcpListener = listener
	nm.tcpAddr = addr

	// Start accepting connections
	nm.wg.Add(1)
	go nm.tcpAcceptRoutine()

	log.Printf("TCP listener started on port %d", nm.tcpPort)
	return nil
}

// tcpAcceptRoutine accepts incoming TCP connections
func (nm *NetworkManager) tcpAcceptRoutine() {
	defer nm.wg.Done()

	for {
		select {
		case <-nm.stopChan:
			return
		default:
			nm.tcpListener.SetDeadline(time.Now().Add(time.Second))
			conn, err := nm.tcpListener.AcceptTCP()
			if err != nil {
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					continue
				}
				log.Printf("Error accepting TCP connection: %v", err)
				continue
			}

			// Handle connection in separate goroutine
			go nm.handleTCPConnection(conn)
		}
	}
}

// handleTCPConnection handles an incoming TCP connection
func (nm *NetworkManager) handleTCPConnection(conn *net.TCPConn) {
	defer conn.Close()

	// Set timeouts
	conn.SetReadDeadline(time.Now().Add(30 * time.Second))
	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

	buffer := make([]byte, 4096)
	n, err := conn.Read(buffer)
	if err != nil {
		log.Printf("Error reading from TCP connection: %v", err)
		return
	}

	// Parse message
	var msg Message
	if err := json.Unmarshal(buffer[:n], &msg); err != nil {
		log.Printf("Error parsing TCP message: %v", err)
		return
	}

	// Process message
	nm.processIncomingMessage(msg)
}

// processIncomingMessage processes an incoming message
func (nm *NetworkManager) processIncomingMessage(msg Message) {
	log.Printf("Received message from %s: %s", msg.SenderID, msg.Content)

	// Save to database if available
	if nm.db != nil {
		// This would call the database SaveMessage function
		// For now, we'll just log it
		log.Printf("Would save message to database: %+v", msg)
	}

	// Emit event to frontend
	if nm.ctx != nil {
		runtime.EventsEmit(nm.ctx, "messageReceived", msg)
	}
}

// SendMessageToPeer sends a message to a specific peer via TCP
func (nm *NetworkManager) SendMessageToPeer(peerID, content string) error {
	nm.peersMutex.RLock()
	peer, exists := nm.activePeers[peerID]
	nm.peersMutex.RUnlock()

	if !exists {
		return fmt.Errorf("peer %s not found", peerID)
	}

	// Create message
	msg := Message{
		Type:      "message",
		PeerID:    peerID,
		SenderID:  nm.localPeerID,
		Content:   content,
		Timestamp: time.Now(),
	}

	// Marshal to JSON
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	// Connect to peer
	addr := fmt.Sprintf("%s:%d", peer.IP, peer.Port)
	conn, err := net.DialTimeout("tcp", addr, 5*time.Second)
	if err != nil {
		return fmt.Errorf("failed to connect to peer %s: %w", peerID, err)
	}
	defer conn.Close()

	// Set write deadline
	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

	// Send message
	_, err = conn.Write(data)
	if err != nil {
		return fmt.Errorf("failed to send message to peer %s: %w", peerID, err)
	}

	// Ensure data is sent
	if tcpConn, ok := conn.(*net.TCPConn); ok {
		tcpConn.CloseWrite()
	}

	log.Printf("Message sent to peer %s (%s)", peer.Name, peerID)
	return nil
}

// GetActivePeers returns a copy of active peers
func (nm *NetworkManager) GetActivePeers() map[string]*PeerInfo {
	nm.peersMutex.RLock()
	defer nm.peersMutex.RUnlock()

	peers := make(map[string]*PeerInfo)
	for k, v := range nm.activePeers {
		peerCopy := *v // Shallow copy
		peers[k] = &peerCopy
	}

	return peers
}

// BroadcastMessage sends a message to all active peers
func (nm *NetworkManager) BroadcastMessage(content string) error {
	nm.peersMutex.RLock()
	peers := make([]*PeerInfo, 0, len(nm.activePeers))
	for _, peer := range nm.activePeers {
		peerCopy := *peer
		peers = append(peers, &peerCopy)
	}
	nm.peersMutex.RUnlock()

	var lastErr error
	for _, peer := range peers {
		if err := nm.SendMessageToPeer(peer.PeerID, content); err != nil {
			log.Printf("Failed to send to peer %s: %v", peer.PeerID, err)
			lastErr = err
		}
	}

	return lastErr
}
