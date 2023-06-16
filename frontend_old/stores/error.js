import { defineStore } from 'pinia';
import { ref } from 'vue';

export const useErrorStore = defineStore('ErrorStore', () => {
  const message = ref(null);

  function reset() {
    message.value = null;
  }

  return { message, reset };
});
