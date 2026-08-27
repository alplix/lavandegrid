import './style.css'

const $ = s => document.querySelector(s)
const $$ = s => [...document.querySelectorAll(s)]

let state = {
  page: 'dashboard',
  hosts: [],
  snaps: {},
  hist: {},
  selectedHost: null,
  theme: localStorage.getItem('lavande-theme') || 'dark',
  toasts: [],
  searchQuery: '',
  modal: null,
  daemonStatus: 'unknown',
  daemonInfo: null,
  stats: {},
  xfers: {},
  disk: {}
}

function toast(msg, type = 'info') {
  const id = Date.now()
  state.toasts.push({ id, msg, type })
  renderToasts()
  setTimeout(() => {
    state.toasts = state.toasts.filter(t => t.id !== id)
    renderToasts()
  }, 3500)
}

function renderToasts() {
  const wrap = $('.toast-wrap')
  if (!wrap) return
  wrap.innerHTML = state.toasts.map(t => `<div class="toast ${t.type}">${esc(t.msg)}</div>`).join('')
}

function esc(s) { const d = document.createElement('div'); d.textContent = s; return d.innerHTML }

function fmtCredit(v) {
  if (v >= 1e9) return (v / 1e9).toFixed(2) + ' B'
  if (v >= 1e6) return (v / 1e6).toFixed(2) + ' M'
  if (v >= 1e4) return (v / 1e3).toFixed(1) + ' k'
  if (v >= 100) return Math.round(v).toString()
  return v.toFixed(1)
}

function fmtBytes(b) {
  if (b < 1024) return b + ' B'
  const units = ['KB', 'MB', 'GB', 'TB']
  let u = -1, n = b
  do { n /= 1024; u++ } while (n >= 1024 && u < units.length - 1)
  return n.toFixed(1) + ' ' + units[u]
}

function fmtDuration(s) {
  if (s <= 0) return '-'
  const sec = Math.floor(s)
  const d = Math.floor(sec / 86400)
  const h = Math.floor((sec % 86400) / 3600)
  const m = Math.floor((sec % 3600) / 60)
  if (d > 0) return `${d}d ${h}h`
  if (h > 0) return `${h}h ${String(m).padStart(2, '0')}m`
  if (m > 0) return `${m}m ${String(sec % 60).padStart(2, '0')}s`
  return `${sec}s`
}

function fmtAgo(ts) {
  if (!ts) return '-'
  const d = (Date.now() / 1000) - ts
  if (d < 60) return 'just now'
  if (d < 3600) return Math.floor(d / 60) + ' min ago'
  if (d < 86400) return Math.floor(d / 3600) + ' h ago'
  return Math.floor(d / 86400) + ' d ago'
}

function fmtTime(ts) {
  if (!ts) return '-'
  return new Date(ts * 1000).toLocaleString('en-US', { month: 'short', day: '2-digit', hour: '2-digit', minute: '2-digit' })
}

function projColor(url) {
  let h = 0
  for (let i = 0; i < url.length; i++) h = h * 31 + url.charCodeAt(i)
  return `hsl(${Math.abs(h) % 360}, 65%, 55%)`
}

function statusBadge(s) {
  const map = { running: 'running', paused: 'paused', queued: 'queued', downloading: 'download', uploading: 'upload', error: 'error', ready: 'ready' }
  return `<span class="badge ${map[s] || 'queued'}">${s}</span>`
}

async function api(method, ...args) {
  try {
    const fn = window.go.main.App[method]
    if (!fn) throw new Error(`Method ${method} not found`)
    return await fn(...args)
  } catch (e) {
    throw e
  }
}

async function refreshHosts() {
  state.hosts = await api('GetHosts')
  for (const h of state.hosts) {
    try {
      const snap = await api('GetSnapshot', h.id)
      if (snap) state.snaps[h.id] = snap
      const hist = await api('GetHistory', h.id)
      if (hist) state.hist[h.id] = hist
    } catch (e) {}
  }
  render()
}

function setPage(page) {
  state.page = page
  state.modal = null
  render()
}

function render() {
  renderShell()
  renderContent()
}

