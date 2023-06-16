<script setup>
import { RouterView } from 'vue-router';

import HeaderNav from './components/HeaderNav.vue';
import Footer from './components/Footer.vue';
</script>

<template>
  <div class="wrapper">
    <header class="header">
      <HeaderNav />
    </header>
    <main class="main">
      <div class="main__content">
        <RouterView />
      </div>
    </main>
    <Footer />
  </div>
</template>

<style scoped>
.wrapper {
  min-height: 100vh;
  display: grid;
  grid-template-rows: auto 1fr auto;
}

.header {
  position: sticky;
  top: 0;
}

.main {
  width: 100%;
  height: auto;
  display: flex;
  justify-content: center;
}

.main__content {
  max-width: 1024px;
  width: calc(100% - 30px * 2);
}
</style>
