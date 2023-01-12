<template>
  <div class="content">
    <div class="message">Quizzus <span class="handshake">🤝</span> Google</div>
  </div>
</template>

<script setup>
import axios from 'axios';
import { onMounted } from 'vue';
import { useRoute, navigateTo, useRuntimeConfig } from '#imports';
import { useAuthStore } from '~/stores/auth';

const config = useRuntimeConfig();
const authStore = useAuthStore();

onMounted(async () => {
  if (authStore.isAuthed) return navigateTo('/');
  const { query } = useRoute();
  const googleQuery = query.state && query.code && query.scope && query.authuser && query.prompt;
  if (googleQuery) {
    const { data } = await axios.get('/auth/google/callback', {
      baseURL: config.public.apiUrl,
      withCredentials: true,
      params: query,
    });
    localStorage.setItem('token', data.token);
    authStore.authenticate(data.token);
    navigateTo('/');
  }
});
</script>

<style scoped>
.content {
  width: 100%;
  height: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
  flex-direction: column;
  position: fixed;
  top: 0;
  bottom: 0;
  left: 0;
  right: 0;
}

.content::after {
  content: '';
  position: absolute;
  top: 0;
  bottom: 0;
  left: 0;
  right: 0;
  background: white;
  z-index: -1;
}

.message {
  font-size: var(--font-primary-size);
  color: var(--font-primary-color);
  font-family: 'Inter-Bold', sans-serif;
  margin-bottom: 20px;
}

.handshake {
  display: inline-block;
  animation: handshake 500ms ease-in-out infinite;
}

@keyframes handshake {
  0% {
    transform: translateY(0px);
  }
  50% {
    transform: translateY(2.5px);
  }
  to {
    transform: translateY(0px);
  }
}
</style>
