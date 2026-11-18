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
  // Mobile only: fixed sheet height so it doesn't resize between steps.
  fillHeight?: boolean
  // Desktop only: a wider dialog, for content laid out in columns.
  wide?: boolean
}

// Centered dialog on desktop, bottom sheet below the md breakpoint.
//
// Layout contract:
// - The title and actions always stay visible and only the body scrolls.
//   Callers supply title, children and actions, and must not size or scroll
//   inside children, or the actions can be pushed off screen.
// - Don't change the height of content while the user interacts (filtering,
//   toggling a row, async lists). Reserve the space instead, or the dialog
//   jumps. Desktop is top-anchored, so residual growth only moves the bottom.
export default function ResponsiveDialog({
  open,
  onClose,
  title,
  actions,
  children,
  fillHeight = false,
  wide = false,
}: ResponsiveDialogProps) {
  const theme = useTheme()
  const isDesktop = useMediaQuery(theme.breakpoints.up('md'))

  if (isDesktop) {
    return (
      <Dialog
        open={open}
        onClose={onClose}
        fullWidth
        maxWidth={wide ? 'md' : 'sm'}
        scroll="paper"
        // Anchored to the top, not centred, so a height change only moves the
        // bottom edge. The title and first lines of content stay put.
        sx={{ '& .MuiDialog-container': { alignItems: 'flex-start' } }}
        slotProps={{
          paper: {
            sx: { mt: '10vh', maxHeight: 'calc(100% - 10vh - 32px)' },
          },
        }}
      >
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
          sx: {
            borderTopLeftRadius: 16,
            borderTopRightRadius: 16,
            display: 'flex',
            flexDirection: 'column',
            // dvh follows the on-screen keyboard (see interactive-widget in index.html).
            maxHeight: '90dvh',
            height: fillHeight ? '70dvh' : undefined,
          },
        },
      }}
    >
      <Box sx={{ px: 3, pt: 3, pb: 2, flexShrink: 0 }}>
        <Typography variant="h6" component="h2">
          {title}
        </Typography>
      </Box>
      <Box sx={{ px: 3, flex: '1 1 auto', minHeight: 0, overflowY: 'auto' }}>
        {children}
      </Box>
      {actions && (
        <Stack
          direction="row"
          spacing={1}
          sx={{
            flexShrink: 0,
            px: 3,
            pt: 2,
            pb: 'calc(16px + env(safe-area-inset-bottom))',
            borderTop: 1,
            borderColor: 'divider',
            '& > *': { flex: 1 },
          }}
        >
          {actions}
        </Stack>
      )}
    </Drawer>
  )
}
