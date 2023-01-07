// import regeneratorRuntime from 'regenerator-runtime';
import ConfettiExplosion from 'vue-confetti-explosion';
import { defineNuxtPlugin } from '#imports';

export default defineNuxtPlugin((nuxtApp) => {
  nuxtApp.vueApp.use(ConfettiExplosion);
});
