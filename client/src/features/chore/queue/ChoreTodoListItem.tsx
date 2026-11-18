import { ListItemText } from '@mui/material'
import ClickableCardButton from '../../../components/ClickableCardButton'
import type { ChoreInQueue } from '../../../domain/choreQueue'
import DuePill from './DuePill'

type ChoreTodoListItemProps = {
  choreInQueue: ChoreInQueue
  now: Date
  onClick: () => void
}

export default function ChoreTodoListItem({
  choreInQueue,
  now,
  onClick,
}: ChoreTodoListItemProps) {
  return (
    <ClickableCardButton onClick={onClick}>
      <ListItemText
        primary={choreInQueue.choreName}
        secondary={<DuePill dueAt={choreInQueue.dueAt} now={now} />}
        slotProps={{ secondary: { component: 'div' } }}
      />
    </ClickableCardButton>
  )
}
