<template>
  <div class="nav__wrapper" :class="stateStyle">
    <Transition name="backdrop">
      <div v-if="isOpen" class="backdrop" @click="closeIsland"></div>
    </Transition>
    <nav class="nav">
      <div class="header" @click.capture="openIsland">
        <div class="header__status">Welcome to Quizzus</div>
        <div v-if="isOpen" class="header__icon">
          <nuxt-img style="cursor: pointer" src="svg/icon-close.svg" @click="closeIsland" />
        </div>
      </div>
      <Transition name="body">
        <div v-if="isOpen" class="body">
          <div class="body__navigation">
            <NuxtLink class="body__navlink" to="/new">
              <island-button type="createroom" />
            </NuxtLink>
            <NuxtLink class="body__navlink" to="/">
              <island-button type="workshop" />
            </NuxtLink>
            <island-button type="join">
              <enter-code-input />
            </island-button>
            <NuxtLink class="body__navlink" to="/">
              <island-button type="settings" />
            </NuxtLink>
          </div>
          <div class="body__footer">
            <span class="body__copyright">© 2022 Quizzus All rights reserved</span>
            <span class="body__users"> <span></span> 20 active users</span>
            <span class="body__helplinks">
              <NuxtLink to="/" class="body__helplink">Help</NuxtLink>
              <NuxtLink to="/" class="body__helplink">Privacy policy</NuxtLink>
              <NuxtLink to="/" class="body__helplink">Terms and conditions</NuxtLink>
            </span>
          </div>
        </div>
      </Transition>
    </nav>
  </div>
</template>

<script setup>
import { computed } from 'vue';
import { useDynamicIslandStore } from '~/stores/dynamicIsland';

const island = useDynamicIslandStore();

// defining specific styles depended on state
const states = {
  default: 'nav--default',
  active: 'nav--active',
  game_on: 'nav--game_on',
  game_waiting: 'nav--game_waiting',
};
const stateStyle = computed(() => states[island.state]);

const isOpen = computed(() => island.state === 'active');

const openIsland = () => {
  if (island.state === 'default') {
    island.active();
  }
};

const closeIsland = () => {
  island.default();
};
</script>

<style scoped>
.backdrop {
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  position: fixed;
  width: 100vw;
  height: 100vh;
  background: rgba(0, 0, 0, 0.2);
  transition: all 250ms ease;
}

.backdrop-enter-active,
.backdrop-leave-active {
  transition: opacity 0.5s ease;
}

.backdrop-enter-from,
.backdrop-leave-to {
  opacity: 0;
}

.nav__wrapper {
  text-align: center;
  transition: all 250ms ease;
}

.nav--active {
  left: 0px !important;
  right: 0px !important;
  z-index: 10;
}

.nav--active > .nav {
  border-radius: 20px;
}

.nav--default {
  cursor: pointer;
}

.nav {
  position: relative;
  padding: 15px 10px;
  outline: solid 3px var(--border-color);
  outline-offset: -3px;
  border-radius: 40px;
  background: var(--background-main-color);
  transition: all 250ms ease;
}

.header {
  width: 100%;
  height: auto;
  display: flex;
  justify-content: center;
  align-items: center;
  position: relative;
}

.header__icon {
  position: absolute;
  right: 0px;
  width: 30px;
  height: 30px;
}

.header__icon img {
  width: 100%;
  height: 100%;
  object-fit: contain;
}

.header__status {
  font-family: 'Inter-Medium', sans-serif;
  font-size: var(--font-secondary-size);
}

.body {
  padding: 20px;
}

.body-enter-active,
.body-leave-active {
  transition: opacity 125ms ease;
}

.body-enter-from,
.body-leave-to {
  opacity: 0;
}

.body__navigation {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 30px;
  margin-bottom: 30px;
}

.body__footer {
  color: var(--font-secondary-color);
  font-size: var(--font-quaternary-size);
  font-family: 'Inter-Medium', sans-serif;
  display: flex;
  justify-content: space-between;
}

.body__users {
  display: flex;
  align-items: center;
  gap: 5px;
}

.body__users span {
  display: inline-block;
  width: 5px;
  height: 5px;
  background: var(--green-color);
  border-radius: 50%;
}

.body__navlink {
  text-decoration: none;
}

.body__helplinks {
  display: flex;
  gap: 10px;
}

.body__helplink {
  color: var(--font-secondary-color);
}
</style>
