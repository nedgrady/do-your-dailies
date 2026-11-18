import FlowDialog from '../../../components/flow/FlowDialog'
import type { Chore } from '../../../domain/chore'
import { useChoreFlow } from './useChoreFlow'

interface ChoreFlowProps {
  chore: Chore
  isCreating: boolean
  onClose: () => void
}

// Mounted only while a chore is selected, keyed by chore, so each chore starts
// with fresh values and mutation state.
export default function ChoreFlow({
  chore,
  isCreating,
  onClose,
}: ChoreFlowProps) {
  const definition = useChoreFlow({ chore, isCreating, onDone: onClose })
  return <FlowDialog definition={definition} onClose={onClose} />
}
