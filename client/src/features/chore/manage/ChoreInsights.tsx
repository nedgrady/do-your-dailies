import { Box, Typography } from '@mui/material'
import { orderBy } from 'es-toolkit'
import type { Chore } from '../../../domain/chore'
import { useChoreInsightsQuery } from '../../../domain/insights'
import { formatCadence } from '../../Cadence'

export function ChoreInsights() {
  const { data } = useChoreInsightsQuery()

  if (!data) return <>Loading...</>

  const {
    minCapacityToKeepUtilizationRatioGreaterThanOne,
    userDesiredCapacity,
    utilizationRatio,
    choreProjections,
  } = data

  const willBeSippage = utilizationRatio > 1

  const sorted = orderBy(
    choreProjections,
    [
      (proctedChore) => proctedChore.chore.cadenceInDays,
      (proctedChore) => proctedChore.chore.createdAt,
    ],
    ['asc', 'asc'],
  )

  const first = sorted[0]
  const last = sorted[sorted.length - 1]
  const middle = sorted[Math.floor((sorted.length - 1) / 2)]

  // Dedupe by chore id, preserving order: first -> middle -> last
  const picked = [first, middle, last]
  const seen = new Set<Chore['id']>()
  const deduped = picked.filter((p) => {
    if (!p) return false
    if (seen.has(p.chore.id)) return false
    seen.add(p.chore.id)
    return true
  })

  return (
    <Box>
      <Typography>
        You want to do {userDesiredCapacity} chores per day
      </Typography>
      <Typography>
        Your chore list would require{' '}
        {minCapacityToKeepUtilizationRatioGreaterThanOne.toFixed(1)} chores per
        day
      </Typography>
      {willBeSippage && (
        <>
          <Typography>This means there will be slippage.</Typography>
          <Typography>Example:</Typography>
          {deduped.map((choreProjection) => (
            <Typography key={choreProjection.chore.id}>
              You wanted to complete "{choreProjection.chore.name}"
              {formatCadence(choreProjection.chore.cadenceInDays)}, but you will
              only be able to complete it{' '}
              {formatCadence(choreProjection.projectedCadence)}.
            </Typography>
          ))}
        </>
      )}
      {!willBeSippage && (
        <Typography>
          This means you will complete all your chores in good time.
        </Typography>
      )}
    </Box>
  )
}
