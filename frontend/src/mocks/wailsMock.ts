// Mock Wails functions for development
export const Greet = async (name: string): Promise<string> => {
  return `Hello ${name}! (Development Mode)`;
};

export const SendMessage = async (peerId: string, message: string): Promise<void> => {
  console.log(`Sending message to ${peerId}: ${message}`);
};

export const BroadcastMessage = async (message: string): Promise<void> => {
  console.log(`Broadcasting message: ${message}`);
};

export const GetActivePeers = async (): Promise<Record<string, any>> => {
  return {
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
  };
};

export const GetLocalPeerInfo = async (): Promise<Record<string, string>> => {
  return {
    peer_id: "local_peer_789",
    name: "You (Dev Mode)"
  };
};

export const GetMessageHistory = async (peerId: string, limit: number): Promise<any[]> => {
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

  return allMessages.slice(0, limit);
};