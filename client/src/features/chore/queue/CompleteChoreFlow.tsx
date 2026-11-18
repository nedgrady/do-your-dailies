import FlowDialog from '../../../components/flow/FlowDialog'
import { useCompleteChoreFlow } from './useCompleteChoreFlow'

interface CompleteChoreFlowProps {
  presetChoreId?: number
  onClose: () => void
}

// Mounted only while the flow is open, so its mutation and step state start
// fresh each time.
export default function CompleteChoreFlow({
  presetChoreId,
  onClose,
}: CompleteChoreFlowProps) {
  const definition = useCompleteChoreFlow({ presetChoreId, onDone: onClose })
  return <FlowDialog definition={definition} onClose={onClose} />
}
