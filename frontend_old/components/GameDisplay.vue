<template>
  <div class="wrapper" :style="{ 'outline-color': state === 'game-correct' ? correctColor : null }">
    <div class="content">
      <div class="header">
        <div class="participants">
          <div class="indicator">
            <nuxt-img class="indicator__img" src="svg/icon-participants.svg" />
            {{ gameStore.participantsNumber }}
          </div>
        </div>
        <div class="info info__top">
          <span v-if="state === 'game-over-admin'">Game Over</span>
          <span
            v-if="state === 'game-graph-admin'"
            style="color: var(--font-primary-color); font-family: 'Inter-SemiBold', sans-serif"
            >{{ gameStore.currentQuestion.name }}</span
          >
          <span v-else>{{ gameStore.topic }}</span>
        </div>
        <div class="points">
          <div class="indicator">
            <nuxt-img class="indicator__img" src="svg/icon-points.svg" />
            {{ gameStore.points }}
          </div>
        </div>
      </div>
      <div class="main">
        <!-- GAME WAITING -->
        <div v-if="state === 'game-wait'" class="waiting__screen">
          <span>{{ gameStore.countdown ? gameStore.countdown : gameStore.topic }}</span>
        </div>

        <!-- GAME WAITING FOR ADMIN -->
        <div v-if="state === 'game-wait-admin'" class="waiting__screen waiting__screen-admin">
          <span>{{ gameStore.countdown ? gameStore.countdown : gameStore.topic }}</span>
        </div>

        <!-- GAME IN -->
        <div v-if="state === 'game-in'" class="in__screen">
          <span>{{ gameStore.currentQuestion.name }}</span>
        </div>

        <!-- GAME CORRECT -->
        <div v-if="state === 'game-correct'" class="correct__screen">
          <ClientOnly>
            <ConfettiExplosion
              class="confetti"
              :particle-count="100"
              :force="1"
              :duration="10000"
              :colors="[correctColor]"
            />
          </ClientOnly>
          <div class="totalpoints">{{ gameStore.totalPoints }}</div>
          <div class="pluspoints pluspoints--plus">+ {{ gameStore.points }} points</div>
          <div class="question">{{ gameStore.currentQuestion.name }}</div>
          <div class="answer">{{ gameStore.currentQuestion.answer }}</div>
        </div>

        <!-- GAME WRONG -->
        <div v-if="state === 'game-wrong'" class="correct__screen">
          <div class="totalpoints">{{ gameStore.totalPoints }}</div>
          <div class="pluspoints pluspoints--zero">+ 0 points</div>
          <div class="question">{{ gameStore.currentQuestion.name }}</div>
          <div class="answer">{{ gameStore.currentQuestion.answer }}</div>
        </div>

        <!-- GAME OVER -->
        <div v-if="state === 'game-over'" class="over__screen">
          <div class="totalpoints">{{ gameStore.totalPoints }}</div>
          <div class="pluspoints">your points</div>
          <div v-if="gameStore.top3" class="top3">
            <div v-for="(place, id) in gameStore.top3" :key="id" class="top3__item" :class="places[id].class">
              <nuxt-img :src="places[id].src" />
              <span>{{ place.participant }}</span>
            </div>
          </div>
        </div>

        <!-- GAME OVER FOR ADMIN -->
        <div v-if="state === 'game-over-admin'" class="over__screen over__screen-admin">
          <div class="topic">{{ gameStore.topic }}</div>
          <div v-if="gameStore.top3" class="top3">
            <div v-for="(place, id) in gameStore.top3" :key="id" class="top3__item" :class="places[id].class">
              <nuxt-img :src="places[id].src" />
              <span>{{ place.participant }}</span>
            </div>
          </div>
        </div>

        <!-- GAME GRAPH FOR ADMIN -->
        <div v-if="state === 'game-graph-admin'" class="graph__screen">
          <div class="graph">
            <div class="graph__col">
              <div class="graph__data">
                <div class="graph__block graph__red" :style="{ height: `${gameStore.graphAnswers[0]}%` }"></div>
              </div>
              <div class="graph__number">{{ gameStore.submittedAnswers[0].answers }}</div>
            </div>
            <div class="graph__col">
              <div class="graph__data">
                <div class="graph__block graph__blue" :style="{ height: `${gameStore.graphAnswers[1]}%` }"></div>
              </div>
              <div class="graph__number">{{ gameStore.submittedAnswers[1].answers }}</div>
            </div>
            <div class="graph__col">
              <div class="graph__data">
                <div class="graph__block graph__green" :style="{ height: `${gameStore.graphAnswers[2]}%` }"></div>
              </div>
              <div class="graph__number">{{ gameStore.submittedAnswers[2].answers }}</div>
            </div>
            <div class="graph__col">
              <div class="graph__data">
                <div class="graph__block graph__yellow" :style="{ height: `${gameStore.graphAnswers[3]}%` }"></div>
              </div>
              <div class="graph__number">{{ gameStore.submittedAnswers[3].answers }}</div>
            </div>
          </div>
        </div>
      </div>
      <div class="footer">
        <div class="time">
          <div class="indicator">
            <nuxt-img class="indicator__img" src="svg/icon-time.svg" />
            {{ gameStore.roundTime }}
          </div>
        </div>
        <div class="info info__bottom">
          <span v-if="state !== 'game-over' && state !== 'game-over-admin'">{{ gameStore.hint }}</span>
          <div v-else class="indicator" style="cursor: pointer" @click="openLeaderboard">
            <nuxt-img class="indicator__img" src="svg/icon-leaderboard.svg" />
            Leaderboard
            <nuxt-img class="indicator__img" src="svg/icon-leaderboard-arrow.svg" />
          </div>
        </div>
        <div class="questions">
          <div class="indicator">
            <nuxt-img class="indicator__img" src="svg/icon-questions.svg" />
            {{ gameStore.questionsNumber }}
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue';
import { useGameStore } from '../stores/game';

