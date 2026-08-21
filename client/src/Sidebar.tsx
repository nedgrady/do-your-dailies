import ChecklistIcon from '@mui/icons-material/Checklist'
import TodayIcon from '@mui/icons-material/Today'
import {
  Drawer,
  List,
  ListItemButton,
  ListItemIcon,
  ListItemText,
  Toolbar,
} from '@mui/material'
import { Link } from '@tanstack/react-router'

export const drawerWidth = 220

export default function Sidebar() {
  return (
    <Drawer
      variant="permanent"
      sx={{
        width: drawerWidth,
        flexShrink: 0,
        '& .MuiDrawer-paper': {
          width: drawerWidth,
          boxSizing: 'border-box',
        },
      }}
    >
      <Toolbar />
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
    </Drawer>
  )
}
