import { ListItemText } from '@mui/material'
import ClickableCardButton from '../../../components/ClickableCardButton'

type CompletionListItemProps = {
  choreName: string
  completorNames: string
  onClick: () => void
}

export default function CompletionListItem({
  choreName,
  completorNames,
  onClick,
}: CompletionListItemProps) {
  return (
    <ClickableCardButton onClick={onClick}>
      <ListItemText
        primary={choreName}
        secondary={completorNames || undefined}
      />
    </ClickableCardButton>
  )
}
