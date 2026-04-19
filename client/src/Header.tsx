import DarkModeRounded from '@mui/icons-material/DarkModeRounded'
import LightModeRounded from '@mui/icons-material/LightModeRounded'
import { AppBar, IconButton, Toolbar, Tooltip, Typography } from '@mui/material'
import { ColorModeContext } from './integrations/mui/theme-provider'

export default function Header() {
  return (
    <AppBar color="transparent" elevation={0} position="sticky">
      <ColorModeContext.Consumer>
        {({ mode, toggleMode }) => (
          <Toolbar sx={{ borderBottom: 1, borderColor: 'divider' }}>
            <Typography sx={{ flexGrow: 1 }} variant="h6">
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
