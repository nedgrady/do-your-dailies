import type { ReactNode } from 'react'
import type { z } from 'zod'

export type SetFlowValue<TValues> = <TKey extends keyof TValues>(
  key: TKey,
  value: TValues[TKey],
) => void

// One screen of a flow. `fields` lists the values the step collects; the step
// is complete when the schema has no issues on any of them.
export interface FlowStep<TValues> {
  id: string
  title: string
  fields: ReadonlyArray<keyof TValues & string>
  render: (values: TValues, setValue: SetFlowValue<TValues>) => ReactNode
}

// What a flow collects and what it does with the result. The presenter decides
// how to show it (one dialog on desktop, stepped bottom sheet on mobile).
export interface FlowDefinition<TValues> {
  title: string
  schema: z.ZodType<TValues>
  defaultValues: TValues
  steps: ReadonlyArray<FlowStep<TValues>>
  submitLabel: string
  // 'error' for destructive submits such as "Mark not done".
  submitColor?: 'primary' | 'error'
  submitFailedMessage?: string
  // Desktop only: show the steps side by side in a wider dialog, instead of
  // stacked. Use when the steps are both lists the user works through at once.
  desktopColumns?: boolean
  onSubmit: (values: TValues) => void
  isSubmitting: boolean
  submitFailed: boolean
}
