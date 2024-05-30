import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import VueI18nPlugin from '@intlify/unplugin-vue-i18n/vite'
import path from 'path';

export default defineConfig({
  plugins: [
    vue(),
    VueI18nPlugin({
      include: [path.resolve(__dirname, './src/locales/**')],
    }),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  }
})
