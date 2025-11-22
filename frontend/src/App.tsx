import { useState, useEffect } from 'react'
import './App.css'
import { Greet, SendMessage, BroadcastMessage, GetActivePeers, GetLocalPeerInfo } from '../wailsjs/go/main/App'

interface Peer {
  peer_id: string
  name: string
  ip: string
  port: number
  last_seen: string
  is_online: boolean
}

interface Message {
  id: number
  peer_id: string
  sender_id: string
  content: string
  timestamp: string
  is_read: boolean
}

function App() {
  const [name, setName] = useState('')
  const [greeting, setGreeting] = useState('')
  const [message, setMessage] = useState('')
  const [selectedPeer, setSelectedPeer] = useState('')
  const [peers, setPeers] = useState<Record<string, Peer>>({})
  const [messages, setMessages] = useState<Message[]>([])
  const [localPeer, setLocalPeer] = useState({ peer_id: '', name: '' })

  useEffect(() => {
    // Get local peer info
    GetLocalPeerInfo().then((info: Record<string, string>) => {
      setLocalPeer({
        peer_id: info.peer_id || '',
        name: info.name || ''
      })
    })

    // Listen for network events
    const wailsRuntime = (window as any).runtime
    if (wailsRuntime) {
      wailsRuntime.EventsOn('peerDiscovered', (peer: Peer) => {
        setPeers(prev => ({ ...prev, [peer.peer_id]: peer }))
      })

      wailsRuntime.EventsOn('peerOffline', (peerID: string) => {
        setPeers(prev => {
          const updated = { ...prev }
          delete updated[peerID]
          return updated
        })
      })

      wailsRuntime.EventsOn('messageReceived', (msg: Message) => {
        setMessages(prev => [...prev, msg])
      })
    }

    // Refresh peers periodically
    const interval = setInterval(() => {
      GetActivePeers().then(setPeers)
    }, 5000)

    return () => clearInterval(interval)
  }, [])

  const handleGreet = async () => {
    try {
      const result = await Greet(name)
      setGreeting(result)
    } catch (error) {
      console.error('Error calling Greet:', error)
    }
  }

  const handleSendMessage = async () => {
    if (!message.trim() || !selectedPeer) return

    try {
      await SendMessage(selectedPeer, message)
      setMessage('')
    } catch (error) {
      console.error('Error sending message:', error)
    }
  }

  const handleBroadcastMessage = async () => {
    if (!message.trim()) return

    try {
      await BroadcastMessage(message)
      setMessage('')
    } catch (error) {
      console.error('Error broadcasting message:', error)
    }
  }

  return (
    <div className="App">
      <header className="App-header">
        <h1>LanvoChat</h1>
        <p>P2P LAN Chat Application</p>
        <p className="peer-info">You: {localPeer.name} ({localPeer.peer_id})</p>

        <div className="main-content">
          <div className="peers-section">
            <h3>Active Peers ({Object.keys(peers).length})</h3>
            <div className="peers-list">
              {Object.values(peers).map(peer => (
                <div
                  key={peer.peer_id}
                  className={`peer-item ${selectedPeer === peer.peer_id ? 'selected' : ''}`}
                  onClick={() => setSelectedPeer(peer.peer_id)}
                >
                  <div className="peer-name">{peer.name}</div>
                  <div className="peer-details">{peer.ip}:{peer.port}</div>
                  <div className={`peer-status ${peer.is_online ? 'online' : 'offline'}`}>
                    {peer.is_online ? '● Online' : '● Offline'}
                  </div>
                </div>
              ))}
            </div>
          </div>

          <div className="chat-section">
            <h3>Chat</h3>
            <div className="messages">
              {messages.map((msg, index) => (
                <div key={index} className="message">
                  <strong>{msg.sender_id}:</strong> {msg.content}
                  <small>({new Date(msg.timestamp).toLocaleTimeString()})</small>
                </div>
              ))}
            </div>

            <div className="message-input">
              <input
                type="text"
                value={message}
                onChange={(e) => setMessage(e.target.value)}
                placeholder="Type a message..."
                onKeyPress={(e) => e.key === 'Enter' && handleSendMessage()}
              />
              <button onClick={handleSendMessage} disabled={!selectedPeer}>
                Send to {selectedPeer ? peers[selectedPeer]?.name : 'Peer'}
              </button>
              <button onClick={handleBroadcastMessage}>
                Broadcast to All
              </button>
            </div>
          </div>
        </div>

        <div className="debug-section">
          <div className="input-box">
            <input
              type="text"
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder="Enter your name"
              onKeyPress={(e) => e.key === 'Enter' && handleGreet()}
            />
            <button onClick={handleGreet}>Greet</button>
          </div>

          {greeting && <p className="greeting">{greeting}</p>}
        </div>
      </header>
    </div>
  )
}

export default App
