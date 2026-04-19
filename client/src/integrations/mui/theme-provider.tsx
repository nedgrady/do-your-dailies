import type { PaletteMode } from '@mui/material'
import { CssBaseline, ThemeProvider, createTheme } from '@mui/material'
import type { ReactNode } from 'react'
import { createContext, useCallback, useEffect, useMemo, useState } from 'react'

const colorModeStorageKey = 'dyd-color-mode-preference'

type ColorModePreference = PaletteMode | 'system'

interface ColorModeContextValue {
  mode: PaletteMode
  preference: ColorModePreference
  toggleMode: () => void
}

export const ColorModeContext = createContext<ColorModeContextValue>({
  mode: 'light',
  preference: 'system',
  toggleMode: () => undefined,
})

function getSystemMode(): PaletteMode {
  if (
    typeof window === 'undefined' ||
    typeof window.matchMedia !== 'function'
  ) {
    return 'light'
  }

  return window.matchMedia('(prefers-color-scheme: dark)').matches
    ? 'dark'
    : 'light'
}

function readStoredPreference(): ColorModePreference {
  if (typeof window === 'undefined') {
    return 'system'
  }

  const storedValue = window.localStorage.getItem(colorModeStorageKey)
  if (
    storedValue === 'light' ||
    storedValue === 'dark' ||
    storedValue === 'system'
  ) {
    return storedValue
  }

  return 'system'
}

export function AppThemeProvider({ children }: { children: ReactNode }) {
  const [preference, setPreference] = useState<ColorModePreference>(() =>
    readStoredPreference(),
  )
  const [systemMode, setSystemMode] = useState<PaletteMode>(() =>
    getSystemMode(),
  )

  useEffect(() => {
    if (
      typeof window === 'undefined' ||
      typeof window.matchMedia !== 'function'
    ) {
      return
    }

    const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)')
    setSystemMode(mediaQuery.matches ? 'dark' : 'light')

    const handleSystemModeChange = (event: MediaQueryListEvent) => {
      setSystemMode(event.matches ? 'dark' : 'light')
    }

    mediaQuery.addEventListener('change', handleSystemModeChange)

    return () => {
      mediaQuery.removeEventListener('change', handleSystemModeChange)
    }
  }, [])

  useEffect(() => {
    if (typeof window === 'undefined') {
      return
    }

    window.localStorage.setItem(colorModeStorageKey, preference)
  }, [preference])

  const toggleMode = useCallback(() => {
    setPreference((currentPreference) => {
      const resolvedMode =
        currentPreference === 'system' ? systemMode : currentPreference

      return resolvedMode === 'dark' ? 'light' : 'dark'
    })
  }, [systemMode])

  const mode = preference === 'system' ? systemMode : preference

  const theme = useMemo(
    () =>
      createTheme({
        palette: {
          mode,
        },
        typography: {
          fontFamily:
            "'Manrope', system-ui, -apple-system, 'Segoe UI', Roboto, 'Helvetica Neue', Arial",
        },
      }),
    [mode],
  )

  const contextValue = useMemo<ColorModeContextValue>(
    () => ({
      mode,
      preference,
      toggleMode,
    }),
    [mode, preference, toggleMode],
  )

  return (
    <ColorModeContext.Provider value={contextValue}>
      <ThemeProvider theme={theme}>
        <CssBaseline />
        {children}
      </ThemeProvider>
    </ColorModeContext.Provider>
  )
}
