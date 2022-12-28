<template>
  <div class="main">
    <div class="content">
      <component :is="screen" @next="handleNextScreen" @back="handleBackScreen" />
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue';
import { useNewGameStore } from '../../stores/new';
import { resolveComponent } from '#imports';

const newGameStore = useNewGameStore();

const mode = ref('create');
const screen = computed(() => {
  if (mode.value === 'create') return CreateGame;
  if (mode.value === 'share') return ShareGame;
  return 'ErrorCreateGameComponent';
});

const CreateGame = resolveComponent('CreateGame');
const ShareGame = resolveComponent('ShareGame');

const handleNextScreen = (game) => {
  newGameStore.createGame(game);
  mode.value = 'share';
};
const handleBackScreen = () => (mode.value = 'create');
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
