<script setup lang="ts">
import { computed } from 'vue';
import ButtonSecondary from './buttons/ButtonSecondary.vue';
import { RouterLink, useRoute } from 'vue-router';

const navs = [
  { path: '/discover', name: 'Discover' },
  { path: '/dashboard', name: 'Dashboard' },
  { path: '/library', name: 'Library' },
  { path: '/play', name: 'Join Game' },
  { path: '/create', name: 'Create Game' },
];

const currentPath = computed(() => useRoute().path);
</script>

<template>
  <nav class="nav">
    <RouterLink to="/" class="nav--item nav--logo"><h3>Quizzus</h3></RouterLink>
    <RouterLink v-for="{ path, name } in navs" :key="path" :to="path" class="nav--item nav--link">
      <ButtonSecondary
        :style="currentPath === path ? 'active' : path === '/play' ? 'action' : 'default'"
        >{{ name }}</ButtonSecondary
      >
    </RouterLink>
    <!-- TODO: when authed, show username -->
    <RouterLink to="/login" class="nav--item nav--link"
      ><ButtonSecondary>Login</ButtonSecondary></RouterLink
    >
  </nav>
</template>

<style scoped>
.nav {
  display: flex;
  padding: 20px 30px;
  background: var(--color-background);
}

.nav--logo {
  display: flex;
  align-items: center;
  text-decoration: none;
  color: var(--color-heading);
  margin-right: 20px;
}

.nav--logo * {
  font-size: 18px;
  font-weight: 700;
}

.nav--link {
  text-decoration: none;
}

.nav--item:nth-of-type(5) {
  margin-left: auto;
}
</style>
