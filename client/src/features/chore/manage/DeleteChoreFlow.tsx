import FlowDialog from '../../../components/flow/FlowDialog'
import type { Chore } from '../../../domain/chore'
import { useDeleteChoreFlow } from './useDeleteChoreFlow'

interface DeleteChoreFlowProps {
  chore: Chore
  onClose: () => void
}

// Mounted only while a delete is being confirmed, so the mutation state starts
// fresh each time.
export default function DeleteChoreFlow({
  chore,
  onClose,
}: DeleteChoreFlowProps) {
  const definition = useDeleteChoreFlow({ chore, onDone: onClose })
  return <FlowDialog definition={definition} onClose={onClose} />
}
