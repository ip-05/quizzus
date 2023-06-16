<template>
  <div class="console__wrapper">
    <div class="content">
      <div class="main">
        <game-display class="display" />
        <Transition name="leaderboard">
          <game-table v-if="gameStore.leaderboardShown" mode="leaderboard" class="leaderboard">Leaderboard</game-table>
        </Transition>
        <game-table mode="questions" class="question" />
      </div>
      <div class="sidebar">
        <div class="code__wrapper">
          <span>Game code</span>
          <div class="code" @click="saveToClipboard">
            <span>{{ gameStore.inviteCode }}</span>
            <nuxt-img src="svg/icon-copy.svg" />
          </div>
        </div>
        <div class="button__wrapper">
          <NuxtLink :to="gameStore.gameStarted ? false : `/edit/${gameStore.inviteCode}`">
            <regular-button :disabled="gameStore.gameStarted" style="width: 100%">Edit</regular-button>
          </NuxtLink>
          <regular-button active :disabled="disabled" @click="handleGame">{{
            gameStore.gameStarted ? 'Next Question' : 'Start Game'
          }}</regular-button>
          <span>After starting game you can not edit it</span>
        </div>
        <div class="separator"></div>
        <div class="participants__wrapper">
          <div class="participants__header">
            <nuxt-img src="svg/icon-participants-filled.svg" />
            <span>Participants</span>
          </div>
          <div class="participants">
            <nuxt-img v-for="(participant, id) in gameStore.participants" :key="id" :src="participant.profilePicture" />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted, computed } from 'vue';
import { useGameStore } from '../../stores/game';
import { useSocketStore } from '../../stores/socket';
import { useRoute } from '#imports';

const { params } = useRoute();
const gameStore = useGameStore();

const disabled = computed(() => {
  if (!gameStore.gameStarted) return false;
  if (gameStore.canContinue) return false;
  return true;
});

onMounted(() => {
  gameStore.inviteCode = params.invite_code;
  gameStore.getGame(params);
  const socketStore = useSocketStore();
  socketStore.joinGame(params.invite_code);
  // gameStore.joinGame(params.invite_code);
});

const handleGame = () => {
  if (!gameStore.gameStarted) {
    gameStore.startGame();
    return;
  }
  if (gameStore.canContinue) {
    gameStore.nextQuestion();
  }
};

const saveToClipboard = () => {
  const copyText = gameStore.inviteCode;
  navigator.clipboard.writeText(copyText);

  alert('Copied the text: ' + copyText);
};
</script>

<style scoped>
.console__wrapper {
  margin-top: 100px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  margin-bottom: 20px;
}

.content {
  max-width: 1080px;
  width: calc(100% - 20px);
  height: auto;
  display: grid;
  grid-template-columns: 9fr 3fr;
  gap: 20px;
}

.display,
.leaderboard,
.quesions {
  margin-bottom: 20px;
}

.leaderboard-enter-active,
.leaderboard-leave-active {
  transition: opacity 0.5s ease;
}

.leaderboard-enter-from,
.leaderboard-leave-to {
  opacity: 0;
}

.sidebar {
  /* display: flex;
  flex-direction: column; */
}

.code__wrapper {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  font-family: 'Inter-SemiBold', sans-serif;
  font-size: var(--font-tertiary-size);
  color: var(--font-secondary-color);
  gap: 10px;
  line-height: 20px;
  margin-top: 20px;
  margin-bottom: 30px;
}

.code {
  font-size: 20px;
  color: var(--font-primary-color);
  display: flex;
  gap: 5px;
  cursor: pointer;
}

.code img {
  width: 18px;
  height: 18px;
  object-fit: contain;
}

.button__wrapper {
  display: flex;
  flex-direction: column;
  text-align: center;
  gap: 10px;
  margin-bottom: 20px;
  font-size: 12px;
  font-family: 'Inter-Medium', sans-serif;
  color: var(--font-secondary-color);
  line-height: 20px;
}

.separator {
  width: 100%;
  height: 2px;
  background: var(--background-secondary-color);
  border-radius: 50px;
  margin-bottom: 20px;
}

.participants__header {
  display: flex;
  align-items: center;
  gap: 5px;
  font-size: var(--font-quaternary-size);
  font-family: 'Inter-SemiBold', sans-serif;
  color: var(--font-primary-color);
  line-height: 20px;
  margin-bottom: 10px;
}

.participants__header img {
  width: 24px;
  height: 24px;
  object-fit: contain;
}

.participants {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(40px, 40px));
  gap: 5px;
}

.participants img {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  object-fit: contain;
}
</style>
