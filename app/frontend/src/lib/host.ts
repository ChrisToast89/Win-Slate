// Host bridge — exposes window.slate matching the Electron preload API.
// In the Wails binary, methods call Go bindings. In browser-only dev without
// Wails, installDevMock() still provides a stub.

import type {
  BrainBackend,
  BrainRequest,
  BrainResult,
  BrainStatus,
  CircledTake,
  LocalModelInfo,
  Project,
  ProjectMeta,
  SlateApi,
  AudioFingerprint
} from '../shared/types'

// Wails generates window.go.main.App at build time. We declare a loose shape
// so TypeScript compiles before `wails generate module` runs.
type GoApp = {
  ListProjects(): Promise<ProjectMeta[]>
  CreateProject(name: string): Promise<Project>
  OpenProject(id: string): Promise<Project | null>
  SaveProject(project: Project): Promise<void>
  DeleteProject(id: string): Promise<void>
  RevealProject(id: string): Promise<void>
  BrainStatus(localEndpoint: string): Promise<BrainStatus>
  LocalModels(endpoint: string): Promise<{ endpoint: string | null; models: LocalModelInfo[] }>
  BrainRun(req: BrainRequest): Promise<BrainResult>
  BrainCancel(id: string): Promise<void>
  BrainTest(backend: string, local: { endpoint?: string; model?: string }): Promise<BrainResult>
  StillsDiscover(): Promise<CircledTake[]>
  StillsExtract(
    projectId: string,
    mediaPath: string,
    inSec: number | null,
    outSec: number | null
  ): Promise<string[]>
  PickMedia(): Promise<string[]>
  PickAudio(): Promise<string[]>
  IngestMedia(projectId: string, path: string): Promise<{ kind: 'image' | 'video'; frames: string[] }>
  AnalyzeAudio(path: string): Promise<AudioFingerprint>
  CopyText(text: string): Promise<void>
  FileAsDataURL(path: string): Promise<string>
  AppVersion(): Promise<string>
  OpenExternal(url: string): Promise<void>
  GetPlatform(): Promise<string>
}

declare global {
  interface Window {
    go?: { main?: { App?: GoApp } }
    runtime?: {
      EventsOn(event: string, cb: (...args: unknown[]) => void): () => void
      EventsEmit(event: string, ...args: unknown[]): void
    }
    slate: SlateApi
  }
}

function goApp(): GoApp | null {
  return window.go?.main?.App ?? null
}

function eventsOn(event: string, cb: () => void): () => void {
  if (window.runtime?.EventsOn) {
    return window.runtime.EventsOn(event, () => cb())
  }
  return () => {}
}

/** Install the real Wails-backed slate API. Returns true if bindings exist. */
export function installWailsHost(): boolean {
  const app = goApp()
  if (!app) return false

  const api: SlateApi = {
    listProjects: () => app.ListProjects(),
    createProject: (name) => app.CreateProject(name),
    openProject: async (id) => {
      const p = await app.OpenProject(id)
      return p ?? null
    },
    saveProject: (project) => app.SaveProject(project),
    deleteProject: (id) => app.DeleteProject(id),
    revealProject: (id) => app.RevealProject(id),
    brainStatus: (localEndpoint?: string) => app.BrainStatus(localEndpoint ?? ''),
    localModels: (endpoint?: string) => app.LocalModels(endpoint ?? ''),
    stillsDiscover: () => app.StillsDiscover(),
    stillsExtract: (projectId, mediaPath, inSec, outSec) =>
      app.StillsExtract(projectId, mediaPath, inSec ?? null, outSec ?? null),
    brainRun: (req) => {
      // Backend is taken from project defaults in callers via brainRunWith when needed.
      return app.BrainRun(req)
    },
    brainRunWith: (req, backend) => app.BrainRun({ ...req, backend }),
    brainCancel: (id) => app.BrainCancel(id),
    brainTest: (backend, local) =>
      app.BrainTest(backend, { endpoint: local?.endpoint, model: local?.model }),
    pickMedia: () => app.PickMedia(),
    pickAudio: () => app.PickAudio(),
    ingestMedia: (projectId, path) => app.IngestMedia(projectId, path),
    analyzeAudio: (path) => app.AnalyzeAudio(path),
    pathForFile: () => {
      // WebView2 does not expose local paths for drag-dropped File objects.
      // Prefer pickMedia / pickAudio dialogs. Return empty so callers fall back.
      return ''
    },
    copyText: (text) => app.CopyText(text),
    onProjectsChanged: (cb) => eventsOn('projects:changed', cb),
    onAboutOpen: (cb) => eventsOn('about:open', cb),
    onHelpOpen: (cb) => eventsOn('help:open', cb),
    onBrainRefresh: (cb) => eventsOn('brain:refresh', cb),
    fileAsDataURL: (path) => app.FileAsDataURL(path),
    appVersion: () => app.AppVersion()
  }

  window.slate = api
  return true
}

/** Resolve a local disk path to something an <img> can load. */
const dataURLCache = new Map<string, string>()

export async function localMediaSrc(path: string): Promise<string> {
  if (!path) return ''
  if (path.startsWith('data:') || path.startsWith('http://') || path.startsWith('https://') || path.startsWith('blob:')) {
    return path
  }
  const cached = dataURLCache.get(path)
  if (cached) return cached
  if (window.slate && 'fileAsDataURL' in window.slate && typeof window.slate.fileAsDataURL === 'function') {
    try {
      const url = await window.slate.fileAsDataURL(path)
      dataURLCache.set(path, url)
      return url
    } catch {
      /* fall through */
    }
  }
  // Last resort (Electron-style); usually blocked in WebView2
  return `file:///${path.replace(/\\/g, '/')}`
}
