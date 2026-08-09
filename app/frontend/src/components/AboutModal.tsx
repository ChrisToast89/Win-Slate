// About Slate — brand art, credit, and support links.
// Windows port retains full attribution required by Apache-2.0 NOTICE.

import React, { useEffect, useState } from 'react'
import brandArt from '../assets/brand.webp'

const FALLBACK_VERSION = '0.3.2-win.1'

export default function AboutModal({ onClose }: { onClose(): void }): React.JSX.Element {
  const [version, setVersion] = useState(FALLBACK_VERSION)

  useEffect(() => {
    const vfn = window.slate && 'appVersion' in window.slate ? window.slate.appVersion : null
    if (typeof vfn === 'function') {
      void vfn().then((v: string) => setVersion(v || FALLBACK_VERSION)).catch(() => {})
    }
  }, [])

  return (
    <div className="modal-scrim" onClick={onClose}>
      <div className="modal about-modal" onClick={(e) => e.stopPropagation()}>
        <div className="modal-body" style={{ textAlign: 'center', padding: '22px 26px' }}>
          <img className="about-brand" src={brandArt} alt="Slate" />
          <p className="about-text">
            Slate is the prompt studio for AI filmmaking — plan shots, direct coverage, spot your
            score, cast your voices, and keep continuity across an entire film while an AI brain
            helps you craft every prompt. Compile each one for your generator of choice; Slate
            makes the prompts, your generators make the picture and sound.
          </p>
          <p className="about-meta">
            Version {version} · Apache-2.0 · Created by Sam Wasserman
            <br />
            Windows port: personal derivative of Slate 0.3.2 (Go + Wails host).
            <br />
            Brain: your Claude Code or Codex sign-in, or a local model — no API keys.
            <br />
            <a href="https://wassermanproductions.com" target="_blank" rel="noreferrer">
              wassermanproductions.com
            </a>{' '}
            ·{' '}
            <a href="https://wasserman.ai" target="_blank" rel="noreferrer">
              wasserman.ai
            </a>
            <br />
            <a href="https://github.com/wassermanproductions/slate" target="_blank" rel="noreferrer">
              Original source on GitHub
            </a>
          </p>
          <div style={{ display: 'flex', gap: 8, justifyContent: 'center', flexWrap: 'wrap' }}>
            <a
              className="btn btn-key"
              href="https://ko-fi.com/samwasserman"
              target="_blank"
              rel="noreferrer"
              style={{ textDecoration: 'none' }}
            >
              ♥ Support on Ko-fi
            </a>
            <button className="btn" onClick={onClose}>
              Close
            </button>
          </div>
        </div>
      </div>
    </div>
  )
}
