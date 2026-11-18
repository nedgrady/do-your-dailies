import type {
  FlowDefinition,
  FlowStep,
} from '../../../components/flow/flowTypes'
import type { Chore } from '../../../domain/chore'
import {
  useCreateChoreMutation,
  useUpdateChoreMutation,
} from '../../../domain/chore'
import { toAmount, toDays } from '../../Cadence'
import { CadenceStep, NameStep } from './ChoreFieldSteps'
import type { ChoreFormValues } from './choreForm'
import { choreFormSchema } from './choreForm'

interface UseChoreFlowOptions {
  chore: Chore
  isCreating: boolean
  onDone: () => void
}

// Describes "add or edit a chore" as a flow: name, then cadence (amount and
// unit as one value). Delete is not part of this flow yet.
export function useChoreFlow({
  chore,
  isCreating,
  onDone,
}: UseChoreFlowOptions): FlowDefinition<ChoreFormValues> {
  const createMutation = useCreateChoreMutation()
  const updateMutation = useUpdateChoreMutation()
  const saveMutation = isCreating ? createMutation : updateMutation

  const steps: FlowStep<ChoreFormValues>[] = [
    {
      id: 'name',
      title: 'Name',
      fields: ['name'],
      render: (values, setValue) => (
        <NameStep values={values} setValue={setValue} />
      ),
    },
    {
      id: 'cadence',
      title: 'How often',
      fields: ['amount', 'unit'],
      render: (values, setValue) => (
        <CadenceStep values={values} setValue={setValue} />
      ),
    },
  ]

  return {
    title: isCreating ? 'Add chore' : 'Edit chore',
    schema: choreFormSchema,
    defaultValues: isCreating
      ? { name: '', amount: '1', unit: 'DAY' }
      : {
          name: chore.name,
          amount: String(toAmount(chore.cadenceInDays, chore.displayUnit)),
          unit: chore.displayUnit,
        },
    steps,
    submitLabel: 'Save',
    onSubmit: (values) => {
      const data = {
        name: values.name.trim(),
        cadenceInDays: toDays(Number(values.amount), values.unit),
        displayUnit: values.unit,
      }
      if (isCreating) {
        createMutation.mutate({ data }, { onSuccess: onDone })
      } else {
        updateMutation.mutate({ id: chore.id, data }, { onSuccess: onDone })
      }
    },
    isSubmitting: saveMutation.isPending,
    submitFailed: saveMutation.isError,
  }
}
