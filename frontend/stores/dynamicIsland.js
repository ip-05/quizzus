import { defineStore } from 'pinia';

export const useDynamicIslandStore = defineStore('DynamicIslandStore', {
  state: () => ({ state: 'default' }),
  getters: {
    getState: (state) => state.state,
  },
  actions: {
    default() {
      this.state = 'default';
    },
    active() {
      this.state = 'active';
    },
    gameOn() {
      this.state = 'game_on';
    },
    gameWaiting() {
      this.state = 'game_waiting';
    },
  },
});
