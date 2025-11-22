import { useState, useEffect } from 'react'
import './App.css'
import './components/animations.css'
import PeerList from './components/PeerList'
import ChatWindow from './components/ChatWindow'
import NotificationPopup from './components/NotificationPopup'

// Declare Wails runtime for TypeScript
declare global {
  interface Window {
    runtime?: any;
  }

  interface ImportMeta {
    env: {
      DEV: boolean;
    };
  }
}

interface Peer {
  peer_id: string
  peer_name: string
  ip: string
  port: number
  last_seen: string
  is_online: boolean
}

interface Notification {
  id: string
  sender: string
  content: string
  timestamp: Date
}

// Mock functions for development
const mockGreet = async (name: string): Promise<string> => {
  return `Hello ${name}! (Development Mode)`;
};

const mockSendMessage = async (peerId: string, message: string): Promise<void> => {
  console.log(`Sending message to ${peerId}: ${message}`);
};

const mockBroadcastMessage = async (message: string): Promise<void> => {
  console.log(`Broadcasting message: ${message}`);
};

const mockGetActivePeers = async (): Promise<Record<string, any>> => {
  console.log('mockGetActivePeers called')
  const result = {
    "peer_123": {
      peer_id: "peer_123",
      peer_name: "Alice",
      ip: "192.168.1.100",
      port: 8080,
      last_seen: new Date().toISOString(),
      is_online: true
    },
    "peer_456": {
      peer_id: "peer_456",
      peer_name: "Bob",
      ip: "192.168.1.101",
      port: 8080,
      last_seen: new Date().toISOString(),
      is_online: true
    }
  }
  console.log('mockGetActivePeers returning:', result)
  return result
};

const mockGetLocalPeerInfo = async (): Promise<Record<string, string>> => {
  console.log('mockGetLocalPeerInfo called')
  const result = {
    peer_id: "local_peer_789",
    name: "You (Dev Mode)"
  }
  console.log('mockGetLocalPeerInfo returning:', result)
  return result
};

const mockGetMessageHistory = async (peerId: string, limit: number): Promise<any[]> => {
  console.log('mockGetMessageHistory called with peerId:', peerId, 'limit:', limit)
  const allMessages = [
    {
      id: 1,
      peer_id: peerId || "broadcast",
      sender_id: "peer_123",
      content: "Hello from Alice!",
      timestamp: new Date(Date.now() - 300000).toISOString(),
      is_read: true
    },
    {
      id: 2,
      peer_id: peerId || "broadcast",
      sender_id: "peer_456",
      content: "Hey everyone!",
      timestamp: new Date(Date.now() - 120000).toISOString(),
      is_read: true
    },
    {
      id: 3,
      peer_id: peerId || "broadcast",
      sender_id: "local_peer_789",
      content: "Hi there!",
      timestamp: new Date(Date.now() - 60000).toISOString(),
      is_read: true
    }
  ];

  const result = allMessages.slice(0, limit)
  console.log('mockGetMessageHistory returning:', result)
  return result
};

// Define functions - always use mocks for now to ensure they work
const Greet = mockGreet
const SendMessage = mockSendMessage
const BroadcastMessage = mockBroadcastMessage
const GetActivePeers = mockGetActivePeers
const GetLocalPeerInfo = mockGetLocalPeerInfo
const GetMessageHistory = mockGetMessageHistory

