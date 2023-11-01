<script setup lang="ts">
import { computed } from 'vue';

import { Icons } from '@/types';

interface Props {
  icon: Icons;
}

const props = defineProps<Props>();

const getImageUrl = (name: string) =>
  new URL(`../../assets/svg/icon-${name}.svg`, import.meta.url).href;

const source = computed(() => getImageUrl(props.icon));
console.log(props.icon);
</script>

<template>
  <img :src="source" alt="Icon" />
</template>
