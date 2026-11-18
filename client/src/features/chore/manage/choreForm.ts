import { z } from 'zod'

// Form values are what the user types: the amount stays a string until submit.
export const choreFormSchema = z.object({
  name: z
    .string()
    .refine((value) => value.trim().length > 0, 'Name is required'),
  amount: z.string().refine((value) => {
    const parsed = Number(value)
    return Number.isInteger(parsed) && parsed > 0
  }, 'Must be a positive whole number'),
  unit: z.enum(['DAY', 'WEEK', 'MONTH', 'YEAR']),
})

export type ChoreFormValues = z.infer<typeof choreFormSchema>
