import type { z } from 'zod'

// True when the schema has no issues on any of the given fields. Issues on
// other fields (not yet visited) don't block this step.
export function isStepValid<TValues>(
  schema: z.ZodType<TValues>,
  values: TValues,
  fields: ReadonlyArray<string>,
): boolean {
  const result = schema.safeParse(values)
  if (result.success) return true
  return !result.error.issues.some((issue) =>
    fields.includes(String(issue.path[0])),
  )
}

export function isFormValid<TValues>(
  schema: z.ZodType<TValues>,
  values: TValues,
): boolean {
  return schema.safeParse(values).success
}
