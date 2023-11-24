<script setup lang="ts">
import { MessageType } from '@/types';
import { Icons } from '@/types';
import Icon from './icons/Icon.vue';
import { computed } from 'vue';

interface Props {
  type: MessageType | undefined;
  label: string | undefined;
  message: string | undefined;
}

const props = withDefaults(defineProps<Props>(), {
  label: 'Label',
  message: 'Message',
});

const currentStyle = computed(() => `toast--${props.type}`);
const currentIcon = computed<Icons>(() => {
  if (props.type === 'error') return 'error-red';
  if (props.type === 'warning') return 'warning-yellow';
  if (props.type === 'info') return 'info-blue';
  if (props.type === 'success') return 'success-green';
  return 'info';
});
</script>

<template>
  <div class="toast" :class="currentStyle">
    <Icon :icon="currentIcon" class="icon"></Icon>
    <div class="text">
      <h3 class="flask__label">{{ label }}</h3>
      <p class="flask__message">{{ message }}</p>
    </div>
  </div>
</template>

<style scoped>
.toast {
  padding: 15px 20px;
  background: var(--color-background);
  border-radius: 20px;
  position: fixed;
  bottom: 50px;
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  gap: 15px;
  align-items: center;
  line-height: 1;
  outline: 2px solid var(--c-white);
  outline-offset: -2px;
  box-shadow:
    0px 4px 9px 0px rgba(0, 0, 0, 0.02),
    0px 17px 17px 0px rgba(0, 0, 0, 0.02),
    0px 39px 23px 0px rgba(0, 0, 0, 0.01),
    0px 69px 28px 0px rgba(0, 0, 0, 0),
    0px 108px 30px 0px rgba(0, 0, 0, 0);
}

.icon {
  width: 24px;
  height: 24px;
  object-fit: contain;
}

.flask__label {
  margin-bottom: 5px;
}

.toast--error {
  outline-color: var(--c-red);
}

.toast--warning {
  outline-color: var(--c-yellow);
}

.toast--info {
  outline-color: var(--c-blue);
}

.toast--success {
  outline-color: var(--c-green);
}
</style>
