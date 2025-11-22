package network

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// startMulticastDiscovery starts UDP multicast discovery
func (nm *NetworkManager) startMulticastDiscovery() error {
	// Parse multicast address
	addr, err := net.ResolveUDPAddr("udp", nm.multicastAddr)
	if err != nil {
		return fmt.Errorf("failed to resolve multicast address: %w", err)
	}

	// Create UDP connection for listening
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on multicast address: %w", err)
	}

	nm.multicastConn = conn

	// Join multicast group
	if err := nm.joinMulticastGroup(); err != nil {
		conn.Close()
		return fmt.Errorf("failed to join multicast group: %w", err)
	}

	// Start listening goroutine
	nm.wg.Add(1)
	go nm.multicastListenRoutine()

	// Start broadcasting goroutine
	nm.wg.Add(1)
	go nm.multicastBroadcastRoutine()

	log.Printf("Multicast discovery started on %s", nm.multicastAddr)
	return nil
}

// joinMulticastGroup joins the multicast group
func (nm *NetworkManager) joinMulticastGroup() error {
	// For simplicity, we'll skip explicit multicast group joining
	// and rely on the system routing. In production, you might want
	// to use socket options to properly join the multicast group.
	return nil
}

// multicastListenRoutine listens for incoming multicast messages
func (nm *NetworkManager) multicastListenRoutine() {
	defer nm.wg.Done()

	buffer := make([]byte, 2048)

	for {
		select {
		case <-nm.stopChan:
			return
		default:
			nm.multicastConn.SetReadDeadline(time.Now().Add(time.Second))
			n, srcAddr, err := nm.multicastConn.ReadFromUDP(buffer)
			if err != nil {
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					continue
				}
				log.Printf("Error reading multicast: %v", err)
				continue
			}

			// Parse message
			var msg DiscoveryMessage
			if err := json.Unmarshal(buffer[:n], &msg); err != nil {
				continue // Ignore invalid messages
			}

			// Skip our own messages
			if msg.PeerID == nm.localPeerID {
				continue
			}

			// Update peer information
			nm.updatePeerInfo(msg, srcAddr.IP.String())
		}
	}
}

// multicastBroadcastRoutine broadcasts our presence periodically
func (nm *NetworkManager) multicastBroadcastRoutine() {
	defer nm.wg.Done()

	ticker := time.NewTicker(30 * time.Second) // Broadcast every 30 seconds
	defer ticker.Stop()

	for {
		select {
		case <-nm.stopChan:
			return
		case <-ticker.C:
			nm.broadcastPresence()
		}
	}
}

// broadcastPresence sends a discovery message
func (nm *NetworkManager) broadcastPresence() {
	msg := DiscoveryMessage{
		Type:   "discovery",
		PeerID: nm.localPeerID,
		Name:   nm.localName,
		IP:     nm.localIP,
		Port:   nm.tcpPort,
		Status: "online",
	}

	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error marshaling discovery message: %v", err)
		return
	}

	addr, err := net.ResolveUDPAddr("udp", nm.multicastAddr)
	if err != nil {
		log.Printf("Error resolving multicast address: %v", err)
		return
	}

	_, err = nm.multicastConn.WriteToUDP(data, addr)
	if err != nil {
		log.Printf("Error sending discovery message: %v", err)
	}
}

// updatePeerInfo updates peer information from discovery message
func (nm *NetworkManager) updatePeerInfo(msg DiscoveryMessage, srcIP string) {
	nm.peersMutex.Lock()
	defer nm.peersMutex.Unlock()

	peer := &PeerInfo{
		PeerID:   msg.PeerID,
		Name:     msg.Name,
		IP:       srcIP,
		Port:     msg.Port,
		LastSeen: time.Now(),
		IsOnline: msg.Status == "online",
	}

	nm.activePeers[msg.PeerID] = peer

	// Emit event to frontend
	if nm.ctx != nil {
		runtime.EventsEmit(nm.ctx, "peerDiscovered", peer)
	}

	log.Printf("Peer discovered: %s (%s) at %s:%d", msg.Name, msg.PeerID, srcIP, msg.Port)
}

// peerCleanupRoutine removes inactive peers
func (nm *NetworkManager) peerCleanupRoutine() {
	ticker := time.NewTicker(60 * time.Second) // Check every minute
	defer ticker.Stop()

	for {
		select {
		case <-nm.stopChan:
			return
		case <-ticker.C:
			nm.cleanupInactivePeers()
		}
	}
}

// cleanupInactivePeers removes peers that haven't been seen for 5 minutes
func (nm *NetworkManager) cleanupInactivePeers() {
	nm.peersMutex.Lock()
	defer nm.peersMutex.Unlock()

	now := time.Now()
	timeout := 5 * time.Minute

	for peerID, peer := range nm.activePeers {
		if now.Sub(peer.LastSeen) > timeout {
			delete(nm.activePeers, peerID)
			log.Printf("Peer %s (%s) removed due to inactivity", peer.Name, peerID)

			// Emit offline event
			if nm.ctx != nil {
				runtime.EventsEmit(nm.ctx, "peerOffline", peerID)
			}
		}
	}
}
