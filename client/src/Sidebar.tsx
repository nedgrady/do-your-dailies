import ChecklistIcon from '@mui/icons-material/Checklist'
import DarkModeRounded from '@mui/icons-material/DarkModeRounded'
import LightModeRounded from '@mui/icons-material/LightModeRounded'
import TodayIcon from '@mui/icons-material/Today'
import {
  Box,
  Drawer,
  IconButton,
  List,
  ListItemButton,
  ListItemIcon,
  ListItemText,
  Toolbar,
  Tooltip,
  Typography,
} from '@mui/material'
import { Link } from '@tanstack/react-router'
import { ColorModeContext } from './integrations/mui/theme-provider'

export const drawerWidth = 220

export default function Sidebar() {
  return (
    <Drawer
      variant="permanent"
      sx={{
        display: { xs: 'none', md: 'block' },
        width: drawerWidth,
        flexShrink: 0,
        '& .MuiDrawer-paper': {
          width: drawerWidth,
          boxSizing: 'border-box',
        },
      }}
    >
      <Toolbar>
        <Typography
          component={Link}
          to="/today"
          variant="h6"
          sx={{ color: 'inherit', textDecoration: 'none' }}
        >
          Do Your Dailies
        </Typography>
      </Toolbar>
      <List>
        <ListItemButton
          component={Link}
          to="/today"
          sx={{
            '&[data-status="active"]': {
              backgroundColor: 'action.selected',
              fontWeight: 700,
            },
          }}
        >
          <ListItemIcon>
            <TodayIcon />
          </ListItemIcon>
          <ListItemText primary="Today" />
        </ListItemButton>
        <ListItemButton
          component={Link}
          to="/chores/manage"
          sx={{
            '&[data-status="active"]': {
              backgroundColor: 'action.selected',
              fontWeight: 700,
            },
          }}
        >
          <ListItemIcon>
            <ChecklistIcon />
          </ListItemIcon>
          <ListItemText primary="Manage chores" />
        </ListItemButton>
      </List>

      <ColorModeContext.Consumer>
        {({ mode, toggleMode }) => (
          <Box sx={{ mt: 'auto', p: 1 }}>
            <Tooltip
              title={
                mode === 'dark' ? 'Switch to light mode' : 'Switch to dark mode'
              }
            >
              <IconButton aria-label="toggle color mode" onClick={toggleMode}>
                {mode === 'dark' ? <LightModeRounded /> : <DarkModeRounded />}
              </IconButton>
            </Tooltip>
          </Box>
        )}
      </ColorModeContext.Consumer>
    </Drawer>
  )
}