function renderShell() {
  const app = $('#app')
  const running = Object.values(state.snaps).reduce((a, s) => a + (s.online ? s.totals.running : 0), 0)
  app.innerHTML = `
    <div class="shell">
      <aside class="sidebar">
        <div class="side-logo">
          <div class="logo-icon">
            <svg width="22" height="22" viewBox="0 0 24 24" fill="currentColor" opacity="0.95"><ellipse cx="12" cy="5.5" rx="2.2" ry="4"/><ellipse cx="8.5" cy="7" rx="2" ry="3.5" transform="rotate(-25 8.5 7)"/><ellipse cx="15.5" cy="7" rx="2" ry="3.5" transform="rotate(25 15.5 7)"/><ellipse cx="7" cy="10.5" rx="1.8" ry="3" transform="rotate(-40 7 10.5)"/><ellipse cx="17" cy="10.5" rx="1.8" ry="3" transform="rotate(40 17 10.5)"/><line x1="12" y1="11" x2="12" y2="22" stroke="currentColor" stroke-width="1.8" fill="none" opacity="0.7"/><path d="M12 15c-2 1.5-4 1-5 0" stroke="currentColor" stroke-width="1.2" fill="none" opacity="0.5"/><path d="M12 17.5c1.8 1.5 3.5 1 4.5 0" stroke="currentColor" stroke-width="1.2" fill="none" opacity="0.5"/></svg>
          </div>
          <div class="logo-text">
            <div class="logo-title">LavandeGrid</div>
            <div class="logo-sub">Camellia Manager</div>
          </div>
        </div>
        <nav class="nav">
          <div class="nav-item ${state.page === 'dashboard' ? 'active' : ''}" onclick="window._setPage('dashboard')">
            <svg width="19" height="19" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="3" width="7" height="7" rx="1"/><rect x="14" y="3" width="7" height="7" rx="1"/><rect x="3" y="14" width="7" height="7" rx="1"/><rect x="14" y="14" width="7" height="7" rx="1"/></svg>
            <span>Dashboard</span>
          </div>
          <div class="nav-item ${state.page === 'tasks' ? 'active' : ''}" onclick="window._setPage('tasks')">
            <svg width="19" height="19" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M9 11l3 3L22 4"/><path d="M21 12v7a2 2 0 01-2 2H5a2 2 0 01-2-2V5a2 2 0 012-2h11"/></svg>
            <span>Tasks</span>
            ${running > 0 ? `<span class="nav-count num">${running}</span>` : ''}
          </div>
          <div class="nav-item ${state.page === 'projects' ? 'active' : ''}" onclick="window._setPage('projects')">
            <svg width="19" height="19" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M22 19a2 2 0 01-2 2H4a2 2 0 01-2-2V5a2 2 0 012-2h5l2 3h9a2 2 0 012 2z"/></svg>
            <span>Projects</span>
          </div>
          <div class="nav-item ${state.page === 'transfers' ? 'active' : ''}" onclick="window._setPage('transfers')">
            <svg width="19" height="19" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M7 16V4m0 0L3 8m4-4l4 4m6 0v12m0 0l4-4m-4 4l-4-4"/></svg>
            <span>Transfers</span>
          </div>
          <div class="nav-item ${state.page === 'messages' ? 'active' : ''}" onclick="window._setPage('messages')">
            <svg width="19" height="19" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z"/><path d="M14 2v6h6"/><path d="M16 13H8"/><path d="M16 17H8"/><path d="M10 9H8"/></svg>
            <span>Messages</span>
          </div>
          <div class="nav-item ${state.page === 'stats' ? 'active' : ''}" onclick="window._setPage('stats')">
            <svg width="19" height="19" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 20V10"/><path d="M18 20V4"/><path d="M6 20v-4"/></svg>
            <span>Stats</span>
          </div>
          <div class="nav-item ${state.page === 'hosts' ? 'active' : ''}" onclick="window._setPage('hosts')">
            <svg width="19" height="19" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="2" y="2" width="20" height="8" rx="2" ry="2"/><rect x="2" y="14" width="20" height="8" rx="2" ry="2"/><line x1="6" y1="6" x2="6.01" y2="6"/><line x1="6" y1="18" x2="6.01" y2="18"/></svg>
            <span>Servers</span>
          </div>
          <div class="nav-item ${state.page === 'settings' ? 'active' : ''}" onclick="window._setPage('settings')">
            <svg width="19" height="19" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="3"/><path d="M19.4 15a1.65 1.65 0 00.33 1.82l.06.06a2 2 0 010 2.83 2 2 0 01-2.83 0l-.06-.06a1.65 1.65 0 00-1.82-.33 1.65 1.65 0 00-1 1.51V21a2 2 0 01-4 0v-.09A1.65 1.65 0 009 19.4a1.65 1.65 0 00-1.82.33l-.06.06a2 2 0 01-2.83-2.83l.06-.06A1.65 1.65 0 004.68 15a1.65 1.65 0 00-1.51-1H3a2 2 0 010-4h.09A1.65 1.65 0 004.6 9a1.65 1.65 0 00-.33-1.82l-.06-.06a2 2 0 012.83-2.83l.06.06A1.65 1.65 0 009 4.68a1.65 1.65 0 001-1.51V3a2 2 0 014 0v.09a1.65 1.65 0 001 1.51 1.65 1.65 0 001.82-.33l.06-.06a2 2 0 012.83 2.83l-.06.06A1.65 1.65 0 0019.4 9a1.65 1.65 0 001.51 1H21a2 2 0 010 4h-.09a1.65 1.65 0 00-1.51 1z"/></svg>
            <span>Settings</span>
          </div>
        </nav>
        <div class="side-foot">
          <span>v1.0 - desktop management</span>
          <span class="credit">Coded by <b>Alperen Yavuz</b></span>
        </div>
      </aside>
      <div class="main">
        <div class="topbar">
          <div style="flex:1">
            <div class="page-title">${pageTitle()}</div>
            <div class="page-sub">${pageSub()}</div>
          </div>
          <button class="btn sm" onclick="window._toggleTheme()" title="Toggle theme">
            ${state.theme === 'dark' ? '☀️' : '🌙'}
          </button>
        </div>
        <div class="content" id="content"></div>
      </div>
    </div>
    <div class="toast-wrap"></div>
    ${state.modal ? state.modal : ''}
  `
  renderToasts()
}

