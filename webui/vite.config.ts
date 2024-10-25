import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'


const githubBase = '/chipper/';
const defaultConfig = {
  plugins: [
    vue(),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  }
};

export default defineConfig(({ mode }) => {
  let base = '/';
  if (mode === 'production') {
    base = githubBase;
  }

  return Object.assign({ base }, defaultConfig)
})
