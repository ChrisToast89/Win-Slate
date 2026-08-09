// Slate for Windows Setup — audit, update check, install to chosen folder.

const STEPS = ['Home', 'Check PC', 'Install', 'Finish']

function el(tag, attrs = {}, children = []) {
  const n = document.createElement(tag)
  for (const [k, v] of Object.entries(attrs)) {
    if (k === 'className') n.className = v
    else if (k === 'text') n.textContent = v
    else if (k.startsWith('on') && typeof v === 'function') n.addEventListener(k.slice(2).toLowerCase(), v)
    else if (k === 'checked') n.checked = !!v
    else if (k === 'value') n.value = v
    else if (v !== false && v != null) n.setAttribute(k, v)
  }
  for (const c of [].concat(children)) {
    if (c == null) continue
    n.appendChild(typeof c === 'string' ? document.createTextNode(c) : c)
  }
  return n
}

const state = {
  view: 'home', // home | wizard | finish | audit-only
  step: 0,
  paths: {},
  installStatus: null,
  audit: null,
  updates: null,
  installDir: '',
  desktop: true,
  busy: false,
  progress: { step: '', detail: '', percent: 0 },
  result: null,
  claudeMsg: '',
  error: null,
  mode: 'install' // install | update
}

function go() {
  return window.go?.main?.App
}

function eventsOn(name, cb) {
  if (window.runtime?.EventsOn) return window.runtime.EventsOn(name, cb)
  return () => {}
}

async function bootstrap() {
  try {
    state.paths = (await go().GetPaths()) || {}
    state.installDir = state.paths.installDir || state.paths.defaultInstall || ''
    state.installStatus = await go().GetInstallStatus()
  } catch (e) {
    console.warn(e)
  }
  eventsOn('install:progress', (p) => {
    state.progress = p || state.progress
    render()
  })
  render()
}

function render() {
  const root = document.getElementById('app')
  root.innerHTML = ''

  root.appendChild(
    el('div', { className: 'header' }, [
      el('h1', { text: '◆  Slate for Windows Setup' }),
      el('p', {
        text: 'Install or update the Windows build of Slate. Your Documents\\Slate projects are never modified.'
      }),
      el('p', { className: 'credit' }, [
        'Slate by ',
        el('a', {
          href: '#',
          onClick: (e) => {
            e.preventDefault()
            go()?.OpenExternal(state.paths.upstreamURL || 'https://github.com/wassermanproductions/slate')
          }
        }, [state.paths.upstreamAuthor || 'Sam Wasserman']),
        ' (Apache-2.0). This Setup packages a Windows port — not an official Wasserman release. ',
        el('a', {
          href: '#',
          onClick: (e) => {
            e.preventDefault()
            go()?.OpenExternal(state.paths.repoURL || 'https://github.com/ChrisToast89/slate-for-windows')
          }
        }, ['Repo'])
      ])
    ])
  )

  if (state.view === 'wizard' || state.view === 'finish') {
    const pills = el('div', { className: 'steps' })
    STEPS.forEach((name, i) => {
      const cls =
        'step-pill' + (i === state.step ? ' on' : '') + (i < state.step || state.view === 'finish' ? ' done' : '')
      pills.appendChild(el('div', { className: cls, text: `${i + 1}. ${name}` }))
    })
    root.appendChild(pills)
  }

  const main = el('div', { className: 'main' })
  if (state.view === 'home') main.appendChild(viewHome())
  else if (state.view === 'audit-only') main.appendChild(viewAuditOnly())
  else if (state.step === 1) main.appendChild(viewAudit())
  else if (state.step === 2) main.appendChild(viewInstall())
  else if (state.step === 3 || state.view === 'finish') main.appendChild(viewFinish())
  else main.appendChild(viewHome())
  root.appendChild(main)
}

