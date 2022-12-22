import axios from 'axios';
import { defineStore } from 'pinia';
import { ref, onMounted } from 'vue';
import { useRuntimeConfig } from '#imports';

export const useAuthStore = defineStore('AuthStore', () => {
  const config = useRuntimeConfig();

  const isAuthed = ref(false);
  const token = ref(null);
  const user = ref(null);

  onMounted(async () => {
    const localToken = localStorage.getItem('token');
    if (localToken) {
      isAuthed.value = true;
      token.value = localToken;
      await getMe();
      console.log(user.value);
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

  async function authenticate(jwt) {
    isAuthed.value = true;
    token.value = jwt;
    localStorage.setItem('token', token.value);
    await getMe();
  }

  return { isAuthed, token, signInGoogle, authenticate };
});
