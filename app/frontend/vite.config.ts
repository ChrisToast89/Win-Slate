import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import { resolve } from 'path'

export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      '@': resolve(__dirname, 'src'),
      '@shared': resolve(__dirname, 'src/shared')
    }
  },
  // Allow importing JSON knowledge bases from frontend/data
  server: {
    port: 34115,
    strictPort: true,
    fs: { allow: [resolve(__dirname, '..'), __dirname] }
  },
  build: {
    outDir: 'dist',
    emptyOutDir: true
  }
})