const gameStore = useGameStore();
const state = computed(() => gameStore.state);

const correct = computed(() => gameStore.correctAnswer);

const correctColor = computed(() => {
  const colors = ['var(--red-color)', 'var(--blue-color)', 'var(--green-color)', 'var(--yellow-color)'];
  for (const [idx, option] of gameStore.currentQuestion.options.entries()) {
    if (option.id === correct.value.id) {
      return colors[idx];
    }
  }
  return colors[0];
});

// const gold = computed(() => gameStore.top3[0].participant);
// const silver = computed(() => gameStore.top3[1].participant);
// const bronze = computed(() => gameStore.top3[2].participant);

const places = [
  { src: 'svg/icon-gold.svg', class: 'gold' },
  { src: 'svg/icon-silver.svg', class: 'silver' },
  { src: 'svg/icon-bronze.svg', class: 'bronze' },
];

const openLeaderboard = () => {
  gameStore.leaderboardShown = !gameStore.leaderboardShown;
};
</script>

<style scoped>
.wrapper {
  padding: 15px;
  outline: solid 3px var(--border-color);
  outline-offset: -3px;
  border-radius: 30px;
  background: var(--background-main-color);
  transition: all 250ms ease;
  overflow: hidden;
  position: relative;
  transition: opacity 250ms 250ms ease;
}

.confetti {
  position: absolute;
  /* opacity: 0.5; */
}

.wrapper--correct {
  outline-color: var(--yellow-color);
  /* animation: color 5s infinite linear; */
}

.header,
.footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.indicator {
  display: flex;
  gap: 5px;
  justify-content: center;
  align-items: center;
  padding: 5px 10px;
  background: var(--background-secondary-color);
  border-radius: 30px;
  font-family: 'Inter-SemiBold', sans-serif;
  font-size: var(--font-quaternary-size);
  color: var(--font-primary-color);
  line-height: 20px;
  user-select: none;
}

.indicator__img {
  width: 14px;
  height: 14px;
  object-fit: contain;
}

.info {
  font-family: 'Inter-Medium', sans-serif;
  font-size: var(--font-quaternary-size);
  color: var(--font-secondary-color);
  line-height: 20px;
}

/* .screen-enter-active,
.screen-leave-active {
  transition: opacity 0.5s ease;
}

.screen-enter-from,
.screen-leave-to {
  opacity: 0;
} */

.main {
}

.waiting__screen,
.in__screen,
.correct__screen,
.over__screen,
.graph__screen {
  font-size: var(--font-primary-size);
  color: var(--font-primary-color);
  font-family: 'Inter-Bold', sans-serif;
  line-height: 30px;
  padding: 202px 0;
  display: flex;
  justify-content: center;
  align-items: center;
  flex-direction: column;
}

.waiting__screen-admin {
  padding: 140px 0;
}

.in__screen {
  padding: 90px 0;
}

.correct__screen {
  padding-top: 94px;
  padding-bottom: 87px;
}

.totalpoints {
  font-size: 48px;
  margin-bottom: 5px;
  line-height: 58px;
}

.pluspoints {
  font-family: 'Inter-Medium', sans-serif;
  font-size: var(--font-secondary-size);
  color: var(--font-primary-color);
  line-height: 30px;
  margin-bottom: 40px;
}

.pluspoints--plus {
  color: var(--green-color);
}

.pluspoints--zero {
  color: var(--red-color);
}

.question {
  color: var(--font-secondary-color);
  font-size: var(--font-secondary-size);
  line-height: 30px;
  margin-bottom: 60px;
}

.answer {
  font-size: 22px;
  line-height: 30px;
  font-family: 'Inter-Medium', sans-serif;
}

.over__screen {
  padding: 94px 0;
}

.over__screen .pluspoints {
  margin-bottom: 40px;
}

.top3 {
  display: flex;
  align-items: center;
  gap: 100px;
}

.top3__item {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  gap: 15px;
  font-size: var(--font-tertiary-size);
  font-family: 'Inter-Medium', sans-serif;
  color: var(--font-secondary-color);
  line-height: 30px;
}

.top3__item img {
  width: 48px;
  height: 48px;
  object-fit: contain;
}

.gold {
  justify-self: center;
  align-self: center;
  padding-bottom: 20px;
}

.over__screen-admin {
  padding: 54px 0;
}

.topic {
  margin-bottom: 59px;
}

.graph__screen {
  padding: 63px 0;
}

.graph {
  display: flex;
  gap: 20px;
}

.graph__col {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 10px;
  flex-direction: column;
}

.graph__data {
  width: 100px;
  min-height: 144px;
  position: relative;
}

.graph__number {
  font-size: var(--font-quaternary-size);
  line-height: 20px;
  font-family: 'Inter-SemiBold', sans-serif;
  padding: 5px 10px;
  background: var(--background-secondary-color);
  border-radius: 30px;
}

.graph__block {
  border-radius: 15px;
  position: absolute;
  left: 0;
  right: 0;
  bottom: 0;
  transition: all 250ms ease;
}

.graph__red {
  background: var(--red-color);
}
.graph__blue {
  background: var(--blue-color);
}
.graph__green {
  background: var(--green-color);
}
.graph__yellow {
  background: var(--yellow-color);
}
</style>