function pageTitle() {
  const map = { dashboard: 'Dashboard', tasks: 'Tasks', projects: 'Projects', transfers: 'Transfers', messages: 'Messages', hosts: 'Servers', settings: 'Settings' }
  return map[state.page] || 'Dashboard'
}

function pageSub() {
  const online = Object.values(state.snaps).filter(s => s.online).length
  return `${state.hosts.length} host(s) configured, ${online} online`
}

function renderContent() {
  const c = $('#content')
  if (!c) return
  switch (state.page) {
    case 'dashboard': c.innerHTML = renderDashboard(); break
    case 'tasks': c.innerHTML = renderTasks(); break
    case 'projects': c.innerHTML = renderProjects(); break
    case 'transfers': c.innerHTML = renderTransfers(); break
    case 'messages': c.innerHTML = renderMessages(); break
    case 'stats': c.innerHTML = renderStats(); break
    case 'hosts': c.innerHTML = renderHosts(); break
    case 'settings': c.innerHTML = renderSettings(); break
    default: c.innerHTML = renderDashboard()
  }
}

function renderDashboard() {
  const list = Object.values(state.snaps).filter(s => s.online)
  const running = list.reduce((a, s) => a + s.totals.running, 0)
  const queued = list.reduce((a, s) => a + s.totals.queued + s.totals.paused, 0)
  const rac = list.reduce((a, s) => a + s.totals.rac, 0)
  const credit = list.reduce((a, s) => a + s.totals.credit, 0)

  let serverCards = ''
  for (const h of state.hosts) {
    const snap = state.snaps[h.id]
    const online = snap?.online
    const projects = (snap?.projects || []).slice(0, 4).map(p =>
      `<span class="chip"><span class="chip-dot" style="background:${projColor(p.url)}"></span><span class="trunc">${esc(p.name || p.url)}</span></span>`
    ).join('')

    serverCards += `
      <div class="card hoverable" style="display:flex;flex-direction:column;gap:12px">
        <div class="spread">
          <div class="row" style="min-width:0">
            <span class="dot ${online ? 'on' : 'off'}" style="flex-shrink:0"></span>
            <div style="min-width:0">
              <div class="trunc" style="font-weight:700;font-size:15px">${esc(h.name)}${h.demo ? '<span class="chip plain" style="margin-left:8px">demo</span>' : ''}</div>
              <div class="faint num">${esc(h.host)}:${h.port}</div>
            </div>
          </div>
          ${online && snap.version ? `<span class="chip indigo num">v${esc(snap.version)}</span>` : ''}
        </div>
        ${online ? `
          <div class="row wrap" style="gap:14px">
            <span><b class="num">${snap.totals.running}</b> <span class="muted">running</span></span>
            <span><b class="num">${snap.totals.paused}</b> <span class="muted">paused</span></span>
            <span><b class="num">${snap.totals.queued}</b> <span class="muted">queued</span></span>
            ${snap.totals.errors > 0 ? `<span class="badge error">${snap.totals.errors} error(s)</span>` : ''}
          </div>
          <div class="spread faint">
            <span>RAC <b class="num muted">${fmtCredit(snap.totals.rac)}</b></span>
            <span>Credit <b class="num muted">${fmtCredit(snap.totals.credit)}</b></span>
          </div>
          <div class="row wrap" style="gap:6px">${projects}</div>
        ` : `
          <div class="empty" style="padding:18px 10px">
            <b>${snap?.error || 'Waiting for connection...'}</b>
          </div>
        `}
      </div>
    `
  }

  let feed = []
  for (const [hid, s] of Object.entries(state.snaps)) {
    if (!s.online) continue
    for (const m of (s.messages || []).slice(-8)) {
      feed.push({ ...m, hostId: hid })
    }
  }
  feed.sort((a, b) => b.time - a.time || b.seq - a.seq)
  const recent = feed.slice(0, 10)

  return `<div class="content-inner">
    <div class="stats-grid">
      <div class="card hoverable">
        <div class="row"><div class="stat-icon"><svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="2" y="2" width="20" height="8" rx="2"/><rect x="2" y="14" width="20" height="8" rx="2"/></svg></div><div><div class="stat-label">Servers</div><div class="stat-value num">${list.length}/${state.hosts.length}</div><div class="faint">${onlineCount()} online</div></div></div>
      </div>
      <div class="card hoverable">
        <div class="row"><div class="stat-icon alt"><svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="22 12 18 12 15 21 9 3 6 12 2 12"/></svg></div><div><div class="stat-label">Active Tasks</div><div class="stat-value num">${running}</div><div class="faint">${queued} pending</div></div></div>
      </div>
      <div class="card hoverable">
        <div class="row"><div class="stat-icon soft"><svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 20V10"/><path d="M18 20V4"/><path d="M6 20v-4"/></svg></div><div><div class="stat-label">Fleet RAC</div><div class="stat-value num gain-pos">${fmtCredit(rac)}</div><div class="faint">recent average credit</div></div></div>
      </div>
      <div class="card hoverable">
        <div class="row"><div class="stat-icon alt"><svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M6 9H4.5a2.5 2.5 0 010-5H6"/><path d="M18 9h1.5a2.5 2.5 0 000-5H18"/><path d="M4 22h16"/><path d="M10 14.66V17c0 .55-.47.98-.97 1.21C7.85 18.75 7 20 7 22"/><path d="M14 14.66V17c0 .55.47.98.97 1.21C16.15 18.75 17 20 17 22"/><path d="M18 2H6v7a6 6 0 0012 0V2z"/></svg></div><div><div class="stat-label">Total Credit</div><div class="stat-value num">${fmtCredit(credit)}</div><div class="faint">across all projects</div></div></div>
      </div>
    </div>
    <div class="cards-grid">${serverCards || '<div class="card"><div class="empty"><b>No servers configured</b><span class="faint">Add a server to get started</span></div></div>'}</div>
    ${recent.length > 0 ? `
      <div class="card">
        <div class="card-head"><h2 class="card-title">Recent Activity</h2></div>
        ${recent.map(m => `<div class="msg-line"><span class="msg-pri-${Math.min(m.pri || 1, 3)} num faint" style="flex-shrink:0">${fmtAgo(m.time)}</span><span class="grow trunc" title="${esc(m.body)}">${esc(m.body)}</span></div>`).join('')}
      </div>
    ` : ''}
  </div>`
}

