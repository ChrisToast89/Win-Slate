import React from 'react'
import ReactDOM from 'react-dom/client'
import App from './App'
import { installWailsHost } from './lib/host'
import { installDevMock } from './lib/devMock'
import './styles/global.css'

// Prefer live Go bindings (Wails). Fall back to in-browser mock for pure Vite UI work.
if (!installWailsHost()) {
  installDevMock()
}

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>
)
