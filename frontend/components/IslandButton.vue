<template>
  <form action="" class="form" :class="{ 'form--join': type === 'join' }">
    <div class="img">
      <nuxt-img :src="img" :alt="alt" />
    </div>
    <slot>
      <button class="button" type="submit">{{ alt }}</button>
    </slot>
  </form>
</template>

<script setup>
import { computed, defineProps } from 'vue';

const props = defineProps({
  type: {
    type: String,
    default: 'default',
  },
});

const imgs = {
  createroom: { src: 'svg/icon-createroom-giant.svg', alt: 'Create Room' },
  workshop: { src: 'svg/icon-workshop-giant.svg', alt: 'Workshop' },
  join: { src: 'svg/icon-entercode-giant.svg', alt: 'Enter Code to Join' },
  settings: { src: 'svg/icon-settings-giant.svg', alt: 'Settings' },
};
const img = computed(() => imgs[props.type].src);
const alt = computed(() => imgs[props.type].alt);
</script>

<style scoped>
.form {
  width: 100%;
  height: 100%;
  background: var(--background-secondary-color);
  border-radius: 30px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
}

:slotted(form) {
  margin-bottom: -15px;
}

.form--join > .img {
  transform: translateY(-32px);
}

.img {
  width: 128px;
  height: 128px;
  margin-top: 48px;
}

.img img {
  width: 100%;
  height: 100%;
  object-fit: contain;
}

.button {
  border: none;
  color: var(--font-primary-color);
  font-size: 22px;
  font-family: 'Inter-Medium', sans-serif;
  margin-top: 47px;
  padding-bottom: 30px;
  line-height: 30px;
}
</style>
