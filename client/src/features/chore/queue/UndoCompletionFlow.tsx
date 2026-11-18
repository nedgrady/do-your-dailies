import FlowDialog from '../../../components/flow/FlowDialog'
import { useUndoCompletionFlow } from './useUndoCompletionFlow'

interface UndoCompletionFlowProps {
  completionId: number
  choreName: string
  completorNames: string
  completedAt: Date
  onClose: () => void
}

// Mounted only while a completion is being viewed, so the undo mutation's
// state starts fresh each time.
export default function UndoCompletionFlow({
  completionId,
  choreName,
  completorNames,
  completedAt,
  onClose,
}: UndoCompletionFlowProps) {
  const definition = useUndoCompletionFlow({
    completionId,
    choreName,
    completorNames,
    completedAt,
    onDone: onClose,
  })
  return <FlowDialog definition={definition} onClose={onClose} />
}
