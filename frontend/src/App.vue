<script setup>
import { RouterView } from 'vue-router';
import HeaderNav from './components/HeaderNav.vue';
</script>

<template>
  <header>
    <div class="wrapper">
      <HeaderNav />
    </div>
  </header>

  <RouterView />
</template>

<style scoped></style>
