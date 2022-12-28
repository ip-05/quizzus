import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import { useCookie } from '#imports';

export const useNewGameStore = defineStore('NewGameStore', () => {
  const tokenCookie = useCookie('token');

  const code = ref(null);
  const link = computed(() => `https://quizzus.fun/game/${code.value}`);
  const game = ref(null);
  const isEmptyGame = computed(() => game.value === null);

  function createGame(data) {
    game.value = data;
  }

  function resetGame() {
    game.value = null;
  }

  // Upload game to db
  async function postNewGame() {}

  // Get game code by id
  async function getGameCode() {}

  return { code, link, game, isEmptyGame, postNewGame, getGameCode, createGame, resetGame };
});
