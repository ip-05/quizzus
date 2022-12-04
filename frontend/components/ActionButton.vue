<template>
  <button class="button" @click="handleClick">
    <div class="button__img">
      <img v-if="type === 'hamburger'" src="@/assets/svg/icon-hamburger-menu.svg" alt="Menu" />
      <img v-if="type === 'avatar'" src="@/assets/svg/icon-login.svg" alt="Avatar" />
    </div>
  </button>
</template>

<script setup>
import { defineProps } from 'vue';
import { useDynamicIslandStore } from '~/stores/dynamicIsland';

const island = useDynamicIslandStore();

const props = defineProps({
  type: {
    type: String,
    default: 'hamburger',
  },
});

const handleClick = () => {
  if (props.type === 'hamburger') {
    island.active();
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

.button_img img {
  width: 100%;
  height: 100%;
  object-fit: contain;
}
</style>
