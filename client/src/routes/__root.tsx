import { Box } from '@mui/material'
import { TanStackDevtools } from '@tanstack/react-devtools'
import type { QueryClient } from '@tanstack/react-query'
import { Outlet, createRootRouteWithContext } from '@tanstack/react-router'
import { TanStackRouterDevtoolsPanel } from '@tanstack/react-router-devtools'
import BackToParent from '../components/BackToParent'
import Footer from '../Footer'
import AppConfigPanel from '../integrations/devtools/AppConfigPanel'
import { AppThemeProvider } from '../integrations/mui/theme-provider'
import TanStackQueryDevtools from '../integrations/tanstack-query/devtools'
import MobileNav from '../MobileNav'
import Sidebar, { drawerWidth } from '../Sidebar'

const devtoolsEnabled = import.meta.env.VITE_ENABLE_DEVTOOLS === 'true'

interface MyRouterContext {
  queryClient: QueryClient
}

// Search params declared here are shared by every route beneath the root.
// `completing` is the chore id whose "Who did it?" sheet is open. A number,
// not a numeric string, for the same serialization reason as chores/manage.
// `completion` is the completion id whose "Mark not done?" sheet is open.
// `picking` is true while the "Add completed chore" flow is open.
export interface RootSearch {
  completing?: number
  completion?: number
  picking?: true
}

export const Route = createRootRouteWithContext<MyRouterContext>()({
  component: RootComponent,
  // Each param is validated on its own so the schema never drops a valid one.
  validateSearch: (search: Record<string, unknown>): RootSearch => {
    const result: RootSearch = {}
    if (typeof search.completing === 'number') {
      result.completing = search.completing
    }
    if (typeof search.completion === 'number') {
      result.completion = search.completion
    }
    if (search.picking === true) {
      result.picking = true
    }
    return result
  },
})

function RootComponent() {
  return (
    <AppThemeProvider>
      <Box sx={{ display: 'flex' }}>
        <Sidebar />
        <Box
          component="main"
          sx={{
            flexGrow: 1,
            width: { md: `calc(100% - ${drawerWidth}px)` },
            pb: { xs: 8, md: 0 },
          }}
        >
          <Box sx={{ mt: 2, mx: 'auto', maxWidth: 720, px: 2 }}>
            <BackToParent />
            <Outlet />
          </Box>
          <Footer />
        </Box>
      </Box>
      <MobileNav />
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