function onlineCount() { return Object.values(state.snaps).filter(s => s.online).length }

function renderTasks() {
  let rows = ''
  for (const [hid, snap] of Object.entries(state.snaps)) {
    if (!snap.online) continue
    const host = state.hosts.find(h => h.id === hid)
    for (const t of snap.tasks) {
      const pct = Math.round((t.progress || 0) * 100)
      rows += `<tr>
        <td class="task-name trunc" title="${esc(t.name)}">${esc(t.name)}</td>
        <td>${esc(t.projectName || '-')}</td>
        <td>${statusBadge(t.status)}</td>
        <td><div class="progress"><div class="fill" style="width:${pct}%"></div></div><span class="faint num">${pct}%</span></td>
        <td class="num">${fmtDuration(t.elapsed)}</td>
        <td class="num">${fmtDuration(t.eta)}</td>
        <td>${esc(t.resources || '-')}</td>
        <td><span class="chip plain trunc">${esc(host?.name || hid)}</span></td>
      </tr>`
    }
  }
  return `<div class="content-inner">
    <div class="card">
      <div class="card-head"><h2 class="card-title">All Tasks</h2></div>
      ${rows ? `<div class="tbl-wrap"><table class="tbl"><thead><tr><th>Name</th><th>Project</th><th>Status</th><th>Progress</th><th>Elapsed</th><th>ETA</th><th>Resources</th><th>Server</th></tr></thead><tbody>${rows}</tbody></table></div>` : '<div class="empty"><b>No tasks</b></div>'}
    </div>
  </div>`
}

function renderProjects() {
  let cards = ''
  for (const [hid, snap] of Object.entries(state.snaps)) {
    if (!snap.online) continue
    const host = state.hosts.find(h => h.id === hid)
    for (const p of snap.projects) {
      cards += `
        <div class="card hoverable" style="display:flex;flex-direction:column;gap:10px">
          <div class="row">
            <div class="proj-avatar" style="background:${projColor(p.url)};width:36px;height:36px;border-radius:11px;display:grid;place-items:center;color:#fff;font-weight:700;font-size:14px;flex-shrink:0">${(p.name || '?')[0]}</div>
            <div style="min-width:0"><div class="trunc" style="font-weight:700;font-size:14px">${esc(p.name)}</div><div class="faint num trunc" style="font-size:11px">${esc(p.url)}</div></div>
          </div>
          <div class="spread faint num" style="font-size:12px">
            <span>Host RAC: <b>${fmtCredit(p.hostRac)}</b></span>
            <span>Host Credit: <b>${fmtCredit(p.hostCredit)}</b></span>
          </div>
          <div class="row wrap" style="gap:6px">
            <span class="chip plain">${esc(host?.name || hid)}</span>
            ${p.suspended ? '<span class="badge paused">suspended</span>' : ''}
          </div>
        </div>
      `
    }
  }
  return `<div class="content-inner">
    <div class="cards-grid">${cards || '<div class="card"><div class="empty"><b>No projects</b></div></div>'}</div>
  </div>`
}

function renderTransfers() {
  let rows = ''
  for (const [hid, snap] of Object.entries(state.snaps)) {
    if (!snap.online) continue
    const host = state.hosts.find(h => h.id === hid)
    for (const t of snap.transfers) {
      const pct = Math.round((t.progress || 0) * 100)
      rows += `<tr>
        <td class="trunc" title="${esc(t.name)}" style="max-width:220px">${esc(t.name)}</td>
        <td>${esc(t.projectName || '-')}</td>
        <td><span class="badge ${t.upload ? 'upload' : 'download'}">${t.upload ? 'Upload' : 'Download'}</span></td>
        <td><div class="progress"><div class="fill" style="width:${pct}%"></div></div><span class="faint num">${pct}%</span></td>
        <td class="num">${fmtBytes(t.done)} / ${fmtBytes(t.total)}</td>
        <td><span class="chip plain trunc">${esc(host?.name || hid)}</span></td>
      </tr>`
    }
  }
  return `<div class="content-inner">
    <div class="card">
      <div class="card-head"><h2 class="card-title">Transfers</h2></div>
      ${rows ? `<div class="tbl-wrap"><table class="tbl"><thead><tr><th>Name</th><th>Project</th><th>Type</th><th>Progress</th><th>Size</th><th>Server</th></tr></thead><tbody>${rows}</tbody></table></div>` : '<div class="empty"><b>No active transfers</b></div>'}
    </div>
  </div>`
}

