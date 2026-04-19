import {
  createContext,
  useCallback,
  useContext,
  useMemo,
  useState,
  type ReactNode,
} from 'react'

const STORAGE_KEY = 'app_current_project_id'

type ProjectContextValue = {
  /** 当前选中的项目 ID，0 表示「全部」或未选 */
  projectId: number
  setProjectId: (id: number) => void
}

const ProjectContext = createContext<ProjectContextValue | null>(null)

function readStored(): number {
  if (typeof window === 'undefined') return 0
  try {
    const raw = localStorage.getItem(STORAGE_KEY)
    return raw ? Number(raw) || 0 : 0
  } catch {
    return 0
  }
}

export function ProjectProvider({ children }: { children: ReactNode }) {
  const [projectId, setProjectIdState] = useState<number>(readStored)

  const setProjectId = useCallback((id: number) => {
    setProjectIdState(id)
    try {
      if (id > 0) {
        localStorage.setItem(STORAGE_KEY, String(id))
      } else {
        localStorage.removeItem(STORAGE_KEY)
      }
    } catch {
      // ignore
    }
  }, [])

  const value = useMemo(
    () => ({ projectId, setProjectId }),
    [projectId, setProjectId]
  )

  return (
    <ProjectContext.Provider value={value}>
      {children}
    </ProjectContext.Provider>
  )
}

export function useProject() {
  const ctx = useContext(ProjectContext)
  if (ctx == null) {
    throw new Error('useProject must be used within ProjectProvider')
  }
  return ctx
}
