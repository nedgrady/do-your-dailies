import { ListItemText } from '@mui/material'
import ChoreAvatar from '../../../components/ChoreAvatar'
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
      <ChoreAvatar name={choreInQueue.choreName} emoji={choreInQueue.emoji} />
      <ListItemText
        primary={choreInQueue.choreName}
        secondary={<DuePill dueAt={choreInQueue.dueAt} now={now} />}
        slotProps={{ secondary: { component: 'div' } }}
      />
    </ClickableCardButton>
  )
}
