// @vitest-environment jsdom

import type { RenderOptions } from '@testing-library/react'
import { render } from '@testing-library/react'
import type { ReactElement, ReactNode } from 'react'
import { AppThemeProvider } from './integrations/mui/theme-provider'

function AllTheProviders({ children }: { children: ReactNode }) {
  return <AppThemeProvider>{children}</AppThemeProvider>
}

function customRender(
  ui: ReactElement,
  options?: Omit<RenderOptions, 'wrapper'>,
) {
  return render(ui, { wrapper: AllTheProviders, ...options })
}

export * from '@testing-library/react'
export { customRender as render }