function renderMessages() {
  let msgs = []
  for (const [hid, snap] of Object.entries(state.snaps)) {
    if (!snap.online) continue
    for (const m of (snap.messages || [])) {
      msgs.push({ ...m, hostId: hid })
    }
  }
  msgs.sort((a, b) => b.time - a.time || b.seq - a.seq)

  return `<div class="content-inner">
    <div class="card">
      <div class="card-head"><h2 class="card-title">Messages</h2></div>
      ${msgs.length > 0 ? msgs.slice(0, 50).map(m => `<div class="msg-line"><span class="msg-pri-${Math.min(m.pri || 1, 3)} num faint" style="flex-shrink:0">${fmtAgo(m.time)}</span><span class="grow trunc" title="${esc(m.body)}">${esc(m.body)}</span><span class="chip plain" style="flex-shrink:0">${esc(state.snaps[m.hostId]?.host || m.hostId)}:${state.snaps[m.hostId]?.port || ''}</span></div>`).join('') : '<div class="empty"><b>No messages</b></div>'}
    </div>
  </div>`
}

function svgLineChart(data, width, height, color) {
  if (!data || data.length < 2) return ''
  const pad = { l: 50, r: 10, t: 10, b: 28 }
  const cw = width - pad.l - pad.r
  const ch = height - pad.t - pad.b
  const vals = data.map(d => d.v)
  const maxV = Math.max(...vals, 1)
  const niceMax = Math.ceil(maxV / Math.pow(10, Math.floor(Math.log10(maxV)))) * Math.pow(10, Math.floor(Math.log10(maxV)))
  const yMax = (niceMax < maxV ? niceMax * 2 : niceMax) || maxV

  const pts = data.map((d, i) => {
    const x = pad.l + (i / (data.length - 1)) * cw
    const y = pad.t + ch - (d.v / yMax) * ch
    return `${x},${y}`
  })

  const areaPath = `M${pad.l},${pad.t + ch} L${pts.join(' L')} L${pad.l + cw},${pad.t + ch} Z`
  const linePath = `M${pts.join(' L')}`

  const gradId = 'g' + Math.random().toString(36).substr(2, 6)

  let gridLines = ''
  for (let i = 0; i <= 4; i++) {
    const y = pad.t + (ch / 4) * i
    const v = yMax - (yMax / 4) * i
    gridLines += `<line x1="${pad.l}" y1="${y}" x2="${pad.l + cw}" y2="${y}" stroke="var(--border)" stroke-width="0.8" stroke-dasharray="4,4"/>`
    gridLines += `<text x="${pad.l - 6}" y="${y + 4}" text-anchor="end" fill="var(--text-faint)" font-size="10" font-family="JetBrains Mono,monospace">${fmtCredit(v)}</text>`
  }

  const firstLabel = data[0]?.l || ''
  const midLabel = data[Math.floor(data.length / 2)]?.l || ''
  const lastLabel = data[data.length - 1]?.l || ''

  return `<svg viewBox="0 0 ${width} ${height}" class="linechart" preserveAspectRatio="xMidYMid meet">
    <defs><linearGradient id="${gradId}" x1="0" y1="0" x2="0" y2="1"><stop offset="0%" stop-color="${color}" stop-opacity="0.3"/><stop offset="100%" stop-color="${color}" stop-opacity="0.02"/></linearGradient></defs>
    ${gridLines}
    <path d="${areaPath}" fill="url(#${gradId})"/>
    <path d="${linePath}" fill="none" stroke="${color}" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"/>
    <text x="${pad.l}" y="${height - 4}" fill="var(--text-faint)" font-size="10" font-family="JetBrains Mono,monospace">${firstLabel}</text>
    <text x="${pad.l + cw / 2}" y="${height - 4}" text-anchor="middle" fill="var(--text-faint)" font-size="10" font-family="JetBrains Mono,monospace">${midLabel}</text>
    <text x="${pad.l + cw}" y="${height - 4}" text-anchor="end" fill="var(--text-faint)" font-size="10" font-family="JetBrains Mono,monospace">${lastLabel}</text>
  </svg>`
}