function viewHome() {
  const installed = state.installStatus?.installed
  const box = el('div')
  box.appendChild(el('h2', { text: 'Welcome' }))
  box.appendChild(
    el('p', {
      className: 'muted',
      text: installed
        ? `Found install: ${state.installStatus.installDir} (${state.installStatus.version || 'version unknown'})`
        : 'No existing install detected. You can choose where to put Slate for Windows.'
    })
  )

  const cards = el('div', { className: 'cards' })
  cards.appendChild(
    card(installed ? 'Update / reinstall' : 'Install Slate for Windows', 'Check this PC, pick a folder, install the app.', () => {
      state.mode = installed ? 'update' : 'install'
      state.view = 'wizard'
      state.step = 1
      state.error = null
      state.audit = null
      void runAudit()
    })
  )
  cards.appendChild(
    card('Audit this PC', 'Dependencies, Claude Code, existing install — read-only.', () => {
      state.view = 'audit-only'
      state.error = null
      void runAuditOnly()
    })
  )
  cards.appendChild(
    card('Check for updates', 'Compare installed version to GitHub releases.', async () => {
      state.busy = true
      render()
      try {
        state.updates = await go().CheckForUpdates()
      } catch (e) {
        state.error = String(e)
      }
      state.busy = false
      render()
    })
  )
  box.appendChild(cards)

  if (state.updates) {
    box.appendChild(el('div', { className: 'summary', text: state.updates.message || JSON.stringify(state.updates, null, 2) }))
    if (state.updates.updateAvailable) {
      box.appendChild(
        el('div', { className: 'row' }, [
          el('button', {
            className: 'btn btn-key',
            text: 'Apply update…',
            onClick: () => {
              state.mode = 'update'
              state.view = 'wizard'
              state.step = 1
              void runAudit()
            }
          })
        ])
      )
    }
  }
  if (state.error) box.appendChild(el('div', { className: 'err', text: state.error }))
  return box
}

function card(title, body, onClick) {
  return el('div', { className: 'card', onClick }, [el('h3', { text: title }), el('p', { text: body })])
}

async function runAudit() {
  state.busy = true
  state.error = null
  render()
  try {
    state.audit = await go().RunAudit()
    state.updates = await go().CheckForUpdates()
    if (state.audit?.installPath) state.installDir = state.audit.installPath
  } catch (e) {
    state.error = String(e)
  }
  state.busy = false
  render()
}

async function runAuditOnly() {
  state.busy = true
  render()
  try {
    state.audit = await go().RunAudit()
    state.updates = await go().CheckForUpdates()
  } catch (e) {
    state.error = String(e)
  }
  state.busy = false
  render()
}

function renderChecks(parent) {
  if (!state.audit) {
    parent.appendChild(el('p', { className: 'muted', text: state.busy ? 'Checking…' : 'No audit yet.' }))
    return
  }
  parent.appendChild(el('p', { className: 'muted', text: state.audit.summary }))
  const ul = el('ul', { className: 'check-list' })
  for (const c of state.audit.checks || []) {
    ul.appendChild(
      el('li', {}, [
        el('div', { className: 'mark ' + (c.ok ? 'ok' : 'bad'), text: c.ok ? '✓' : '✗' }),
        el('div', {}, [
          el('div', { className: 'label', text: c.label + (c.required ? ' (required)' : '') }),
          el('div', { className: 'detail', text: c.detail || '' }),
          el('div', { className: 'action', text: c.action || '' })
        ])
      ])
    )
  }
  parent.appendChild(ul)
  if (state.updates?.message) {
    parent.appendChild(el('div', { className: 'summary', text: 'Repo check: ' + state.updates.message }))
  }
}

function viewAuditOnly() {
  const box = el('div')
  box.appendChild(el('h2', { text: 'System audit' }))
  renderChecks(box)
  if (state.error) box.appendChild(el('div', { className: 'err', text: state.error }))
  box.appendChild(
    el('div', { className: 'row' }, [
      el('button', {
        className: 'btn',
        text: 'Back',
        onClick: () => {
          state.view = 'home'
          render()
        }
      }),
      el('button', {
        className: 'btn',
        text: 'Re-run',
        disabled: state.busy,
        onClick: () => void runAuditOnly()
      })
    ])
  )
  return box
}

function viewAudit() {
  const box = el('div')
  box.appendChild(el('h2', { text: 'Check this PC' }))
  renderChecks(box)
  if (state.error) box.appendChild(el('div', { className: 'err', text: state.error }))
  box.appendChild(
    el('div', { className: 'row' }, [
      el('button', {
        className: 'btn',
        text: 'Back',
        onClick: () => {
          state.view = 'home'
          state.step = 0
          render()
        }
      }),
      el('button', {
        className: 'btn',
        text: 'Re-check',
        disabled: state.busy,
        onClick: () => void runAudit()
      }),
      el('button', {
        className: 'btn btn-key',
        text: 'Continue',
        disabled: state.busy || (state.audit && !state.audit.canProceed),
        onClick: () => {
          state.step = 2
          render()
        }
      })
    ])
  )
  return box
}

