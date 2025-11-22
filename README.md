# LanvoChat

P2P LAN Chat Desktop Application built with Wails (Go + React/TypeScript)

## Features

- 🔒 P2P encrypted messaging over LAN
- 💾 SQLite database for message history
- 🪟 Frameless window with system tray support
- 🌐 **UDP Multicast peer discovery**
- 📡 **TCP messaging between peers**
- 🔄 **Real-time message events**
- ⚡ Fast and lightweight

## Tech Stack

- **Backend**: Go 1.25.3
- **Frontend**: React + TypeScript
- **Framework**: Wails v2
- **Database**: SQLite3
- **Networking**: UDP Multicast + TCP
- **Package Manager**: Yarn (Node.js 22.21.0)

## Network Architecture

### UDP Multicast Discovery (Port 8081)
- **Purpose**: Auto-discover peers on LAN
- **Address**: `239.255.255.250:1900`
- **Broadcast**: Every 30 seconds
- **Cleanup**: Inactive peers removed after 5 minutes

### TCP Messaging (Port 8080)
- **Purpose**: Reliable message delivery
- **Connection**: Direct peer-to-peer
- **Timeout**: 30s read, 10s write
- **Format**: JSON with message metadata

## Prerequisites

- Go 1.25.3+
- Node.js 22.21.0+
- Yarn
- Wails CLI v2

## Getting Started

### Install Dependencies

```bash
# Install Go dependencies
go mod download

# Install frontend dependencies
cd frontend
yarn install
cd ..
```

### Development

```bash
# Run in development mode (Linux dengan WebKit2GTK 4.1)
./dev.sh
```

### Build

```bash
# Build for production (Linux dengan WebKit2GTK 4.1)
./build.sh

# Jalankan hasil build
./build/bin/lanvochat
```

## Project Structure

```
lanvochat/
├── main.go              # Application entry point
├── app.go               # App structure and API bindings
├── database/            # SQLite database layer
│   └── database.go
├── network/             # Network communication layer
│   ├── network.go       # Main network manager
│   ├── udp_multicast.go # UDP multicast discovery
│   └── tcp_handler.go   # TCP messaging
├── frontend/            # React frontend
│   ├── src/
│   └── package.json
├── wails.json           # Wails configuration
└── go.mod
```

## Database Schema

### Messages Table
- id (INTEGER PRIMARY KEY)
- peer_id (TEXT)
- sender_id (TEXT)
- content (TEXT)
- timestamp (DATETIME)
- is_read (BOOLEAN)

### Peers Table
- id (INTEGER PRIMARY KEY)
- peer_id (TEXT UNIQUE)
- name (TEXT)
- ip_address (TEXT)
- last_seen (DATETIME)
- is_online (BOOLEAN)
- created_at (DATETIME)

## API Methods

### Network Operations
- `SendMessage(peerID, content)` - Send to specific peer
- `BroadcastMessage(content)` - Send to all peers
- `GetActivePeers()` - Get discovered peers
- `GetLocalPeerInfo()` - Get local peer details

### Database Operations
- `SaveMessage(peerID, senderID, content)`
- `GetMessageHistory(peerID, limit)`
- `SavePeer(peerID, name, ipAddress)`
- `GetPeers()`

## License

MIT