function svgBarChart(data, width, height, upColor, downColor) {
  if (!data || data.length === 0) return ''
  const pad = { l: 50, r: 10, t: 10, b: 28 }
  const cw = width - pad.l - pad.r
  const ch = height - pad.t - pad.b
  const maxV = Math.max(...data.map(d => Math.max(d.up, d.down)), 1)
  const yMax = maxV * 1.15
  const barW = Math.max(2, (cw / data.length) - 2)

  let bars = ''
  let gridLines = ''
  for (let i = 0; i <= 4; i++) {
    const y = pad.t + (ch / 4) * i
    const v = yMax - (yMax / 4) * i
    gridLines += `<line x1="${pad.l}" y1="${y}" x2="${pad.l + cw}" y2="${y}" stroke="var(--border)" stroke-width="0.8" stroke-dasharray="4,4"/>`
    gridLines += `<text x="${pad.l - 6}" y="${y + 4}" text-anchor="end" fill="var(--text-faint)" font-size="10" font-family="JetBrains Mono,monospace">${fmtBytes(Math.round(v))}</text>`
  }

  data.forEach((d, i) => {
    const x = pad.l + (i / data.length) * cw + 1
    const hUp = (d.up / yMax) * ch
    const hDown = (d.down / yMax) * ch
    bars += `<rect x="${x}" y="${pad.t + ch - hUp}" width="${barW}" height="${hUp}" fill="${upColor}" rx="2" opacity="0.75"/>`
    bars += `<rect x="${x + barW + 1}" y="${pad.t + ch - hDown}" width="${barW}" height="${hDown}" fill="${downColor}" rx="2" opacity="0.75"/>`
  })

  return `<svg viewBox="0 0 ${width} ${height}" class="linechart" preserveAspectRatio="xMidYMid meet">
    ${gridLines}${bars}
  </svg>`
}

async function loadStats() {
  for (const h of state.hosts) {
    try { state.stats[h.id] = await api('GetStats', h.id) } catch (e) { state.stats[h.id] = [] }
    try { state.xfers[h.id] = await api('GetXferHistory', h.id) } catch (e) { state.xfers[h.id] = [] }
    try { state.disk[h.id] = await api('GetDiskUsage', h.id) } catch (e) { state.disk[h.id] = null }
  }
}

function renderStats() {
  loadStats().then(() => {
    const c = $('#content')
    if (c && state.page === 'stats') c.innerHTML = renderStatsInner()
  })

  return `<div class="content-inner"><div class="card"><div class="empty"><b>Loading stats...</b></div></div></div>`
}

function renderStatsInner() {
  const PALETTE = ['#7c3aed', '#6366f1', '#a78bfa', '#818cf8', '#c084fc', '#8b5cf6', '#4f46e5', '#6d28d9']
  let creditCharts = ''
  let totalCredit = 0
  let totalRAC = 0

  for (const h of state.hosts) {
    const series = state.stats[h.id] || []
    for (const s of series) {
      const color = PALETTE[Math.abs(s.url.length * 47) % PALETTE.length]
      const data = (s.daily || []).map(p => ({ v: p.hostCredit, l: p.day.slice(4, 6) + '/' + p.day.slice(6, 8) }))
      if (data.length > 0) totalCredit = Math.max(totalCredit, data[data.length - 1].v)
      creditCharts += `
        <div class="card" style="display:flex;flex-direction:column;gap:8px">
          <div class="card-head"><h3 class="card-title" style="font-size:14px"><span style="display:inline-block;width:10px;height:10px;border-radius:3px;background:${color};flex-shrink:0"></span> ${esc(s.name || s.url)}</h3><span class="chip plain">${esc(h.name)}</span></div>
          ${svgLineChart(data, 800, 220, color)}
        </div>
      `
    }
  }

  let xferData = []
  let xferChart = ''
  for (const h of state.hosts) {
    const xf = state.xfers[h.id] || []
    xferData = xf.map(d => ({ up: d.up, down: d.down, l: new Date(d.when * 1000).toLocaleDateString('en-US', { month: 'short', day: 'numeric' }) }))
  }
  if (xferData.length > 0) {
    xferChart = `
      <div class="card" style="display:flex;flex-direction:column;gap:8px">
        <div class="card-head"><h3 class="card-title" style="font-size:14px">Transfer History (30 days)</h3></div>
        <div class="row" style="gap:16px;padding-left:16px;margin-bottom:4px">
          <span class="row" style="gap:5px"><span style="width:10px;height:10px;border-radius:3px;background:#6366f1"></span><span class="faint" style="font-size:11px">Download</span></span>
          <span class="row" style="gap:5px"><span style="width:10px;height:10px;border-radius:3px;background:#a78bfa"></span><span class="faint" style="font-size:11px">Upload</span></span>
        </div>
        ${svgBarChart(xferData, 800, 220, '#a78bfa', '#6366f1')}
      </div>
    `
  }

  let diskCards = ''
  for (const h of state.hosts) {
    const du = state.disk[h.id]
    if (!du) continue
    const used = du.total - du.free
    const pct = du.total > 0 ? Math.round((used / du.total) * 100) : 0
    let projectBars = ''
    const colors = ['#7c3aed', '#6366f1', '#a78bfa', '#818cf8', '#c084fc', '#8b5cf6']
    for (let i = 0; i < (du.projects || []).length; i++) {
      const p = du.projects[i]
      const ppct = du.total > 0 ? Math.round((p.diskUsage / du.total) * 100) : 0
      projectBars += `<div style="display:flex;align-items:center;gap:8px;font-size:12px">
        <span style="width:8px;height:8px;border-radius:3px;background:${colors[i % colors.length]};flex-shrink:0"></span>
        <span class="grow trunc">${esc(p.url.replace('https://', '').replace('http://', '').split('/')[0])}</span>
        <span class="num faint">${fmtBytes(p.diskUsage)} (${ppct}%)</span>
      </div>`
    }
    diskCards += `
      <div class="card" style="display:flex;flex-direction:column;gap:10px">
        <div class="card-head"><h3 class="card-title" style="font-size:14px">Disk Usage</h3><span class="chip plain">${esc(h.name)}</span></div>
        <div class="progress striped" style="height:10px"><div class="fill" style="width:${pct}%"></div></div>
        <div class="spread faint num" style="font-size:12px"><span>${fmtBytes(used)} used</span><span>${fmtBytes(du.free)} free</span><span>${fmtBytes(du.total)} total</span></div>
        ${projectBars ? `<div style="display:flex;flex-direction:column;gap:5px;margin-top:6px">${projectBars}</div>` : ''}
      </div>
    `
  }

  const hasData = creditCharts || xferChart || diskCards
  return `<div class="content-inner">
    ${creditCharts ? `<div class="section-title"><h2>Credit History</h2></div><div style="display:flex;flex-direction:column;gap:14px">${creditCharts}</div>` : ''}
    ${xferChart ? `<div style="display:flex;flex-direction:column;gap:14px">${xferChart}</div>` : ''}
    ${diskCards ? `<div class="section-title"><h2>Disk Usage</h2></div><div class="cards-grid">${diskCards}</div>` : ''}
    ${!hasData ? '<div class="card"><div class="empty"><b>No stats available</b><span class="faint">Connect to a host to see statistics</span></div></div>' : ''}
  </div>`
}

