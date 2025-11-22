import React from 'react'

interface Peer {
  peer_id: string
  peer_name: string
  ip: string
  port: number
  last_seen: string
  is_online: boolean
}

interface PeerListProps {
  peers: Record<string, Peer>
  selectedPeer: string
  onPeerSelect: (peerId: string) => void
}

const PeerList: React.FC<PeerListProps> = ({ peers, selectedPeer, onPeerSelect }) => {
  const peerArray = Object.values(peers)

  return (
    <div className="bg-gray-800 rounded-lg p-4 h-96 overflow-y-auto">
      <h3 className="text-lg font-semibold mb-4 text-white">
        Active Peers ({peerArray.length})
      </h3>
      <div className="space-y-2">
        {peerArray.length === 0 ? (
          <div className="text-gray-400 text-center py-8">
            No peers discovered yet
          </div>
        ) : (
          peerArray.map(peer => (
            <div
              key={peer.peer_id}
              className={`p-3 rounded-lg cursor-pointer transition-all duration-200 ${
                selectedPeer === peer.peer_id
                  ? 'bg-blue-600 border border-blue-400'
                  : 'bg-gray-700 hover:bg-gray-600 border border-gray-600'
              }`}
              onClick={() => onPeerSelect(peer.peer_id)}
            >
              <div className="flex items-center justify-between">
                <div className="flex-1">
                  <div className="font-medium text-white">{peer.peer_name}</div>
                  <div className="text-sm text-gray-300">{peer.ip}:{peer.port}</div>
                </div>
                <div className={`w-3 h-3 rounded-full ${
                  peer.is_online ? 'bg-green-400' : 'bg-red-400'
                }`} />
              </div>
              <div className="text-xs text-gray-400 mt-1">
                Last seen: {new Date(peer.last_seen).toLocaleTimeString()}
              </div>
            </div>
          ))
        )}
      </div>
    </div>
  )
}

export default PeerList