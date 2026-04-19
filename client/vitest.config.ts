import { defineConfig } from 'vitest/config'

export default defineConfig({
  test: {
    setupFiles: ['./src/test-setup.ts'],
    env: {
      VITE_API_BASE_URL: 'http://localhost:8080',
    },
  },
})
