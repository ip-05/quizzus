<template>
  <div class="wrapper">
    <table class="content">
      <thead class="header">
        <td class="cell">№</td>
        <td class="cell">Points</td>
        <td class="cell">{{ tableName }}</td>
      </thead>
      <tbody v-if="mode === 'leaderboard'" class="main">
        <table-item
          v-for="({ participant, points }, id) in gameStore.leaderboardSorted"
          :key="id"
          :place="id + 1"
          :name="participant"
          :points="points"
        />
      </tbody>
      <tbody v-if="mode === 'questions'" class="main">
        <table-item
          v-for="({ name }, id) in gameStore.questions"
          :key="id"
          :place="id + 1"
          :name="name"
          :points="gameStore.points"
        />
      </tbody>
    </table>
  </div>
</template>

<script setup>
import { defineProps, computed } from 'vue';
import { useGameStore } from '../stores/game';

const gameStore = useGameStore();

const props = defineProps({
  mode: {
    type: String,
    default: 'table',
  },
});

const tableName = computed(() => {
  const word = props.mode;
  return word.charAt(0).toUpperCase() + word.slice(1);
});

// const tableData = computed(() => {
//   for (const iterator of ) {

//   }
// });
</script>

<style scoped>
.wrapper {
  background: var(--background-secondary-color);
  border-radius: 30px 30px 25px 25px;
  padding: 0 10px 10px 10px;
}

.content {
  width: 100%;
  border-collapse: collapse;
  font-family: 'Inter-SemiBold', sans-serif;
  font-size: var(--font-quaternary-size);
  line-height: 20px;
}

.header {
  color: var(--font-secondary-color);
}

.header td {
  padding: 20px;
}
</style>
