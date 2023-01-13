import { defineConfig } from 'vite';
import Vue from '@vitejs/plugin-vue';
import nuxtimg from '@nuxt/image-edge';

export default defineConfig({
  plugins: [
    Vue({
      template: {
        compilerOptions: {
          isCustomElement: (tag) => tag.startsWith('nuxt'),
        },
      },
    }),
  ],
  test: {
    globals: true,
    environment: 'jsdom',
  },
});
