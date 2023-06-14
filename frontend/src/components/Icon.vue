<script setup>
import { computed } from 'vue';

const props = defineProps({
  icon: {
    type: String,
    default: 'icon-plus.svg',
    required: true,
  },
});

const svgs = {
  plus: 'icon-plus.svg',
};

const getImageUrl = (name) => new URL(`../assets/svg/${name}`, import.meta.url).href;

const source = computed(() => getImageUrl(svgs[props.icon] ?? props.icon));
</script>

<template>
  <img :src="source" alt="Icon" />
</template>
