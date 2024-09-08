import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import VitePluginSitemap from 'vite-plugin-sitemap'
import { copyFileSync } from 'fs'
import { resolve } from 'path'

const routes = [
  '/',
  '/signin',
  '/signup',
  '/subscription',
  '/settings',
  '/subscribe/success',
  '/subscribe/cancel'
].map(url => ({ url, lastmod: new Date().toISOString() }));

export default defineConfig({
  plugins: [
    react(),
    VitePluginSitemap({
      hostname: 'https://watermark-generator.com',
      dynamicRoutes: routes.map(route => route.url),
      changefreq: 'weekly',
      priority: 0.8,
      outDir: './dist',
    }),
    {
      name: 'copy-robots-txt',
      closeBundle() {
        const src = resolve(__dirname, 'public', 'robots.txt')
        const dest = resolve(__dirname, 'dist', 'robots.txt')
        copyFileSync(src, dest)
      }
    }
  ],
  optimizeDeps: {
    include: ['react', 'react-dom'],
  },
  server: {
    proxy: {
      '/api': {
        target: process.env.VITE_API_URL || 'https://watermark-generator.com',
      },
    },
  },
  base: '/',
})