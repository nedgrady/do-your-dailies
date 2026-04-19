// @vitest-environment jsdom

import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import type { RenderOptions } from '@testing-library/react'
import { render } from '@testing-library/react'
import type { ReactElement, ReactNode } from 'react'
import { AppThemeProvider } from './integrations/mui/theme-provider'

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: false,
      staleTime: 0,
      gcTime: 0,
      refetchInterval: false,
    },
    mutations: {
      retry: false,
      gcTime: 0,
    },
  },
})

function AllTheProviders({ children }: { children: ReactNode }) {
  return (
    <AppThemeProvider>
      <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
    </AppThemeProvider>
  )
}

function customRender(
  ui: ReactElement,
  options?: Omit<RenderOptions, 'wrapper'>,
) {
  return render(ui, { wrapper: AllTheProviders, ...options })
}

export * from '@testing-library/react'
export { customRender as render }
