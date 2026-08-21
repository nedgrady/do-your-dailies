import { Box } from '@mui/material'
import { TanStackDevtools } from '@tanstack/react-devtools'
import { Outlet, createRootRouteWithContext } from '@tanstack/react-router'
import { TanStackRouterDevtoolsPanel } from '@tanstack/react-router-devtools'
import Footer from '../Footer'
import AppConfigPanel from '../integrations/devtools/AppConfigPanel'
import { AppThemeProvider } from '../integrations/mui/theme-provider'
import TanStackQueryDevtools from '../integrations/tanstack-query/devtools'
import MobileNav from '../MobileNav'
import Sidebar, { drawerWidth } from '../Sidebar'

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
