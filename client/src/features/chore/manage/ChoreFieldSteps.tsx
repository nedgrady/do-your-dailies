import {
  Stack,
  TextField,
  ToggleButton,
  ToggleButtonGroup,
} from '@mui/material'
import type { SetFlowValue } from '../../../components/flow/flowTypes'
import type { DisplayUnit } from '../../../generated/api/api.schemas'
import { DISPLAY_UNITS, unitLabel } from '../../Cadence'
import type { ChoreFormValues } from './choreForm'

// Each step is one logical value. Cadence (amount + unit) is one value, so it
// is one step.

interface StepProps {
  values: ChoreFormValues
  setValue: SetFlowValue<ChoreFormValues>
}

export function NameStep({ values, setValue }: StepProps) {
  return (
    <TextField
      label="Chore name"
      fullWidth
      value={values.name}
      onChange={(event) => setValue('name', event.target.value)}
    />
  )
}

// "Every [amount] [unit]" as one value.
export function CadenceStep({ values, setValue }: StepProps) {
  return (
    <Stack spacing={2}>
      <TextField
        label="Every"
        type="number"
        fullWidth
        helperText="Whole number, 1 or more"
        value={values.amount}
        onChange={(event) => setValue('amount', event.target.value)}
        slotProps={{ htmlInput: { min: 1, step: 1 } }}
      />
      <ToggleButtonGroup
        value={values.unit}
        exclusive
        onChange={(_event, unit: DisplayUnit | null) => {
          if (unit) setValue('unit', unit)
        }}
        aria-label="cadence unit"
        fullWidth
      >
        {DISPLAY_UNITS.map((unit) => (
          <ToggleButton key={unit} value={unit} aria-label={unitLabel(unit)}>
            {unitLabel(unit)}
          </ToggleButton>
        ))}
      </ToggleButtonGroup>
    </Stack>
  )
}
