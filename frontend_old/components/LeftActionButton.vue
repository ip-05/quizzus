<template>
  <button class="button" @click="handleClick">
    <div class="button__img">
      <nuxt-img class="img" :src="imgSrc || img" :alt="alt" />
    </div>
  </button>
</template>

<script setup>
import { computed, defineProps } from 'vue';
import { useGameStore } from '../stores/game';
import { useDynamicIslandStore } from '~/stores/dynamicIsland';

const island = useDynamicIslandStore();
const gameStore = useGameStore();

const props = defineProps({
  mode: {
    type: String,
    default: 'default',
  },
  imgSrc: {
    type: String,
    default: null,
  },
});

const imgs = {
  default: { src: 'svg/icon-default.svg', alt: 'Icon' },
  leave: { src: 'svg/icon-leave.svg', alt: 'Leave' },
  hamburger: { src: 'svg/icon-hamburger-menu.svg', alt: 'Menu' },
  avatar: { src: 'svg/icon-login.svg', alt: 'Login' },
};
const img = computed(() => imgs[props.mode].src);
const alt = computed(() => imgs[props.mode].alt);

const handleClick = () => {
  if (props.mode === 'hamburger') {
    if (gameStore.active) return;
    island.active();
  }
  if (props.mode === 'leave') {
    gameStore.leaveGame();
  }
};
</script>

<style scoped>
.button {
  background: var(--background-main-color);
  border: none;
  outline: solid 3px var(--border-color);
  outline-offset: -3px;
  border-radius: 50%;
  display: flex;
  justify-content: center;
  align-items: center;
  cursor: pointer;
  width: 50px;
  height: 50px;
}

.button__img {
  width: 30px;
  height: 30px;
}

.img {
  width: 100%;
  height: 100%;
  object-fit: contain;
  border-radius: 50%;
}

.popup {
  height: auto !important;
}
</style>
