import { ref } from 'vue';
import { defineStore } from 'pinia';
import { MessageType } from '@/types';

export const useToastMessageStore = defineStore('toastMessageStore', () => {
  const hasMessage = ref(false);
  const timeout = ref(5000);

  const type = ref<MessageType>();
  const label = ref<string>();
  const message = ref<string>();

  function register(msgType: MessageType, msgLabel: string, msg: string) {
    hasMessage.value = true;
    type.value = msgType;
    label.value = msgLabel;
    message.value = msg;
    reset();
  }

  function reset() {
    setTimeout(forceReset, timeout.value);
  }

  function forceReset() {
    hasMessage.value = false;
    type.value = 'unset';
    label.value = '';
    message.value = '';
  }

  return { type, label, message, hasMessage, register, forceReset };
});
