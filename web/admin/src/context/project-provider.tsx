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
  /** 当前选中的项目 ID，空字符串表示「全部」或未选 */
  projectID: string
  setProjectID: (id: string) => void
}

const ProjectContext = createContext<ProjectContextValue | null>(null)

function readStored(): string {
  if (typeof window === 'undefined') return ''
  try {
    return localStorage.getItem(STORAGE_KEY) ?? ''
  } catch {
    return ''
  }
}

export function ProjectProvider({ children }: { children: ReactNode }) {
  const [projectID, setProjectIDState] = useState<string>(readStored)

  const setProjectID = useCallback((id: string) => {
    setProjectIDState(id)
    try {
      if (id) {
        localStorage.setItem(STORAGE_KEY, id)
      } else {
        localStorage.removeItem(STORAGE_KEY)
      }
    } catch {
      // ignore
    }
  }, [])

  const value = useMemo(
    () => ({ projectID, setProjectID }),
    [projectID, setProjectID]
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