function renderHosts() {
  let cards = ''
  for (const h of state.hosts) {
    const snap = state.snaps[h.id]
    const online = snap?.online
    cards += `
      <div class="card hoverable" style="display:flex;flex-direction:column;gap:10px">
        <div class="spread">
          <div class="row">
            <span class="dot ${online ? 'on' : 'off'}"></span>
            <div><div style="font-weight:700;font-size:15px">${esc(h.name)}</div><div class="faint num">${esc(h.host)}:${h.port}</div></div>
          </div>
          <button class="btn sm danger" onclick="window._removeHost('${h.id}')">Remove</button>
        </div>
        ${online ? `
          <div class="faint num" style="font-size:12px">${esc(snap.hostInfo?.os || '')} · ${esc(snap.hostInfo?.cpu || '')} · ${snap.hostInfo?.cores || 0} cores</div>
          ${snap.hostInfo?.gpus?.length ? `<div class="faint" style="font-size:12px">GPU: ${snap.hostInfo.gpus.map(g => g.names.join(', ')).join(' | ')}</div>` : ''}
        ` : `<div class="faint">${snap?.error || 'Offline'}</div>`}
      </div>
    `
  }
  return `<div class="content-inner">
    <div style="display:flex;justify-content:flex-end"><button class="btn primary" onclick="window._showAddHost()">+ Add Server</button></div>
    <div class="cards-grid">${cards || '<div class="card"><div class="empty"><b>No servers</b><span class="faint">Add a server to manage</span></div></div>'}</div>
    ${state.modal || ''}
  </div>`
}

function renderSettings() {
  const di = state.daemonInfo
  const found = di?.found
  const exe = di?.exe || 'Not found'
  const dataDir = di?.dataDir || 'N/A'
  const hint = di?.hint || ''

  let hwCards = ''
  for (const [hid, snap] of Object.entries(state.snaps)) {
    if (!snap?.online || !snap.hostInfo) continue
    const hi = snap.hostInfo
    const host = state.hosts.find(h => h.id === hid)
    const gpus = (hi.GPUs || []).map(g => `<div class="row" style="gap:6px;font-size:12px"><span style="width:8px;height:8px;border-radius:3px;background:#7c3aed;flex-shrink:0"></span><span class="grow">${esc((g.Names || []).join(', '))}</span><span class="num faint">${fmtBytes(g.VRAM || 0)}</span></div>`).join('')
    hwCards += `
      <div class="card" style="display:flex;flex-direction:column;gap:10px">
        <div class="card-head"><h3 class="card-title" style="font-size:14px">${esc(host?.name || hid)}</h3><span class="chip indigo num">v${esc(snap.version || '?')}</span></div>
        <div style="display:grid;grid-template-columns:1fr 1fr;gap:8px 16px;font-size:12px">
          <div class="faint">OS</div><div class="num">${esc(hi.OS || '?')} ${esc(hi.OSVersion || '')}</div>
          <div class="faint">CPU</div><div class="num">${esc(hi.CPU || '?')}</div>
          <div class="faint">Cores</div><div class="num">${hi.Cores || 0}</div>
          <div class="faint">FLOPS</div><div class="num">${fmtFlops(hi.Flops)}</div>
          <div class="faint">RAM</div><div class="num">${fmtBytes(hi.Memory || 0)}</div>
          <div class="faint">Disk</div><div class="num">${fmtBytes(hi.DiskTotal || 0)} total, ${fmtBytes(hi.DiskFree || 0)} free</div>
        </div>
        ${gpus ? `<div style="display:flex;flex-direction:column;gap:5px;margin-top:4px"><div class="faint" style="font-size:11px;text-transform:uppercase;letter-spacing:0.5px">GPUs</div>${gpus}</div>` : ''}
      </div>
    `
  }

  return `<div class="content-inner">
    <div class="card">
      <div class="card-head"><h2 class="card-title">Local Camellia Client</h2></div>
      <div style="display:flex;flex-direction:column;gap:12px">
        <div class="row" style="gap:12px;flex-wrap:wrap">
          <span class="dot ${state.daemonStatus === 'running' ? 'on' : 'off'}"></span>
          <span style="font-weight:700">Status: <span class="num">${esc(state.daemonStatus)}</span></span>
        </div>
        <div class="faint" style="font-size:12px">Path: <span class="num">${esc(exe)}</span></div>
        <div class="faint" style="font-size:12px">Data: <span class="num">${esc(dataDir)}</span></div>
        ${hint ? `<div class="faint" style="font-size:12px">${esc(hint)}</div>` : ''}
        <div class="row" style="gap:8px;margin-top:6px">
          <button class="btn primary" onclick="window._startDaemon()" ${!found ? 'disabled' : ''}>Start</button>
          <button class="btn danger" onclick="window._stopDaemon()" ${!found ? 'disabled' : ''}>Stop</button>
        </div>
      </div>
    </div>
    ${hwCards ? `<div class="section-title"><h2>Hardware</h2></div><div class="cards-grid">${hwCards}</div>` : ''}
    <div class="card">
      <div class="card-head"><h2 class="card-title">About</h2></div>
      <div style="color:var(--text-soft)">
        <p><b>LavandeGrid v1.0.0</b> - Desktop Camellia Manager</p>
        <p>Built with Wails + Go backend</p>
        <p>Coded by <b>Alperen Yavuz</b></p>
      </div>
    </div>
  </div>`
}

