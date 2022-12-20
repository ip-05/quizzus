<template>
  <button class="button" @click="authStore.signInGoogle">
    <div class="icon">
      <nuxt-img src="svg/icon-google.svg" alt="Sign In with Google" />
    </div>
    <span class="text">Sign In with Google</span>
  </button>
</template>

<script setup>
import { useAuthStore } from '../stores/auth';

const authStore = useAuthStore();
</script>

<style scoped>
.button {
  padding: 15px 20px 15px 15px;
  background: var(--background-main-color);
  border: none;
  outline: solid 3px var(--border-color);
  outline-offset: -3px;
  border-radius: 15px;
  font-family: 'Inter-Medium', sans-serif;
  font-size: 22px;
  display: flex;
  align-items: center;
  gap: 10px;
  color: var(--font-primary-color);
  cursor: pointer;
  transition: all 250ms ease;
}

/* .button:hover {
  box-shadow: 0px 8px 32px 0px rgba(182, 237, 205, 0.2);
} */

.icon {
  width: 24px;
  height: 24px;
}

.icon img {
  width: 100%;
  height: 100%;
  object-fit: contain;
}

.text {
  text-decoration: none;
  line-height: 24px;
}
</style>
