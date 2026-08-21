import DarkModeRounded from '@mui/icons-material/DarkModeRounded'
import LightModeRounded from '@mui/icons-material/LightModeRounded'
import { AppBar, IconButton, Toolbar, Tooltip, Typography } from '@mui/material'
import { Link } from '@tanstack/react-router'
import { drawerWidth } from './Sidebar'
import { ColorModeContext } from './integrations/mui/theme-provider'

export default function Header() {
  return (
    <AppBar
      color="default"
      elevation={0}
      position="fixed"
      sx={{
        width: `calc(100% - ${drawerWidth}px)`,
        ml: `${drawerWidth}px`,
      }}
    >
      <ColorModeContext.Consumer>
        {({ mode, toggleMode }) => (
          <Toolbar sx={{ borderBottom: 1, borderColor: 'divider' }}>
            <Typography
              component={Link}
              to="/today"
              sx={{
                flexGrow: 1,
                color: 'inherit',
                textDecoration: 'none',
              }}
              variant="h6"
            >
              Do Your Dailies
            </Typography>

            <Tooltip
              title={
                mode === 'dark' ? 'Switch to light mode' : 'Switch to dark mode'
              }
            >
              <IconButton
                aria-label="toggle color mode"
                color="inherit"
                onClick={toggleMode}
              >
                {mode === 'dark' ? <LightModeRounded /> : <DarkModeRounded />}
              </IconButton>
            </Tooltip>
          </Toolbar>
        )}
      </ColorModeContext.Consumer>
    </AppBar>
  )
}
