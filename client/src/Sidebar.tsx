import ChecklistIcon from '@mui/icons-material/Checklist'
import TodayIcon from '@mui/icons-material/Today'
import {
  Box,
  List,
  ListItemButton,
  ListItemIcon,
  ListItemText,
} from '@mui/material'
import { Link } from '@tanstack/react-router'

const drawerWidth = 220

export default function Sidebar() {
  return (
    <Box
      component="nav"
      sx={{
        width: drawerWidth,
        flexShrink: 0,
        borderRight: 1,
        borderColor: 'divider',
        position: 'sticky',
        top: 0,
        height: '100vh',
        overflowY: 'auto',
      }}
    >
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
    </Box>
  )
}
