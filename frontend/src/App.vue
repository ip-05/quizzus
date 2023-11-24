<script setup lang="ts">
import { RouterView } from 'vue-router';

import HeaderNav from './components/HeaderNav.vue';
import Footer from './components/Footer.vue';
import ToastMessage from './components/ToastMessage.vue';
import { useToastMessageStore } from './stores/toast';

const toast = useToastMessageStore();
</script>

<template>
  <div class="wrapper">
    <header class="header">
      <HeaderNav />
    </header>
    <main class="main">
      <div class="main__content">
        <transition name="page" mode="out-in">
          <router-view></router-view>
        </transition>
      </div>
    </main>
    <Footer />
  </div>
  <Transition name="toast">
    <ToastMessage
      v-if="toast.hasMessage"
      :type="toast.type"
      :label="toast.label"
      :message="toast.message"
    />
  </Transition>
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

.toast-enter-active,
.toast-leave-active {
  transition: opacity 0.15s ease;
}

.toast-enter-from,
.toast-leave-to {
  opacity: 0;
}

.page-enter-active,
.page-leave-active {
  transition: opacity 0.15s ease;
}

.page-enter-from,
.page-leave-to {
  opacity: 0;
}
</style>
