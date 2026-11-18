import ScheduleIcon from '@mui/icons-material/Schedule'
import { Chip } from '@mui/material'
import ClockExclamationIcon from '../../../components/icons/ClockExclamationIcon'
import { describeDue } from '../../../domain/choreDue'

type DuePillProps = {
  dueAt: Date
  now: Date
}

// Neutral on purpose: overdue changes the icon and label, not the colour.
export default function DuePill({ dueAt, now }: DuePillProps) {
  const { isOverdue, label } = describeDue(dueAt, now)
  return (
    <Chip
      variant="outlined"
      sx={{
        bgcolor: 'transparent',
        color: 'text.secondary',
        borderColor: 'divider',
        '& .MuiChip-icon': { color: 'inherit', ml: 0.5, mr: 0 },
        '& .MuiChip-label': { pl: 0.5, pr: 1 },
      }}
      icon={
        isOverdue ? (
          <ClockExclamationIcon fontSize="small" />
        ) : (
          <ScheduleIcon fontSize="small" />
        )
      }
      label={label}
    />
  )
}
