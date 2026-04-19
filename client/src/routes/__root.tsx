import { TanStackDevtools } from '@tanstack/react-devtools'
import { Outlet, createRootRouteWithContext } from '@tanstack/react-router'
import { TanStackRouterDevtoolsPanel } from '@tanstack/react-router-devtools'
import Footer from '../Footer'
import Header from '../Header'
import AppConfigPanel from '../integrations/devtools/AppConfigPanel'
import { AppThemeProvider } from '../integrations/mui/theme-provider'
import TanStackQueryDevtools from '../integrations/tanstack-query/devtools'

const devtoolsEnabled = import.meta.env.VITE_ENABLE_DEVTOOLS === 'true'

import type { QueryClient } from '@tanstack/react-query'

interface MyRouterContext {
  queryClient: QueryClient
}

export const Route = createRootRouteWithContext<MyRouterContext>()({
  component: RootComponent,
})

function RootComponent() {
  return (
    <AppThemeProvider>
      <Header />
      <main>
        <Outlet />
      </main>
      <Footer />
      {devtoolsEnabled && (
        <TanStackDevtools
          config={{
            position: 'bottom-right',
          }}
          plugins={[
            {
              name: 'Tanstack Router',
              render: <TanStackRouterDevtoolsPanel />,
            },
            TanStackQueryDevtools,
            {
              name: 'App Config',
              render: <AppConfigPanel />,
            },
          ]}
        />
      )}
    </AppThemeProvider>
  )
}
