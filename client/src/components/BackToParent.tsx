import ArrowBackIcon from '@mui/icons-material/ArrowBack'
import Button from '@mui/material/Button'
import { Link, useMatches } from '@tanstack/react-router'
import type { FileRouteTypes } from '../routeTree.gen'

// Each page declares its parent explicitly in its route file, e.g.
//   staticData: { parent: { to: '/chores/manage', label: 'Back to chores' } }
// Pages without a parent (top-level pages) render nothing.
declare module '@tanstack/react-router' {
  interface StaticDataRouteOption {
    parent?: { to: FileRouteTypes['to']; label: string }
  }
}

export default function BackToParent() {
  const parent = useMatches({
    select: (matches) => matches.at(-1)?.staticData.parent,
  })

  if (!parent) return null

  return (
    <Button
      component={Link}
      to={parent.to}
      size="small"
      startIcon={<ArrowBackIcon />}
      sx={{ mb: 1 }}
    >
      {parent.label}
    </Button>
  )
}
