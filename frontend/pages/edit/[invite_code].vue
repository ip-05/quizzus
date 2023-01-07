<template>
  <div class="main">
    <div class="content">
      <component :is="screen" />
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted } from 'vue';
import { useNewGameStore } from '../../stores/new';
import { useErrorStore } from '../../stores/error';
import { resolveComponent, useRoute } from '#imports';

const newGameStore = useNewGameStore();
const errorStore = useErrorStore();
const { params } = useRoute();

onMounted(() => {
  newGameStore.getGame(params);
});

const mode = computed(() => {
  if (errorStore.message) {
    return 'error';
  }
  return 'create';
});
const screen = computed(() => {
  if (mode.value === 'create') return CreateGame;
  return ErrorScreen;
});

const CreateGame = resolveComponent('CreateGame');
const ErrorScreen = resolveComponent('ErrorScreen');
</script>

<style scoped>
.main {
  margin-top: 100px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.content {
  max-width: 1080px;
  width: calc(100% - 20px);
  height: auto;
}
</style>
