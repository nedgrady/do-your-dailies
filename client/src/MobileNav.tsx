import ChecklistIcon from '@mui/icons-material/Checklist'
import DarkModeRounded from '@mui/icons-material/DarkModeRounded'
import LightModeRounded from '@mui/icons-material/LightModeRounded'
import TodayIcon from '@mui/icons-material/Today'
import { BottomNavigation, BottomNavigationAction, Paper } from '@mui/material'
import { Link, useRouterState } from '@tanstack/react-router'
import { ColorModeContext } from './integrations/mui/theme-provider'

export default function MobileNav() {
  const pathname = useRouterState({
    select: (state) => state.location.pathname,
  })

  return (
    <ColorModeContext.Consumer>
      {({ mode, toggleMode }) => (
        <Paper
          elevation={3}
          sx={{
            display: { xs: 'block', md: 'none' },
            position: 'fixed',
            bottom: 0,
            left: 0,
            right: 0,
          }}
        >
          <BottomNavigation showLabels value={pathname}>
            <BottomNavigationAction
              label="Today"
              value="/today"
              icon={<TodayIcon />}
              component={Link}
              to="/today"
            />
            <BottomNavigationAction
              label="Manage"
              value="/chores/manage"
              icon={<ChecklistIcon />}
              component={Link}
              to="/chores/manage"
            />
            <BottomNavigationAction
              label="Theme"
              value="theme"
              icon={
                mode === 'dark' ? <LightModeRounded /> : <DarkModeRounded />
              }
              onClick={toggleMode}
            />
          </BottomNavigation>
        </Paper>
      )}
    </ColorModeContext.Consumer>
  )
}
