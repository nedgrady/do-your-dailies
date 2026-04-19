import type { PaletteMode } from '@mui/material'
import { CssBaseline, ThemeProvider, createTheme } from '@mui/material'
import type { ReactNode } from 'react'
import React, { createContext } from 'react'

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

interface ThemeProviderState {
  preference: ColorModePreference
  systemMode: PaletteMode
}

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

export class AppThemeProvider extends React.Component<
  { children: ReactNode },
  ThemeProviderState
> {
  private mediaQuery?: MediaQueryList

  public constructor(props: { children: ReactNode }) {
    super(props)
    this.state = {
      preference: readStoredPreference(),
      systemMode: getSystemMode(),
    }
  }

  public componentDidMount() {
    if (
      typeof window === 'undefined' ||
      typeof window.matchMedia !== 'function'
    ) {
      return
    }

    this.mediaQuery = window.matchMedia('(prefers-color-scheme: dark)')
    this.setState({ systemMode: this.mediaQuery.matches ? 'dark' : 'light' })
    this.mediaQuery.addEventListener('change', this.handleSystemModeChange)
  }

  public componentDidUpdate(
    _: { children: ReactNode },
    prevState: ThemeProviderState,
  ) {
    if (typeof window === 'undefined') {
      return
    }

    if (prevState.preference !== this.state.preference) {
      window.localStorage.setItem(colorModeStorageKey, this.state.preference)
    }
  }

  public componentWillUnmount() {
    this.mediaQuery?.removeEventListener('change', this.handleSystemModeChange)
  }

  private readonly handleSystemModeChange = (event: MediaQueryListEvent) => {
    this.setState({ systemMode: event.matches ? 'dark' : 'light' })
  }

  private readonly toggleMode = () => {
    this.setState((currentState) => {
      const resolvedMode =
        currentState.preference === 'system'
          ? currentState.systemMode
          : currentState.preference

      return {
        preference: resolvedMode === 'dark' ? 'light' : 'dark',
      }
    })
  }

  public render() {
    const mode =
      this.state.preference === 'system'
        ? this.state.systemMode
        : this.state.preference

    const theme = createTheme({
      palette: {
        mode,
      },
      typography: {
        fontFamily:
          "'Manrope', system-ui, -apple-system, 'Segoe UI', Roboto, 'Helvetica Neue', Arial",
      },
    })

    const contextValue: ColorModeContextValue = {
      mode,
      preference: this.state.preference,
      toggleMode: this.toggleMode,
    }

    return (
      <ColorModeContext.Provider value={contextValue}>
        <ThemeProvider theme={theme}>
          <CssBaseline />
          {this.props.children}
        </ThemeProvider>
      </ColorModeContext.Provider>
    )
  }
}
