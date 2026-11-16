import {
  Box,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  Drawer,
  Stack,
  Typography,
  useMediaQuery,
  useTheme,
} from '@mui/material'
import type { ReactNode } from 'react'

interface ResponsiveDialogProps {
  open: boolean
  onClose: () => void
  title: string
  actions?: ReactNode
  children: ReactNode
}

// Centered dialog on desktop, bottom sheet below the md breakpoint. Both use
// the same title / body / actions structure so callers only write it once.
export default function ResponsiveDialog({
  open,
  onClose,
  title,
  actions,
  children,
}: ResponsiveDialogProps) {
  const theme = useTheme()
  const isDesktop = useMediaQuery(theme.breakpoints.up('md'))

  if (isDesktop) {
    return (
      <Dialog open={open} onClose={onClose} fullWidth maxWidth="xs">
        <DialogTitle>{title}</DialogTitle>
        <DialogContent>{children}</DialogContent>
        {actions && (
          <DialogActions sx={{ px: 3, pb: 2 }}>{actions}</DialogActions>
        )}
      </Dialog>
    )
  }

  return (
    <Drawer
      anchor="bottom"
      open={open}
      onClose={onClose}
      slotProps={{
        paper: {
          sx: { borderTopLeftRadius: 16, borderTopRightRadius: 16 },
        },
      }}
    >
      <Box sx={{ p: 3, pb: 4 }}>
        <Typography variant="h6" component="h2" sx={{ mb: 2 }}>
          {title}
        </Typography>
        {children}
        {actions && (
          <Stack
            direction="row"
            spacing={1}
            sx={{ mt: 3, '& > *': { flex: 1 } }}
          >
            {actions}
          </Stack>
        )}
      </Box>
    </Drawer>
  )
}
