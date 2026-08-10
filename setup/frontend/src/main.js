// Win-Slate Setup — audit, update check, install.

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
  busyKind: '', // 'audit' | 'install' | ''
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
  const onProgress = (p) => {
    if (!p || typeof p !== 'object') return
    state.progress = {
      step: p.step || state.progress.step || '',
      detail: p.detail || state.progress.detail || '',
      percent: typeof p.percent === 'number' ? p.percent : state.progress.percent || 0
    }
    // Update busy widgets in place — full re-render during await can stall WebView updates.
    const stepEl = document.getElementById('busy-step')
    const detailEl = document.getElementById('busy-detail')
    const barEl = document.getElementById('busy-bar')
    const pctEl = document.getElementById('busy-pct')
    if (stepEl && detailEl && barEl && pctEl) {
      stepEl.textContent = state.progress.step || 'Working'
      detailEl.textContent = state.progress.detail || 'Please wait…'
      const pct = Math.max(0, Math.min(100, state.progress.percent || 0))
      barEl.style.width = pct + '%'
      pctEl.textContent = pct + '%'
      return
    }
    render()
  }
  eventsOn('install:progress', onProgress)
  eventsOn('audit:progress', onProgress)
  render()
}

function busyPanel(title) {
  const pct = Math.max(0, Math.min(100, state.progress.percent || 0))
  const step = state.progress.step || 'Working'
  const detail = state.progress.detail || 'Please wait…'
  const wrap = el('div', { className: 'busy-panel', id: 'busy-panel' })
  wrap.appendChild(el('div', { className: 'busy-title', text: title || 'Checking this PC…' }))
  const row = el('div', { className: 'busy-row' })
  row.appendChild(el('div', { className: 'spinner' }))
  row.appendChild(
    el('div', { className: 'busy-text' }, [
      el('div', { className: 'busy-step', id: 'busy-step', text: step }),
      el('div', { className: 'busy-detail', id: 'busy-detail', text: detail })
    ])
  )
  wrap.appendChild(row)
  const bar = el('div', { className: 'progress' })
  const fill = el('div', { className: 'progress-bar', id: 'busy-bar' })
  fill.style.width = pct + '%'
  bar.appendChild(fill)
  wrap.appendChild(bar)
  wrap.appendChild(el('div', { className: 'progress-meta', id: 'busy-pct', text: pct + '%' }))
  return wrap
}

function render() {
  const root = document.getElementById('app')
  root.innerHTML = ''

  root.appendChild(
    el('div', { className: 'header' }, [
      el('h1', { text: '◆  Win-Slate Setup' }),
      el('p', {
        text: 'Install Win-Slate — a standalone Windows binary of Slate. Your Documents\\Slate projects are never modified.'
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
        ' (Apache-2.0). Win-Slate is an unofficial Windows port — not an official Wasserman release. ',
        el('a', {
          href: '#',
          onClick: (e) => {
            e.preventDefault()
            go()?.OpenExternal(state.paths.repoURL || 'https://github.com/ChrisToast89/Win-Slate')
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
        ? `Found Win-Slate: ${state.installStatus.installDir} (${state.installStatus.version || 'version unknown'})`
        : 'No Win-Slate install yet. Default folder: Programs\\Win-Slate.'
    })
  )
  const cards = el('div', { className: 'cards' })
  cards.appendChild(
    card(installed ? 'Update / reinstall Win-Slate' : 'Install Win-Slate', 'Check this PC, pick a folder, install the standalone app.', () => {
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

async function runAuditCore() {
  state.busy = true
  state.busyKind = 'audit'
  state.error = null
  state.progress = { step: 'Starting', detail: 'Starting system checks…', percent: 1 }
  render()

  // Client-side watchdog so the UI never sits silent forever if Go blocks.
  let tick = 1
  const watchdog = setInterval(() => {
    if (!state.busy || state.busyKind !== 'audit') return
    tick += 1
    if (tick > 45) {
      // 45s hard ceiling message; Go side should have timed out earlier
      const detailEl = document.getElementById('busy-detail')
      if (detailEl) {
        detailEl.textContent =
          'Still working… if this never finishes, close Setup and check %TEMP%\\win-slate-setup.log'
      }
    }
  }, 1000)

  try {
    // Yield a frame so the busy panel paints before the Go call.
    await new Promise((r) => requestAnimationFrame(() => r()))
    state.audit = await go().RunAudit()
    state.progress = { step: 'Updates', detail: 'Optional GitHub update check…', percent: 96 }
    const stepEl = document.getElementById('busy-step')
    const detailEl = document.getElementById('busy-detail')
    const barEl = document.getElementById('busy-bar')
    const pctEl = document.getElementById('busy-pct')
    if (stepEl) stepEl.textContent = state.progress.step
    if (detailEl) detailEl.textContent = state.progress.detail
    if (barEl) barEl.style.width = '96%'
    if (pctEl) pctEl.textContent = '96%'
    try {
      state.updates = await go().CheckForUpdates()
    } catch (ue) {
      state.updates = { message: 'Update check skipped — install still works offline.', ok: false }
    }
    if (state.audit?.installPath) state.installDir = state.audit.installPath
  } catch (e) {
    state.error = String(e?.message || e)
  } finally {
    clearInterval(watchdog)
    state.busy = false
    state.busyKind = ''
    render()
  }
}

async function runAudit() {
  await runAuditCore()
}

async function runAuditOnly() {
  await runAuditCore()
}

function renderChecks(parent) {
  if (state.busy && state.busyKind === 'audit') {
    parent.appendChild(busyPanel('Checking this PC…'))
    return
  }
  if (!state.audit) {
    parent.appendChild(el('p', { className: 'muted', text: 'No audit yet.' }))
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
  box.appendChild(el('h2', { text: state.mode === 'update' ? 'Update Win-Slate' : 'Install Win-Slate' }))
  box.appendChild(
    el('p', {
      className: 'muted',
      text: 'Choose the folder for Win-Slate program files. Projects stay in Documents\\Slate and are never overwritten.'
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
