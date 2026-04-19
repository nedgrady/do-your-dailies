import { devtools } from '@tanstack/devtools-vite'
import { defineConfig } from 'vite'

import viteReact from '@vitejs/plugin-react'

const config = defineConfig({
  resolve: {
    dedupe: ['react', 'react-dom'],
  },
  plugins: [devtools(), viteReact()],
  server: {
    proxy: {
      '/api': {
        target: process.env.VITE_API_BASE_URL ?? 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
})

export default config
