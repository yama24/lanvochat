import React, { useEffect, useState } from 'react'

interface Notification {
  id: string
  sender: string
  content: string
  timestamp: Date
}

interface NotificationPopupProps {
  notifications: Notification[]
  onDismiss: (id: string) => void
}

const NotificationPopup: React.FC<NotificationPopupProps> = ({
  notifications,
  onDismiss
}) => {
  const [visibleNotifications, setVisibleNotifications] = useState<Notification[]>([])

  useEffect(() => {
    setVisibleNotifications(notifications)
  }, [notifications])

  const handleDismiss = (id: string) => {
    setVisibleNotifications(prev => prev.filter(n => n.id !== id))
    onDismiss(id)
  }

  if (visibleNotifications.length === 0) return null

  return (
    <div className="fixed bottom-4 right-4 z-50 space-y-2">
      {visibleNotifications.map((notification) => (
        <div
          key={notification.id}
          className="bg-gray-900 border border-gray-700 rounded-lg p-4 shadow-lg max-w-sm animate-slide-in"
          style={{
            animation: 'slideIn 0.3s ease-out'
          }}
        >
          <div className="flex justify-between items-start mb-2">
            <div className="text-sm font-medium text-blue-400">
              New message from {notification.sender}
            </div>
            <button
              onClick={() => handleDismiss(notification.id)}
              className="text-gray-400 hover:text-white ml-2"
            >
              ×
            </button>
          </div>
          <div className="text-white text-sm mb-2">
            {notification.content.length > 100
              ? `${notification.content.substring(0, 100)}...`
              : notification.content
            }
          </div>
          <div className="text-xs text-gray-400">
            {notification.timestamp.toLocaleTimeString()}
          </div>
        </div>
      ))}
    </div>
  )
}

export default NotificationPopup