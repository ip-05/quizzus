import axios from 'axios';
import { defineStore } from 'pinia';
import { ref, onMounted } from 'vue';
import { useRuntimeConfig } from '#imports';

export const useAuthStore = defineStore('AuthStore', () => {
  const config = useRuntimeConfig();

  const isAuthed = ref(false);
  const token = ref(null);
  const user = ref({
    id: null,
    name: null,
    email: null,
    profilePicture: null,
  });

  onMounted(async () => {
    const localToken = localStorage.getItem('token');
    if (localToken) {
      isAuthed.value = true;
      token.value = localToken;
      await getMe();
    }
  });

  async function signInGoogle() {
    const { data } = await axios.get('/auth/google', {
      baseURL: config.public.apiUrl,
      withCredentials: true,
    });
    const { redirectUrl } = data;

    window.open(redirectUrl, '_self');
  }

  async function getMe() {
    const { data } = await axios.get('/auth/me', {
      baseURL: config.public.apiUrl,
      headers: {
        Authorization: 'Bearer ' + token.value,
      },
    });
    user.value = data;
  }

  function logout() {
    localStorage.removeItem('token');
    token.value = null;
    user.value = null;
    isAuthed.value = false;
  }

  async function authenticate(jwt) {
    isAuthed.value = true;
    token.value = jwt;
    localStorage.setItem('token', token.value);
    await getMe();
  }

  return { isAuthed, user, token, signInGoogle, authenticate, logout };
});
