import { apiBaseUrl } from '#/lib/axios'
import { http, HttpResponse } from 'msw'

export const handlers = [
  http.get(`${apiBaseUrl}/api/chore-queue`, () => {
    return HttpResponse.json([])
  }),
  http.all(`${apiBaseUrl}/api/{*path}`, ({ request, params }) => {
    console.warn(
      `[MSW] Unmocked API request: ${request.method} /api/${String(params.path ?? '')}`,
    )

    return HttpResponse.json(
      {
        error: 'No MSW handler configured for this endpoint.',
      },
      { status: 501 },
    )
  }),
]
