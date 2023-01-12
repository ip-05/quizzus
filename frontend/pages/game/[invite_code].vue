<template>
  <div class="main">
    <div class="content">
      <game-display style="margin-bottom: 20px" />
      <game-table v-if="gameStore.leaderboardShown" mode="leaderboard">Participants</game-table>
      <div v-if="gameStore.state === 'game-in'" class="options">
        <form class="form" action="">
          <!-- TODO: refactor for using v-for -->
          <p class="option">
            <input v-model="gameStore.selectedOption" id="a" type="radio" name="option" :value="0" />
            <label class="label" for="a">{{ options.a }}</label>
          </p>
          <p class="option">
            <input v-model="gameStore.selectedOption" id="b" type="radio" name="option" :value="1" />
            <label class="label" for="b"> {{ options.b }}</label>
          </p>
          <p class="option">
            <input v-model="gameStore.selectedOption" id="c" type="radio" name="option" :value="2" />
            <label class="label" for="c"> {{ options.c }}</label>
          </p>
          <p class="option">
            <input v-model="gameStore.selectedOption" id="d" type="radio" name="option" :value="3" />
            <label class="label" for="d"> {{ options.d }}</label>
          </p>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted } from 'vue';
import { useGameStore } from '../../stores/game';
import { useSocketStore } from '../../stores/socket';
import { useRoute } from '#imports';

const gameStore = useGameStore();
const { params } = useRoute();

onMounted(() => {
  console.log('mounted');
  gameStore.inviteCode = params.invite_code;
  const socketStore = useSocketStore();
  socketStore.joinGame(params.invite_code);
});

const options = computed(() => ({
  a: gameStore.currentQuestion.options[0].name,
  b: gameStore.currentQuestion.options[1].name,
  c: gameStore.currentQuestion.options[2].name,
  d: gameStore.currentQuestion.options[3].name,
}));
</script>

<style scoped>
.main {
  margin-top: 100px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.content {
  max-width: 1080px;
  width: calc(100% - 20px);
  height: auto;
}

.options-enter-active,
.options-leave-active {
  transition: opacity 0.5s ease;
}

.options-enter-from,
.options-leave-to {
  opacity: 0;
}

.options {
  margin-top: 20px;
}

.form {
  display: grid;
  grid-template-columns: 1fr 1fr;
  grid-template-rows: auto auto;
  gap: 20px;
}

.option {
  background: var(--background-secondary-color);
  border-radius: 30px;
  position: relative;
}

.label {
  display: flex;
  align-items: center;
  padding: 30px;
  padding-left: 82px;
  font-size: var(--font-primary-size);
  color: var(--font-primary-color);
  font-family: 'Inter-SemiBold', sans-serif;
  outline: solid 3px var(--border-color);
  outline-offset: -3px;
  border-radius: 30px;
  transition: all 250ms ease;
  user-select: none;
  cursor: pointer;
}

.label[for='a']:hover {
  outline-color: var(--red-color);
}

.label[for='b']:hover {
  outline-color: var(--blue-color);
}

.label[for='c']:hover {
  outline-color: var(--green-color);
}

.label[for='d']:hover {
  outline-color: var(--yellow-color);
}

.icon {
  width: 32px;
  height: 32px;
  object-fit: contain;
}

input[type='radio'] {
  -webkit-appearance: none;
  appearance: none;
  background-color: transparent;
  margin: 0;
  position: absolute;
  top: 50%;
  transform: translateY(-50%);
  left: 30px;
  width: 32px;
  height: 32px;
  background-position: center;
  background-size: cover;
  cursor: pointer;
}

input[type='radio'][id='a'] {
  background-image: url('/svg/icon-triangle.svg');
}

input[type='radio'][id='a']:checked {
  background-image: url('/svg/icon-triangle-active.svg');
}

input[type='radio'][id='a']:checked + .label {
  outline-color: var(--red-color);
  background: var(--red-color);
  color: white;
}

input[type='radio'][id='b'] {
  background-image: url('/svg/icon-circle.svg');
}

input[type='radio'][id='b']:checked {
  background-image: url('/svg/icon-circle-active.svg');
}

input[type='radio'][id='b']:checked + .label {
  outline-color: var(--blue-color);
  background: var(--blue-color);
  color: white;
}

input[type='radio'][id='c'] {
  background-image: url('/svg/icon-square.svg');
}

input[type='radio'][id='c']:checked {
  background-image: url('/svg/icon-square-active.svg');
}

input[type='radio'][id='c']:checked + .label {
  outline-color: var(--green-color);
  background: var(--green-color);
  color: white;
}

input[type='radio'][id='d'] {
  background-image: url('/svg/icon-diamond.svg');
}

input[type='radio'][id='d']:checked {
  background-image: url('/svg/icon-diamond-active.svg');
}

input[type='radio'][id='d']:checked + .label {
  outline-color: var(--yellow-color);
  background: var(--yellow-color);
  color: white;
}
</style>