function App() {
  console.log('App component is mounting!')

  const [name, setName] = useState('')
  const [greeting, setGreeting] = useState('')
  const [selectedPeer, setSelectedPeer] = useState('')
  const [peers, setPeers] = useState<Record<string, Peer>>({})
  const [messages, setMessages] = useState<any[]>([])
  const [localPeer, setLocalPeer] = useState({ peer_id: '', name: '' })
  const [notifications, setNotifications] = useState<Notification[]>([])

  // Debug log to check environment
  console.log('Using mock functions for development')

  useEffect(() => {
    console.log('useEffect running')

    // Load message history
    console.log('Calling GetMessageHistory...')
    GetMessageHistory("", 100).then((messages: any[]) => {
      console.log('Messages loaded:', messages)
      setMessages(messages)
    }).catch((error) => {
      console.error('Error loading messages:', error)
    })

    // Load local peer info
    console.log('Calling GetLocalPeerInfo...')
    GetLocalPeerInfo().then((info: Record<string, string>) => {
      console.log('Local peer info:', info)
      setLocalPeer({
        peer_id: info.peer_id || '',
        name: info.name || ''
      })
    }).catch((error) => {
      console.error('Error loading peer info:', error)
    })

    // Load initial peers
    console.log('Calling GetActivePeers...')
    GetActivePeers().then((peers) => {
      console.log('Peers loaded:', peers)
      setPeers(peers)
    }).catch((error) => {
      console.error('Error loading peers:', error)
    })

    // Refresh peers periodically
    const interval = setInterval(() => {
      console.log('Refreshing peers...')
      GetActivePeers().then((peers) => {
        console.log('Peers refreshed:', peers)
        setPeers(peers)
      }).catch((error) => {
        console.error('Error refreshing peers:', error)
      })
    }, 5000)

    return () => clearInterval(interval)
  }, [])

  const handleGreet = async () => {
    if (!Greet) return
    try {
      const result = await Greet(name)
      setGreeting(result)
    } catch (error) {
      console.error('Error calling Greet:', error)
    }
  }

  const handleSendMessage = async (content: string) => {
    if (!SendMessage || !content.trim() || !selectedPeer) return

    try {
      await SendMessage(selectedPeer, content)
    } catch (error) {
      console.error('Error sending message:', error)
    }
  }

  const handleBroadcastMessage = async (content: string) => {
    if (!BroadcastMessage || !content.trim()) return

    try {
      await BroadcastMessage(content)
    } catch (error) {
      console.error('Error broadcasting message:', error)
    }
  }

  const handleDismissNotification = (id: string) => {
    setNotifications(prev => prev.filter(n => n.id !== id))
  }

  return (
    <div className="min-h-screen bg-gray-900 text-white p-6">
      <div className="max-w-6xl mx-auto">
        <header className="mb-8">
          <h1 className="text-3xl font-bold text-center mb-2">LanvoChat</h1>
          <p className="text-gray-400 text-center">P2P LAN Chat Application</p>
          <p className="text-blue-400 text-center mt-2">
            You: {localPeer.name} ({localPeer.peer_id})
          </p>
        </header>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* Peer List */}
          <div className="lg:col-span-1">
            <PeerList
              peers={peers}
              selectedPeer={selectedPeer}
              onPeerSelect={setSelectedPeer}
            />
          </div>

          {/* Chat Window */}
          <div className="lg:col-span-2">
            <ChatWindow
              messages={messages}
              selectedPeer={selectedPeer}
              localPeerId={localPeer.peer_id}
              onSendMessage={handleSendMessage}
              onBroadcastMessage={handleBroadcastMessage}
            />
          </div>
        </div>

        {/* Debug Section */}
        <div className="mt-8 p-4 bg-gray-800 rounded-lg">
          <h3 className="text-lg font-semibold mb-4">Debug</h3>
          <div className="flex space-x-2">
            <input
              type="text"
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder="Enter your name"
              onKeyPress={(e) => e.key === 'Enter' && handleGreet()}
              className="flex-1 px-3 py-2 bg-gray-700 border border-gray-600 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:border-blue-500"
            />
            <button
              onClick={handleGreet}
              className="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-colors duration-200"
            >
              Greet
            </button>
          </div>
          {greeting && <p className="mt-2 text-green-400">{greeting}</p>}
        </div>

        {/* Notification Popups */}
        <NotificationPopup
          notifications={notifications}
          onDismiss={handleDismissNotification}
        />
      </div>
    </div>
  )
}

export default App