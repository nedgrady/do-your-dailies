import { defineConfig } from 'orval'

export default defineConfig({
  choreApiZod: {
    input: '../server/openapi/openapi.yaml',
    output: {
      mode: 'tags-split',
      target: 'src/generated/zod/api.zod.ts',
      client: 'zod',
    },
  },
  choreApi: {
    input: '../server/openapi/openapi.yaml',
    output: {
      mode: 'tags-split',
      target: 'src/generated/api/api.ts',
      client: 'react-query',
      httpClient: 'axios',
      override: {
        query: {
          useSuspenseQuery: true,
        },
        mutator: {
          path: './src/lib/orval-axios-mutator.ts',
          name: 'customInstance',
        },
      },
    },
  },
})
