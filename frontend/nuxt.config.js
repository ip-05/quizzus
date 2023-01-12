import { defineNuxtConfig } from 'nuxt/config';

export default defineNuxtConfig({
  modules: ['@pinia/nuxt', '@nuxt/image-edge'],
  runtimeConfig: {
    public: {
      apiUrl: process.env.API_URL || 'http://localhost:3001',
      socketUrl: process.env.WS_URL || 'ws://localhost:3001/ws',
    },
  },
  app: {
    meta: [
      {
        name: 'viewport',
        content: 'width=device-width, initial-scale=1',
      },
      {
        charset: 'utf-8',
      },
    ],
  },
  plugins: [{ src: '~/plugins/confetti.js', mode: 'client' }],
  css: ['~/assets/css/main.css'],
});
