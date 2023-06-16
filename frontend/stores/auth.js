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

  (async () => {
    // if token in local storage exists
    console.log(localStorage.getItem('token'));
    if (localStorage.getItem('token')) {
      isAuthed.value = true;
      token.value = localStorage.getItem('token');
      await getMe();
    }
  })();

  async function signInGoogle() {
    const { data } = await axios.get('/auth/google', {
      baseURL: config.public.apiUrl,
      withCredentials: true,
    });
    const { redirectUrl } = data;

    window.open(redirectUrl, '_self');
  }

  async function getMe() {
    console.log(config.public.apiUrl);
    const { data } = await axios.get('/users/me', {
      baseURL: config.public.apiUrl,
      headers: {
        Authorization: 'Bearer ' + token.value,
      },
    });
    user.value = data;
  }

  function logout() {
    // tokenCookie.value = null; // reset token cookie
    localStorage.removeItem('token');
    token.value = null;
    user.value = null;
    isAuthed.value = false;
  }

  async function authenticate(jwt) {
    isAuthed.value = true;
    token.value = jwt;
    localStorage.setItem('token', jwt);
    // tokenCookie.value = token.value; // set token cookie
    await getMe();
  }

  return { isAuthed, user, token, signInGoogle, authenticate, logout };
});
