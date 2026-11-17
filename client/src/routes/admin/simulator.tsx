import { createFileRoute } from '@tanstack/react-router'
import SimulatorPage from '../../features/admin/simulator/SimulatorPage'

export const Route = createFileRoute('/admin/simulator')({
  component: SimulatorPage,
  staticData: { parent: { to: '/chores/manage', label: 'Back to chores' } },
})