function fmtFlops(v) {
  if (!v) return '?'
  const f = typeof v === 'number' ? v : parseFloat(v) || 0
  if (f >= 1e12) return (f / 1e12).toFixed(1) + ' TFLOPS'
  if (f >= 1e9) return (f / 1e9).toFixed(1) + ' GFLOPS'
  if (f >= 1e6) return (f / 1e6).toFixed(1) + ' MFLOPS'
  return f.toFixed(0) + ' FLOPS'
}

window._setPage = setPage

window._toggleTheme = () => {
  state.theme = state.theme === 'dark' ? 'light' : 'dark'
  localStorage.setItem('lavande-theme', state.theme)
  applyTheme()
  render()
}

function applyTheme() {
  document.documentElement.setAttribute('data-theme', state.theme)
}

window._removeHost = async (id) => {
  if (!confirm('Remove this server?')) return
  await api('RemoveHost', id)
  delete state.snaps[id]
  await refreshHosts()
  toast('Server removed', 'ok')
}

window._showAddHost = () => {
  state.modal = `
    <div class="modal-overlay" onclick="if(event.target===this)window._closeModal()">
      <div class="modal" style="max-width:420px">
        <div class="modal-head"><h3>Add Server</h3><button class="btn icon" onclick="window._closeModal()">✕</button></div>
        <div class="field"><label class="label">Name</label><input class="input" id="add-name" placeholder="My Server"></div>
        <div class="field"><label class="label">Host</label><input class="input" id="add-host" placeholder="192.168.0.5"></div>
        <div class="field"><label class="label">Port</label><input class="input num" id="add-port" value="31416"></div>
        <div class="field"><label class="label">Password</label><input class="input" id="add-pass" type="password" placeholder="gui_rpc_auth.cfg password"></div>
        <div class="modal-foot">
          <button class="btn" onclick="window._closeModal()">Cancel</button>
          <button class="btn primary" onclick="window._doAddHost()">Add</button>
        </div>
      </div>
    </div>
  `
  render()
}

window._closeModal = () => { state.modal = null; render() }

window._doAddHost = async () => {
  const name = $('#add-name')?.value || 'Unnamed'
  const host = $('#add-host')?.value || 'localhost'
  const port = parseInt($('#add-port')?.value) || 31416
  const pass = $('#add-pass')?.value || ''
  try {
    await api('AddHost', name, host, port, pass)
    state.modal = null
    await refreshHosts()
    toast('Server added', 'ok')
  } catch (e) {
    toast(e.message || 'Failed to add server', 'err')
  }
}

window._startDaemon = async () => {
  try {
    await api('StartDaemon')
    state.daemonStatus = await api('GetDaemonStatus')
    render()
    toast('Client started', 'ok')
  } catch (e) {
    toast(e.message || 'Failed to start', 'err')
  }
}

window._stopDaemon = async () => {
  try {
    await api('StopDaemon')
    state.daemonStatus = await api('GetDaemonStatus')
    render()
    toast('Client stopped', 'ok')
  } catch (e) {
    toast(e.message || 'Failed to stop', 'err')
  }
}

async function init() {
  applyTheme()
  render()
  try {
    state.daemonStatus = await api('GetDaemonStatus')
    state.daemonInfo = await api('DetectDaemon')
  } catch (e) {}
  await refreshHosts()
  setInterval(refreshHosts, 4000)
}

init()
