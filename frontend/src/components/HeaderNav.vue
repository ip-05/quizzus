<script setup>
import ButtonMedium from './buttons/ButtonMedium.vue';
import { RouterLink } from 'vue-router';
</script>

<template>
  <nav class="nav">
    <div class="left">
      <RouterLink to="/discover" class="logo"><h3>Quizzus</h3></RouterLink>
      <RouterLink to="/discover" class="link"><ButtonMedium>Discover</ButtonMedium></RouterLink>
      <RouterLink to="/dashboard" class="link"><ButtonMedium>Dashboard</ButtonMedium></RouterLink>
      <RouterLink to="/library" class="link"><ButtonMedium>Library</ButtonMedium></RouterLink>
    </div>
    <div class="right">
      <RouterLink to="/play" class="link"
        ><ButtonMedium state="rainbow">Join Game</ButtonMedium></RouterLink
      >
      <RouterLink to="/create" class="link"><ButtonMedium>Create Game</ButtonMedium></RouterLink>
      <!-- TODO: when authed, show username -->
      <RouterLink to="/login" class="link"><ButtonMedium>Login</ButtonMedium></RouterLink>
    </div>
  </nav>
</template>

<style scoped>
.nav {
  display: flex;
  justify-content: space-between;
  padding: 20px 30px;
  background: var(--color-background);
}

.logo {
  display: flex;
  align-items: center;
  text-decoration: none;
  color: var(--color-heading);
  margin-right: 20px;
}

.logo * {
  font-size: 18px;
  font-weight: 700;
}

.link {
  text-decoration: none;
}

.left,
.right {
  display: flex;
}
</style>
