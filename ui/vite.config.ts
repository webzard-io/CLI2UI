import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  define: {
    // react-codemirror2 needs this
    global: "globalThis",
    // https://github.com/vitejs/vite/issues/1973
    'process.env': {}
  }
})
