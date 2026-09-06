import { createContext, useContext, useEffect, useMemo, useState } from 'react'
import type { ReactNode } from 'react'
import axiosInstance from './axios'

// Matches server/internal/auth's SimulatedUserIDHeader. Only honored by the
// server when SIMULATION_MODE_ENABLED is set, so sending this header
// elsewhere is harmless - it's simply ignored.
export const SIMULATED_USER_ID_HEADER = 'X-Simulated-User-Id'

interface SimulatedUserContextValue {
  userId: number | null
  setUserId: (id: number | null) => void
}

// Default (no provider above): "acting as self" - real pages read userId as
// null, and never call setUserId, so its default just guards misuse.
const SimulatedUserContext = createContext<SimulatedUserContextValue>({
  userId: null,
  setUserId: () => {
    throw new Error(
      'setSimulatedUser called without a SimulatedUserProvider above it',
    )
  },
})

function applyHeader(id: number | null) {
  if (id === null) {
    delete axiosInstance.defaults.headers.common[SIMULATED_USER_ID_HEADER]
  } else {
    axiosInstance.defaults.headers.common[SIMULATED_USER_ID_HEADER] = String(id)
  }
}

// Owns the impersonation state and keeps axios in sync with it. The header
// update happens inside setUserId itself, synchronously - not in a separate
// effect - so it's always applied before React even schedules the re-render
// that mounts any consumer's data-fetching components. Effects (even
// layout effects) fire bottom-up post-commit, so a parent's effect can never
// reliably beat a child's - doing the sync here sidesteps that entirely.
//
// Scoped locally to the simulator page (not the whole app): leaving the page
// unmounts this provider, whose own cleanup below clears the header - real
// pages read the outer default (userId: null) the instant they mount, with
// no separate route-based cleanup needed.
export function SimulatedUserProvider({ children }: { children: ReactNode }) {
  const [userId, setUserIdState] = useState<number | null>(null)

  const setUserId = (id: number | null) => {
    applyHeader(id)
    setUserIdState(id)
  }

  useEffect(() => () => applyHeader(null), [])

  const value = useMemo(() => ({ userId, setUserId }), [userId])

  return (
    <SimulatedUserContext.Provider value={value}>
      {children}
    </SimulatedUserContext.Provider>
  )
}

// Safe to call from anywhere, including outside the simulator (e.g. from the
// shared query hooks in domain/*.ts, which use this to scope cache keys by
// acting user) - it resolves to the "acting as self" default when there's no
// SimulatedUserProvider above.
export function useSimulatedUser(): SimulatedUserContextValue {
  return useContext(SimulatedUserContext)
}
