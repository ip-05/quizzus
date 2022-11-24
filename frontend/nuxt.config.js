import { defineNuxtConfig } from 'nuxt/config';

export default defineNuxtConfig({
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
  css: ['~/assets/css/main.css'],
});
