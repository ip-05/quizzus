<template>
  <form action="" class="form" :class="{ formHasIcon: hasIcon }" @submit.prevent="handleJoin">
    <div v-if="hasIcon" class="icon">
      <nuxt-img src="svg/icon-search.svg" alt="Search Icon" />
    </div>
    <input v-model="code" class="input" type="text" placeholder="Enter Code To Join" />
  </form>
</template>

<script setup>
import { ref, defineProps } from 'vue';
import { useGameStore } from '../stores/game';
import { navigateTo } from '#imports';

const code = ref(null);

const gameStore = useGameStore();

const { hasIcon } = defineProps({ hasIcon: Boolean });

const handleJoin = async () => {
  const gameAccessTo = await gameStore.getGame({ invite_code: code.value });
  if (gameAccessTo === 'player') return (window.location.pathname = `/game/${code.value}`);
  if (gameAccessTo === 'admin') return (window.location.pathname = `/console/${code.value}`);
  console.log('error: no such game');
};
</script>

<style scoped>
.form {
  padding: 15px 20px;
  background: var(--background-main-color);
  outline: solid 3px var(--border-color);
  outline-offset: -3px;
  border-radius: 15px;
  display: flex;
  align-items: center;
  gap: 10px;
  transition: all 250ms ease;
}

.formHasIcon {
  padding: 15px 20px 15px 15px;
}

.form:focus-within {
  outline-color: var(--font-active-color);
}

.icon {
  width: 24px;
  height: 24px;
}

.icon img {
  width: 100%;
  height: 100%;
  object-fit: contain;
}

.input {
  font-family: 'Inter-Medium', sans-serif;
  font-size: 22px;
  max-width: 200px;
  line-height: 24px;
  width: 100%;
  border: none;
  color: var(--font-active-color);
  outline: none;
}

.input::placeholder {
  color: var(--font-primary-color);
}
</style>
