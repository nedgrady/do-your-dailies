import AutorenewIcon from '@mui/icons-material/Autorenew'
import CheckCircleIcon from '@mui/icons-material/CheckCircle'
import ErrorIcon from '@mui/icons-material/Error'
import { Box, Stack, Typography, type SvgIconProps } from '@mui/material'
import type { UseMutationResult } from '@tanstack/react-query'
import { differenceInSeconds, formatDistance } from 'date-fns'
import { useEffect, useState, type ComponentType } from 'react'

export interface SaveStatus {
  Icon?: ComponentType<SvgIconProps>
  iconSx: SvgIconProps['sx']
  color: string
  label: string
  role: 'status' | 'alert'
  live: 'polite' | 'assertive'
  name: string
  startedAt: Date
}

export const idle = (startedAt: Date): SaveStatus =>
  Object.freeze<SaveStatus>({
    Icon: undefined,
    iconSx: undefined,
    color: 'text.disabled',
    label: 'No changes saved yet',
    role: 'status',
    live: 'polite',
    name: 'idle',
    startedAt,
  })

export const success = (startedAt: Date): SaveStatus =>
  Object.freeze<SaveStatus>({
    Icon: CheckCircleIcon,
    iconSx: undefined,
    color: 'success.main',
    label: 'Saved',
    role: 'status',
    live: 'polite',
    name: 'success',
    startedAt,
  })

export const error = (startedAt: Date): SaveStatus =>
  Object.freeze<SaveStatus>({
    Icon: ErrorIcon,
    iconSx: undefined,
    color: 'error.main',
    label: 'Save failed',
    role: 'alert',
    live: 'assertive',
    name: 'error',
    startedAt,
  })

export const pending = (startedAt: Date): SaveStatus =>
  Object.freeze<SaveStatus>({
    Icon: AutorenewIcon,
    iconSx: {
      animation: 'spin 1.5s linear infinite',
      '@keyframes spin': {
        '0%': { transform: 'rotate(0deg)' },
        '100%': { transform: 'rotate(360deg)' },
      },
    },
    color: 'info.main',
    label: 'Saving changes...',
    role: 'status',
    live: 'polite',
    name: 'pending',
    startedAt,
  })

export function mutationStatus<
  TData = unknown,
  TError = Error,
  TVariables = void,
  TOnMutateResult = unknown,
>(
  mutation: UseMutationResult<TData, TError, TVariables, TOnMutateResult>,
): SaveStatus {
  // `submittedAt` is 0 until the first `mutate()` call, so idle falls back
  // to "now" rather than the Unix epoch.
  const startedAt = new Date(mutation.submittedAt || Date.now())

  switch (mutation.status) {
    case 'error':
      return error(startedAt)
    case 'idle':
      return idle(startedAt)
    case 'success':
      return success(startedAt)
    case 'pending':
      return pending(startedAt)
  }
}

/**
 * Returns a humanized "x ago" string for `date`, re-rendering periodically
 * so the text stays fresh (e.g. "less than a minute ago" -> "1 minute ago").
 * Ticks faster while the save is recent, then backs off.
 */
function useRelativeTime(date: Date): string {
  const secondsAgo = differenceInSeconds(new Date(), date)
  const [, forceTick] = useState(0)
  useEffect(() => {
    const msSinceSave = Date.now() - date.getTime()
    // Tick every 10s for the first minute (so "just now" updates promptly),
    // then fall back to once a minute.
    const intervalMs = Math.min(msSinceSave < 60_000 ? 10_000 : 60_000, 60_000)

    const id = setInterval(() => forceTick((t) => t + 1), intervalMs)
    return () => clearInterval(id)
  }, [date])

  if (secondsAgo < 60) {
    return 'just now'
  }

  return formatDistance(date, new Date(), { addSuffix: true })
}

export function SaveStatusIndicator({ status }: { status: SaveStatus }) {
  const { Icon, iconSx, color, label, role, live, startedAt } = status

  const relativeTime = useRelativeTime(startedAt)
  const showRelativeTime = status.name === 'success'
  const accessibleLabel = showRelativeTime ? `${label} ${relativeTime}` : label

  return (
    <Stack
      direction="row"
      spacing={0.5}
      role={role}
      aria-live={live}
      aria-atomic="true"
      sx={{ position: 'relative', color, minWidth: 32, alignItems: 'center' }}
    >
      {Icon && <Icon fontSize="small" sx={iconSx} />}

      {showRelativeTime && (
        <Typography
          variant="caption"
          component="span"
          title={startedAt.toLocaleString()}
          sx={{ color: 'text.secondary', whiteSpace: 'nowrap' }}
        >
          Saved {relativeTime}
        </Typography>
      )}

      {/* Visually hidden but announced to screen readers */}
      <Box
        component="span"
        sx={{
          position: 'absolute',
          width: 1,
          height: 1,
          overflow: 'hidden',
          clip: 'rect(0 0 0 0)',
          whiteSpace: 'nowrap',
        }}
      >
        {accessibleLabel}
      </Box>
    </Stack>
  )
}
