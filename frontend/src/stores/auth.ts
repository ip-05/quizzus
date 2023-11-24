import { ref, onMounted } from 'vue';
import { defineStore } from 'pinia';

export const useAuthStore = defineStore('authStore', () => {
  const isAuthenticated = ref(false);
  const token = ref(null);

  onMounted(() => {
    console.log('auth store mounted');
  });

  async function signInGoogle() {}
  async function signInDiscord() {}

  function logout() {}

  return { isAuthenticated, token, signInGoogle, signInDiscord, logout };
});