function viewInstall() {
  const box = el('div')
  box.appendChild(el('h2', { text: state.mode === 'update' ? 'Update Slate for Windows' : 'Install Slate for Windows' }))
  box.appendChild(
    el('p', {
      className: 'muted',
      text: 'Choose the folder for program files. Projects stay in Documents\\Slate and are never overwritten.'
    })
  )

  const pathRow = el('div', { className: 'path-box' })
  const input = el('input', {
    value: state.installDir,
    onInput: (e) => {
      state.installDir = e.target.value
    }
  })
  pathRow.appendChild(input)
  pathRow.appendChild(
    el('button', {
      className: 'btn',
      text: 'Browse…',
      disabled: state.busy,
      onClick: async () => {
        const p = await go().PickInstallFolder(state.installDir)
        if (p) {
          state.installDir = p
          render()
        }
      }
    })
  )
  box.appendChild(pathRow)

  box.appendChild(
    el('label', { className: 'chk' }, [
      el('input', {
        type: 'checkbox',
        checked: state.desktop,
        onChange: (e) => {
          state.desktop = e.target.checked
        }
      }),
      el('span', { text: 'Create Desktop shortcut' })
    ])
  )

  if (state.busy) {
    const bar = el('div', { className: 'progress' }, [
      el('div', { className: 'progress-bar', style: `width:${state.progress.percent || 0}%` })
    ])
    // style attr workaround
    bar.firstChild.style.width = `${state.progress.percent || 0}%`
    box.appendChild(bar)
    box.appendChild(
      el('div', {
        className: 'progress-meta',
        text: `${state.progress.percent || 0}% — ${state.progress.step || ''} ${state.progress.detail || ''}`
      })
    )
  }

  if (state.error) box.appendChild(el('div', { className: 'err', text: state.error }))

  box.appendChild(
    el('div', { className: 'row' }, [
      el('button', {
        className: 'btn',
        text: 'Back',
        disabled: state.busy,
        onClick: () => {
          state.step = 1
          render()
        }
      }),
      el('button', {
        className: 'btn btn-key',
        text: state.mode === 'update' ? 'Update now' : 'Install now',
        disabled: state.busy || !state.installDir,
        onClick: () => void doInstall()
      })
    ])
  )
  return box
}

async function doInstall() {
  state.busy = true
  state.error = null
  state.progress = { step: 'Starting', detail: '', percent: 1 }
  render()
  try {
    const out = await go().StartInstall(state.installDir, state.desktop, state.mode === 'update')
    state.result = out
    state.claudeMsg = out?.claudeMsg || ''
    state.step = 3
    state.view = 'finish'
  } catch (e) {
    state.error = String(e?.message || e)
  }
  state.busy = false
  render()
}

function viewFinish() {
  const box = el('div')
  box.appendChild(el('h2', { text: 'Done' }))
  const summary = state.result?.result?.summary || 'Install finished.'
  box.appendChild(el('div', { className: 'summary', text: summary }))
  if (state.claudeMsg) {
    box.appendChild(el('h2', { text: 'Claude Code (AI brain)', style: 'margin-top:16px' }))
    box.appendChild(el('div', { className: 'summary', text: state.claudeMsg }))
  }
  box.appendChild(
    el('p', {
      className: 'muted',
      text: 'Slate by Sam Wasserman. Sign in to Claude Code with `claude auth login` if you use the Claude brain — Setup never stores API keys.'
    })
  )
  box.appendChild(
    el('div', { className: 'row' }, [
      el('button', {
        className: 'btn btn-key',
        text: 'Launch Slate',
        onClick: async () => {
          try {
            await go().LaunchApp()
          } catch (e) {
            state.error = String(e)
            render()
          }
        }
      }),
      el('button', {
        className: 'btn',
        text: 'Install Claude Code',
        onClick: async () => {
          state.busy = true
          render()
          try {
            const r = await go().InstallClaude()
            state.claudeMsg = r?.message || ''
          } catch (e) {
            state.error = String(e)
          }
          state.busy = false
          render()
        }
      }),
      el('button', {
        className: 'btn',
        text: 'Claude login…',
        onClick: async () => {
          try {
            const msg = await go().LaunchClaudeLogin()
            state.claudeMsg = msg
            render()
          } catch (e) {
            state.error = String(e)
            render()
          }
        }
      }),
      el('button', {
        className: 'btn',
        text: 'Projects folder',
        onClick: () => go()?.OpenProjectsFolder()
      }),
      el('button', {
        className: 'btn btn-ghost',
        text: 'Home',
        onClick: () => {
          state.view = 'home'
          state.step = 0
          void bootstrap()
        }
      })
    ])
  )
  if (state.error) box.appendChild(el('div', { className: 'err', text: state.error }))
  return box
}

bootstrap()
