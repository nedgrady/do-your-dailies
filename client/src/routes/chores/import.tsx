import { createFileRoute } from '@tanstack/react-router'
import ImportChoresPage from '../../features/chore/import/ImportChoresPage'

export const Route = createFileRoute('/chores/import')({
  component: ImportChoresPage,
  staticData: { parent: { to: '/chores/manage', label: 'Back to chores' } },
})
