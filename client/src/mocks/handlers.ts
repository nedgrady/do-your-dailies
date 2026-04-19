import { http, HttpResponse } from 'msw'

export const handlers = [
  http.get('/api/chore-queue', () => {
    return HttpResponse.json([])
  }),
  http.all('/api/{*path}', ({ request, params }) => {
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
