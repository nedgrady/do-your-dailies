import { devtools } from '@tanstack/devtools-vite'
import { defineConfig } from 'vite'

import viteReact from '@vitejs/plugin-react'

const config = defineConfig({
  resolve: {
    dedupe: ['react', 'react-dom'],
  },
  plugins: [devtools(), viteReact()],
})

export default config
