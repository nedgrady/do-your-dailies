import {
  Button,
  Stack,
  Typography,
  useMediaQuery,
  useTheme,
} from '@mui/material'
import { useForm } from '@tanstack/react-form'
import { useState } from 'react'
import ResponsiveDialog from '../ResponsiveDialog'
import { isFormValid, isStepValid } from './flowSteps'
import type { FlowDefinition, SetFlowValue } from './flowTypes'

interface FlowDialogProps<TValues> {
  definition: FlowDefinition<TValues>
  onClose: () => void
}

// Presents a flow definition. Desktop shows every step as a section of one
// dialog; below md the steps appear one at a time in a bottom sheet with
// Back / Next. The steps and the data they collect don't change between them.
export default function FlowDialog<TValues extends object>({
  definition,
  onClose,
}: FlowDialogProps<TValues>) {
  const theme = useTheme()
  const isDesktop = useMediaQuery(theme.breakpoints.up('md'))
  const [stepIndex, setStepIndex] = useState(0)

  const form = useForm({
    defaultValues: definition.defaultValues,
    onSubmit: ({ value }) => definition.onSubmit(value),
  })

  // The step's keys are the form's own top-level fields, which TanStack types
  // as deep keys. The cast is confined to this one bridge.
  const setValue: SetFlowValue<TValues> = (key, value) => {
    form.setFieldValue(
      key as Parameters<typeof form.setFieldValue>[0],
      value as never,
    )
  }

  return (
    <form.Subscribe selector={(state) => state.values}>
      {(values) => {
        const valid = isFormValid(definition.schema, values)
        const current = definition.steps.at(stepIndex)
        const isLastStep = stepIndex === definition.steps.length - 1
        const currentValid =
          current !== undefined &&
          isStepValid(definition.schema, values, current.fields)

        const failedNotice = definition.submitFailed && (
          <Typography color="error" role="alert" sx={{ mb: 1 }}>
            {definition.submitFailedMessage ?? "Couldn't save, try again."}
          </Typography>
        )
        const submitColor = definition.submitColor ?? 'primary'

        if (isDesktop) {
          return (
            <ResponsiveDialog
              open
              onClose={onClose}
              title={definition.title}
              actions={
                <>
                  <Button size="large" onClick={onClose}>
                    Cancel
                  </Button>
                  <Button
                    size="large"
                    variant="contained"
                    color={submitColor}
                    disabled={!valid || definition.isSubmitting}
                    onClick={() => void form.handleSubmit()}
                  >
                    {definition.submitLabel}
                  </Button>
                </>
              }
            >
              {failedNotice}
              {definition.steps.map((step) => (
                <Stack key={step.id} spacing={1} sx={{ mb: 2 }}>
                  <Typography variant="subtitle2" color="text.secondary">
                    {step.title}
                  </Typography>
                  {step.render(values, setValue)}
                </Stack>
              ))}
            </ResponsiveDialog>
          )
        }

        if (!current) return null

        return (
          <ResponsiveDialog
            open
            onClose={onClose}
            title={definition.title}
            // A single step needs no fixed height, so the sheet fits its content.
            fillHeight={definition.steps.length > 1}
            actions={
              <>
                <Button
                  size="large"
                  onClick={
                    stepIndex === 0
                      ? onClose
                      : () => setStepIndex((index) => index - 1)
                  }
                >
                  {stepIndex === 0 ? 'Cancel' : 'Back'}
                </Button>
                {isLastStep ? (
                  <Button
                    size="large"
                    variant="contained"
                    color={submitColor}
                    disabled={!valid || definition.isSubmitting}
                    onClick={() => void form.handleSubmit()}
                  >
                    {definition.submitLabel}
                  </Button>
                ) : (
                  <Button
                    size="large"
                    variant="contained"
                    disabled={!currentValid}
                    onClick={() => setStepIndex((index) => index + 1)}
                  >
                    Next
                  </Button>
                )}
              </>
            }
          >
            {failedNotice}
            <Typography
              variant="subtitle2"
              color="text.secondary"
              sx={{ mb: 1 }}
            >
              {definition.steps.length > 1
                ? `${current.title} · Step ${stepIndex + 1} of ${definition.steps.length}`
                : current.title}
            </Typography>
            {current.render(values, setValue)}
          </ResponsiveDialog>
        )
      }}
    </form.Subscribe>
  )
}
