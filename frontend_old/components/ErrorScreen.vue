<template>
  <div class="wrapper">
    <div class="content">
      <div class="error">{{ errorStore.message ? errorStore.message : 'Unexpected error occured' }}</div>
      <NuxtLink to="/">
        <regular-button @click="errorStore.reset">Take me home</regular-button>
      </NuxtLink>
    </div>
  </div>
</template>

<script setup>
import { useErrorStore } from '../stores/error';
const errorStore = useErrorStore();
</script>

<style scoped>
.wrapper {
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
  z-index: 100;
}

.content {
  display: flex;
  justify-content: center;
  align-items: center;
  flex-direction: column;
  gap: 20px;
}

.error {
  font-family: 'Inter-SemiBold', sans-serif;
  color: var(--font-primary-color);
  font-size: var(--font-primary-size);
}
</style>
