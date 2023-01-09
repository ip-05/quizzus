<template>
  <div class="wrapper">
    <header class="header__wrapper">
      <div class="header">
        <left-action-button :mode="gameStore.active ? 'leave' : 'hamburger'" class="header__first" />
        <dynamic-island class="header__center" />
        <profile-popup class="header__last" />
      </div>
    </header>
    <main class="main">
      <NuxtPage />
    </main>
    <footer class="footer">
      <div>
        <a class="footer__link" href="https://github.com/ip-05/quizzus">
          <nuxt-img class="footer__icon" src="svg/icon-github.svg" alt="GitHub Icon" />
        </a>
      </div>
      <span class="footer__copyright">© 2022 Quizzus. All rights reserved</span>
    </footer>
  </div>
</template>

<script setup>
import { useGameStore } from './stores/game';

const gameStore = useGameStore();
</script>

<style scoped>
.wrapper {
  width: 100%;
  min-height: 100vh;
  display: grid;
  grid-template-rows: 70px 1fr auto;
}

.header {
  max-width: 1280px;
  width: calc(100% - 20px);
  height: auto;
  display: grid;
  grid-template-columns: 50px 1fr 50px;
  gap: 10px;
  margin: 0 auto;
  padding-top: 10px;
  position: relative;
}

.header__first {
  position: absolute;
  top: 10px;
  left: 0;
}

.header__center {
  position: absolute;
  top: 10px;
  left: 60px;
  right: 60px;
}

.header__last {
  position: absolute;
  top: 10px;
  right: 0;
  z-index: 5;
  display: flex;
  flex-direction: column;
  align-items: flex-end;
}

.main {
  /* margin-top: 70px; */
  width: 100%;
  height: 100%;
}

.main__content {
  width: 100%;
  height: 100%;
}

.footer {
  width: 100%;
  height: auto;
  display: flex;
  flex-direction: column;
  justify-content: center;
  gap: 10px;
  align-items: center;
  font-family: 'Inter-Medium', sans-serif;
  font-size: var(--font-quaternary-size);
  color: var(--font-secondary-color);
  padding-bottom: 30px;
}

.footer__link {
  display: block;
  width: 24px;
  height: 24px;
}

.footer__icon {
  width: 100%;
  height: 100%;
  object-fit: contain;
}
</style>
