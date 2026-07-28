// domain/dateCoercion.ts
import { z } from 'zod'

export const isoDateTime = z.iso
  .datetime({ offset: true })
  .pipe(z.coerce.date())
