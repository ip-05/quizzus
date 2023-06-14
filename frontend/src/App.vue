<script setup>
import { RouterLink, RouterView } from 'vue-router';
</script>

<template>
  <header>
    <div class="wrapper">
      <nav>
        <RouterLink to="/discover">Discover</RouterLink>
        <RouterLink to="/dashboard">Dashboard</RouterLink>
        <RouterLink to="/library">Library</RouterLink>
        <RouterLink to="/starred">starred</RouterLink>
        <RouterLink to="/recent">Recently Played</RouterLink>
        <RouterLink to="/create">Create Game</RouterLink>
        <RouterLink to="/edit">Edit Game</RouterLink>
        <RouterLink to="/console">Console</RouterLink>
        <RouterLink to="/game">Game Page</RouterLink>
        <RouterLink to="/play">Play Game</RouterLink>
      </nav>
    </div>
  </header>

  <RouterView />
</template>

<style scoped></style>
