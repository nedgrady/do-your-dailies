import { ListItemText } from '@mui/material'
import ChoreAvatar from '../../../components/ChoreAvatar'
import ClickableCardButton from '../../../components/ClickableCardButton'
import type { Emoji } from '../../../domain/emoji'

type CompletionListItemProps = {
  choreName: string
  emoji: Emoji | null
  completorNames: string
  onClick: () => void
}

export default function CompletionListItem({
  choreName,
  emoji,
  completorNames,
  onClick,
}: CompletionListItemProps) {
  return (
    <ClickableCardButton onClick={onClick}>
      <ChoreAvatar name={choreName} emoji={emoji} />
      <ListItemText
        primary={choreName}
        secondary={completorNames || undefined}
      />
    </ClickableCardButton>
  )
}
