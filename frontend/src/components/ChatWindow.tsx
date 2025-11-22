import React, { useState, useRef, useEffect } from 'react'

interface ChatWindowProps {
  messages: any[]
  selectedPeer: string
  localPeerId: string
  onSendMessage: (content: string) => void
  onBroadcastMessage: (content: string) => void
}

const ChatWindow: React.FC<ChatWindowProps> = ({
  messages,
  selectedPeer,
  localPeerId,
  onSendMessage,
  onBroadcastMessage
}) => {
  const [message, setMessage] = useState('')
  const messagesEndRef = useRef<HTMLDivElement>(null)

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' })
  }

  useEffect(() => {
    scrollToBottom()
  }, [messages])

  const handleSend = () => {
    if (!message.trim()) return

    if (selectedPeer) {
      onSendMessage(message)
    } else {
      onBroadcastMessage(message)
    }
    setMessage('')
  }

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault()
      handleSend()
    }
  }

  const filteredMessages = selectedPeer
    ? messages.filter(msg => msg.peer_id === selectedPeer || msg.sender_id === selectedPeer)
    : messages

  return (
    <div className="bg-gray-800 rounded-lg p-4 h-96 flex flex-col">
      <h3 className="text-lg font-semibold mb-4 text-white">
        Chat {selectedPeer ? `with Peer` : '(Broadcast)'}
      </h3>

      {/* Messages */}
      <div className="flex-1 overflow-y-auto mb-4 space-y-3">
        {filteredMessages.length === 0 ? (
          <div className="text-gray-400 text-center py-8">
            {selectedPeer ? 'No messages with this peer yet' : 'No messages yet'}
          </div>
        ) : (
          filteredMessages.map((msg) => (
            <div
              key={msg.id}
              className={`flex ${msg.sender_id === localPeerId ? 'justify-end' : 'justify-start'}`}
            >
              <div
                className={`max-w-xs lg:max-w-md px-4 py-2 rounded-lg ${
                  msg.sender_id === localPeerId
                    ? 'bg-blue-600 text-white'
                    : 'bg-gray-700 text-white'
                }`}
              >
                <div className="text-sm font-medium mb-1">
                  {msg.sender_id === localPeerId ? 'You' : msg.sender_id}
                </div>
                <div className="text-sm">{msg.content}</div>
                <div className="text-xs opacity-70 mt-1">
                  {new Date(msg.timestamp).toLocaleTimeString()}
                </div>
              </div>
            </div>
          ))
        )}
        <div ref={messagesEndRef} />
      </div>

      {/* Message Input */}
      <div className="flex space-x-2">
        <input
          type="text"
          value={message}
          onChange={(e) => setMessage(e.target.value)}
          onKeyPress={handleKeyPress}
          placeholder={selectedPeer ? `Message to peer...` : "Broadcast message..."}
          className="flex-1 px-3 py-2 bg-gray-700 border border-gray-600 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:border-blue-500"
        />
        <button
          onClick={handleSend}
          disabled={!message.trim()}
          className="px-4 py-2 bg-blue-600 hover:bg-blue-700 disabled:bg-gray-600 disabled:cursor-not-allowed text-white rounded-lg transition-colors duration-200"
        >
          {selectedPeer ? 'Send' : 'Broadcast'}
        </button>
      </div>
    </div>
  )
}

export default ChatWindow