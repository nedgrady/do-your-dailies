import type {
  FlowDefinition,
  FlowStep,
} from '../../../components/flow/flowTypes'
import { useAllChoresQuery } from '../../../domain/chore'
import {
  useChoreCompletionsBetweenQuery,
  useCreateChoreCompletionMutation,
} from '../../../domain/choreCompletion'
import {
  useChoreCompletorsQuery,
  useCreateChoreCompletorMutation,
} from '../../../domain/choreCompletor'
import { useCurrentDayChoreQueueQuery } from '../../../domain/choreQueue'
import { CreateChoreCompletionBody } from '../../../generated/zod/default/default'
import type { z } from 'zod'
import ChoreStep from './ChoreStep'
import PeopleStep from './PeopleStep'

export type CompleteChoreValues = z.infer<typeof CreateChoreCompletionBody>

interface UseCompleteChoreFlowOptions {
  // Set when the chore is already known (tapped from the queue). That skips
  // the chore step.
  presetChoreId?: number
  onDone: () => void
}

// Describes "record a chore completion" as a flow: which values it collects,
// which steps show them, and what happens on submit. Presentation is left to
// FlowDialog.
export function useCompleteChoreFlow({
  presetChoreId,
  onDone,
}: UseCompleteChoreFlowOptions): FlowDefinition<CompleteChoreValues> {
  const { data: allChores = [] } = useAllChoresQuery()
  const { data: choresInQueue = [] } = useCurrentDayChoreQueueQuery()
  const { data: completions = [] } = useChoreCompletionsBetweenQuery(
    new Date(new Date().setHours(0, 0, 0, 0)),
    new Date(new Date().setHours(24, 0, 0, 0)),
  )
  const { data: completors = [] } = useChoreCompletorsQuery()
  const completeChore = useCreateChoreCompletionMutation()
  const addPerson = useCreateChoreCompletorMutation()

  // Ad-hoc picks exclude chores already in the queue (those are tapped on the
  // To do list) and chores already completed today.
  const queuedIds = new Set(choresInQueue.map((chore) => chore.choreId))
  const doneTodayIds = new Set(
    completions.map((completion) => completion.choreId),
  )
  const pickableChores = allChores.filter(
    (chore) => !queuedIds.has(chore.id) && !doneTodayIds.has(chore.id),
  )

  const presetChore = allChores.find((chore) => chore.id === presetChoreId)

  const handleAddPerson = async (name: string) => {
    const created = await addPerson.mutateAsync({ data: { name } })
    return created.id
  }

  const choreStep: FlowStep<CompleteChoreValues> = {
    id: 'chore',
    title: 'Which chore?',
    fields: ['choreId'],
    render: (values, setValue) => (
      <ChoreStep
        chores={pickableChores}
        selectedId={values.choreId}
        onSelect={(choreId) => setValue('choreId', choreId)}
      />
    ),
  }

  const peopleStep: FlowStep<CompleteChoreValues> = {
    id: 'people',
    title: 'Who did it?',
    fields: ['choreCompletorIds'],
    render: (values, setValue) => (
      <PeopleStep
        completors={completors}
        selectedIds={values.choreCompletorIds}
        onChange={(ids) => setValue('choreCompletorIds', ids)}
        onAddPerson={handleAddPerson}
        isAddingPerson={addPerson.isPending}
        addPersonFailed={addPerson.isError}
      />
    ),
  }

  return {
    title: presetChore
      ? `Who did "${presetChore.name}"?`
      : 'Add completed chore',
    schema: CreateChoreCompletionBody,
    defaultValues: {
      choreId: presetChoreId ?? 0,
      choreCompletorIds: [],
    },
    steps: presetChoreId === undefined ? [choreStep, peopleStep] : [peopleStep],
    submitLabel: 'Done',
    onSubmit: (values) => {
      completeChore.mutate({ data: values }, { onSuccess: onDone })
    },
    isSubmitting: completeChore.isPending,
    submitFailed: completeChore.isError,
  }
}
