<script setup lang="ts">
import { computed } from 'vue';
import ButtonSecondary from './buttons/ButtonSecondary.vue';
import { RouterLink, useRoute, useRouter } from 'vue-router';
import { useAuthStore } from '@/stores/auth';
import { useToastMessageStore } from '@/stores/toast';

const authStore = useAuthStore();
const router = useRouter();
const currentPath = computed(() => useRoute().path);

const navs = [
  { path: '/discover', name: 'Discover' },
  { path: '/dashboard', name: 'Dashboard' },
  { path: '/library', name: 'Library' },
  { path: '/join', name: 'Join Game' },
  { path: '/create', name: 'Create Game' },
];

const unauthorizedPopup = (path: string) => {
  const toast = useToastMessageStore();
  const route = router.resolve(path);
  if (route.meta.requiresAuth && !authStore.isAuthenticated) {
    toast.register('error', 'Error', 'You are not logged in.');
  }
};
</script>

<template>
  <nav class="nav">
    <RouterLink to="/" class="nav--item nav--logo"><h3>Quizzus</h3></RouterLink>
    <RouterLink v-for="{ path, name } in navs" :key="path" :to="path" class="nav--item nav--link">
      <ButtonSecondary
        :style="currentPath === path ? 'active' : path === '/join' ? 'action' : 'default'"
        @click="unauthorizedPopup(path)"
        >{{ name }}</ButtonSecondary
      >
    </RouterLink>
    <!-- TODO: when authed, show username -->
    <RouterLink to="/login" class="nav--item nav--link"
      ><ButtonSecondary @click="unauthorizedPopup('/login')">Login</ButtonSecondary></RouterLink
    >
  </nav>
</template>

<style scoped>
.nav {
  display: flex;
  padding: 20px 30px;
}

.nav--logo {
  display: flex;
  align-items: center;
  text-decoration: none;
  color: var(--color-heading);
  margin-right: 20px;
}

.nav--logo h3 {
  font-family: 'Daydream', serif;
  font-size: 14px;
  letter-spacing: 1px;
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
@/stores/toast
