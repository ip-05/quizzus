<script setup lang="ts">
import ButtonLarge from '@/components/buttons/ButtonLarge.vue';
import GameOfLife from '@/components/GameOfLife.vue';
import { useAuthStore } from '../stores/auth';

const authStore = useAuthStore();
</script>

<template>
  <GameOfLife />
  <main class="main">
    <h1 class="name">Quizzus</h1>
    <p class="slogan">
      Easy platform for creating quizzes and <br />having a blast playing them with your pals
    </p>
    <div class="buttons">
      <ButtonLarge icon="google" @click="authStore.signInGoogle">Sign in with Google</ButtonLarge>
      <ButtonLarge icon="discord" @click="authStore.signInDiscord">
        Sign in with Discord
      </ButtonLarge>
    </div>
  </main>
</template>

<style scoped>
.main {
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.name {
  margin-bottom: 20px;
  font-size: 48px;
  font-family: 'Daydream', serif;
}

.slogan {
  color: var(--color-text);
  text-align: center;
  margin-bottom: 40px;
}

.buttons {
  width: 320px;
  display: flex;
  flex-direction: column;
  gap: 10px;
}
</style>
