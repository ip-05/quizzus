<template>
  <div class="wrapper">
    <div class="content">
      <div class="header">
        <div class="title">Share game</div>
        <div class="description">Generate and share link to this quiz</div>
        <div class="bars">
          <div class="bar bar--active"></div>
          <div class="bar bar--active"></div>
        </div>
      </div>
    </div>
    <div class="main">
      <div class="code__wrapper">
        <div class="code__content">
          <span class="code__title">Game code</span>
          <span class="code">{{ code }}</span>
        </div>
        <regular-button @click="clipboardCode">Copy</regular-button>
      </div>
      <div class="link__wrapper">
        <div class="link__title">Extended link</div>
        <div class="link__form">
          <nuxt-img src="svg/icon-link.svg" alt="Link Icon" />
          <div class="link__text">{{ link }}</div>
          <nuxt-img class="link__clipboard" src="svg/icon-clipboard.svg" alt="Clipboard Icon" @click="clipboardLink" />
        </div>
        <div class="info">
          <nuxt-img src="svg/icon-info.svg" alt="Info Icon" />
          <span>This game will be saved to Workshop, where you can edit or/and start game</span>
        </div>
      </div>
    </div>
    <div class="footer">
      <regular-button @click="$emit('back')">Back</regular-button>
      <regular-button>Save to Workshop and Leave</regular-button>
      <NuxtLink :to="game" class="footer__last">
        <regular-button active>Go to Game</regular-button>
      </NuxtLink>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue';

/**
 * TODO:
 * get code and generate link from CreateGame store
 */

const code = ref('7851-1699-0760');
const link = ref('https://quizzus.fun/join/7851-1699-0760');
const game = ref('/');

const saveToClipboard = (data) => {
  navigator.clipboard.writeText(data);
  alert('Copied: ' + data);
};

const clipboardCode = () => saveToClipboard(code.value);
const clipboardLink = () => saveToClipboard(link.value);
</script>

<style scoped>
.wrapper {
  padding: 30px;
  outline: solid 3px var(--border-color);
  outline-offset: -3px;
  border-radius: 30px;
  background: var(--background-main-color);
}

.title {
  font-size: var(--font-primary-size);
  font-family: 'Inter-SemiBold', sans-serif;
  color: var(--font-primary-color);
  line-height: 30px;
  margin-bottom: 10px;
}

.description {
  font-size: var(--font-tertiary-size);
  font-family: 'Inter-SemiBold', sans-serif;
  color: var(--font-secondary-color);
  line-height: 20px;
  margin-bottom: 30px;
}

.bars {
  display: flex;
  gap: 20px;
}

.bar {
  width: 100%;
  height: 6px;
  background: var(--background-secondary-color);
  border-radius: 6px;
  margin-bottom: 30px;
}

.bar--active {
  background: var(--green-color);
}

.code__wrapper {
  padding: 30px 0px 30px 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  position: relative;
  margin-bottom: 30px;
}

.code__wrapper::after,
.code__wrapper::before {
  content: '';
  position: absolute;
  left: 0;
  right: 0;
  height: 3px;
  border-radius: 3px;
  background: var(--background-secondary-color);
}

.code__wrapper::after {
  top: 0;
}

.code__wrapper::before {
  bottom: 0;
}

.code__content {
  display: flex;
  flex-direction: column;
  gap: 10px;
  color: var(--font-secondary-color);
  font-size: var(--font-tertiary-size);
  font-family: 'Inter-SemiBold', sans-serif;
  line-height: 20px;
}

.code {
  color: var(--font-primary-color);
  font-size: 20px;
  font-family: 'Inter-SemiBold', sans-serif;
  line-height: 30px;
}

.link__wrapper {
  color: var(--font-primary-color);
  font-size: var(--font-tertiary-size);
  font-family: 'Inter-SemiBold', sans-serif;
  line-height: 20px;
}

.link__title {
  margin-bottom: 10px;
}

.link__form {
  display: flex;
  align-items: center;
  gap: 10px;
  background: var(--background-secondary-color);
  border-radius: 10px;
  padding: 10px 15px;
  margin-bottom: 15px;
}

.link__clipboard {
  margin-left: auto;
  cursor: pointer;
}

.info {
  display: flex;
  align-items: center;
  gap: 5px;
  font-size: var(--font-quaternary-size);
  color: var(--font-secondary-color);
  margin-bottom: 60px;
}

.footer {
  display: flex;
  gap: 10px;
}

.footer__last {
  margin-left: auto;
}
</style>
