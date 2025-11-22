import { useState } from 'react'
import './App.css'
import { Greet } from '../wailsjs/go/main/App'

function App() {
  const [name, setName] = useState('')
  const [greeting, setGreeting] = useState('')

  const handleGreet = async () => {
    try {
      const result = await Greet(name)
      setGreeting(result)
    } catch (error) {
      console.error('Error calling Greet:', error)
    }
  }

  return (
    <div className="App">
      <header className="App-header">
        <h1>LanvoChat</h1>
        <p>P2P LAN Chat Application</p>
        
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
      </header>
    </div>
  )
}

export default App
