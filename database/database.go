package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Database represents the SQLite database connection
type Database struct {
	db *sql.DB
}

// Message represents a chat message
type Message struct {
	ID        int64     `json:"id"`
	PeerID    string    `json:"peer_id"`
	SenderID  string    `json:"sender_id"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
	IsRead    bool      `json:"is_read"`
}

// Peer represents a chat peer/contact
type Peer struct {
	ID        int64     `json:"id"`
	PeerID    string    `json:"peer_id"`
	Name      string    `json:"name"`
	IPAddress string    `json:"ip_address"`
	LastSeen  time.Time `json:"last_seen"`
	IsOnline  bool      `json:"is_online"`
	CreatedAt time.Time `json:"created_at"`
}

// NewDatabase creates a new database connection and initializes the schema
func NewDatabase(dbPath string) (*Database, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	database := &Database{db: db}

	// Initialize schema
	if err := database.initSchema(); err != nil {
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return database, nil
}

// initSchema creates the necessary tables
func (d *Database) initSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS messages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		peer_id TEXT NOT NULL,
		sender_id TEXT NOT NULL,
		content TEXT NOT NULL,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
		is_read BOOLEAN DEFAULT 0,
		FOREIGN KEY (peer_id) REFERENCES peers(peer_id)
	);

	CREATE INDEX IF NOT EXISTS idx_messages_peer_id ON messages(peer_id);
	CREATE INDEX IF NOT EXISTS idx_messages_timestamp ON messages(timestamp);

	CREATE TABLE IF NOT EXISTS peers (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		peer_id TEXT UNIQUE NOT NULL,
		name TEXT NOT NULL,
		ip_address TEXT NOT NULL,
		last_seen DATETIME DEFAULT CURRENT_TIMESTAMP,
		is_online BOOLEAN DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_peers_peer_id ON peers(peer_id);
	CREATE INDEX IF NOT EXISTS idx_peers_is_online ON peers(is_online);
	`

	_, err := d.db.Exec(schema)
	return err
}

// SaveMessage saves a new message to the database
func (d *Database) SaveMessage(peerID, senderID, content string) error {
	query := `
		INSERT INTO messages (peer_id, sender_id, content, timestamp)
		VALUES (?, ?, ?, ?)
	`

	_, err := d.db.Exec(query, peerID, senderID, content, time.Now())
	if err != nil {
		return fmt.Errorf("failed to save message: %w", err)
	}

	return nil
}

// GetMessageHistory retrieves message history for a specific peer
func (d *Database) GetMessageHistory(peerID string, limit int) ([]Message, error) {
	if limit <= 0 {
		limit = 50 // default limit
	}

	query := `
		SELECT id, peer_id, sender_id, content, timestamp, is_read
		FROM messages
		WHERE peer_id = ?
		ORDER BY timestamp DESC
		LIMIT ?
	`

	rows, err := d.db.Query(query, peerID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query messages: %w", err)
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		err := rows.Scan(&msg.ID, &msg.PeerID, &msg.SenderID, &msg.Content, &msg.Timestamp, &msg.IsRead)
		if err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}
		messages = append(messages, msg)
	}

	// Reverse to get chronological order
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}

// SavePeer saves or updates a peer in the database
func (d *Database) SavePeer(peerID, name, ipAddress string) error {
	query := `
		INSERT INTO peers (peer_id, name, ip_address, last_seen, is_online)
		VALUES (?, ?, ?, ?, 1)
		ON CONFLICT(peer_id) DO UPDATE SET
			name = excluded.name,
			ip_address = excluded.ip_address,
			last_seen = excluded.last_seen,
			is_online = excluded.is_online
	`

	_, err := d.db.Exec(query, peerID, name, ipAddress, time.Now())
	if err != nil {
		return fmt.Errorf("failed to save peer: %w", err)
	}

	return nil
}

// GetPeers retrieves all peers from the database
func (d *Database) GetPeers() ([]Peer, error) {
	query := `
		SELECT id, peer_id, name, ip_address, last_seen, is_online, created_at
		FROM peers
		ORDER BY last_seen DESC
	`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query peers: %w", err)
	}
	defer rows.Close()

	var peers []Peer
	for rows.Next() {
		var peer Peer
		err := rows.Scan(&peer.ID, &peer.PeerID, &peer.Name, &peer.IPAddress,
			&peer.LastSeen, &peer.IsOnline, &peer.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan peer: %w", err)
		}
		peers = append(peers, peer)
	}

	return peers, nil
}

// UpdatePeerStatus updates the online status of a peer
func (d *Database) UpdatePeerStatus(peerID string, isOnline bool) error {
	query := `
		UPDATE peers
		SET is_online = ?, last_seen = ?
		WHERE peer_id = ?
	`

	_, err := d.db.Exec(query, isOnline, time.Now(), peerID)
	if err != nil {
		return fmt.Errorf("failed to update peer status: %w", err)
	}

	return nil
}

// MarkMessageAsRead marks a message as read
func (d *Database) MarkMessageAsRead(messageID int64) error {
	query := `UPDATE messages SET is_read = 1 WHERE id = ?`

	_, err := d.db.Exec(query, messageID)
	if err != nil {
		return fmt.Errorf("failed to mark message as read: %w", err)
	}

	return nil
}

// Close closes the database connection
func (d *Database) Close() error {
	return d.db.Close()
}
