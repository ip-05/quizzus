<template>
  <div class="wrapper">
    <div class="content">
      <div class="header">
        <div class="title">Create new game</div>
        <div class="description">
          Prepare questions, share link with your friends and start the game. Fill all fields to continue
        </div>
        <div class="bars">
          <div class="bar bar--active"></div>
          <div class="bar"></div>
        </div>
      </div>
      <div class="main">
        <div class="about">
          <div class="about__title">Topic of the Game</div>
          <div class="about__title">Points per question</div>
          <div class="about__title">Round time</div>
          <regular-input v-model="game.topic" placeholder="Enter the topic" img="svg/icon-pen-darker.svg" />
          <regular-input v-model="game.points" placeholder="From 0.1 to 100" img="svg/icon-star.svg" />
          <regular-input v-model="game.time" placeholder="From 10 to 60 seconds" img="svg/icon-clock-darker.svg" />
        </div>
        <div class="questions">
          <div class="questions__title">List of questions</div>
          <TransitionGroup name="questions-list" tag="div">
            <question-block
              v-for="(question, idx) in game.questions"
              :key="question.id"
              :generated-id="idx"
              :question="question"
              :removable="removable"
              class="question"
              @update-questions="handleUpdateQuestions"
              @remove="removeQuestionBlock"
            />
          </TransitionGroup>
          <div class="questions__add">
            <medium-button minimalistic src="svg/icon-add.svg" @click="addQuestionBlock">Add question</medium-button>
          </div>
        </div>
      </div>
      <div class="footer">
        <NuxtLink to="/">
          <regular-button @click="newGameStore.resetGame">Cancel</regular-button>
        </NuxtLink>
        <regular-button
          active
          :style="{ opacity: nextable ? '1' : '0.5' }"
          @click="nextable ? $emit('next', game) : null"
          >Next</regular-button
        >
      </div>
    </div>
  </div>
</template>

<script setup>
import { reactive, computed } from 'vue';
import { useNewGameStore } from '../stores/new';

const game = reactive({
  topic: null,
  points: null,
  time: null,
  questions: [{ id: 0, name: null, answer: null, optionA: null, optionB: null, optionC: null, optionD: null }],
});

// Load saved game from store
const newGameStore = useNewGameStore();
if (!newGameStore.isEmptyGame) {
  const { game: storedGame } = newGameStore;
  const keys = Object.keys(storedGame);
  for (const key of keys) {
    game[key] = storedGame[key];
  }
}

// Can question be removed
const removable = computed(() => game.questions.length <= 1);

// Next step is avaliable when all fields are filled
const nextable = computed(() => {
  if (!game.topic || !game.points || !game.time) return false;
  return game.questions.every(
    ({ id, name, optionA, optionB, optionC, optionD, answer }) =>
      id.toString() && name && optionA && optionB && optionC && optionD && answer
  );
});

// Triggers when question-block is updated
const handleUpdateQuestions = (data) => {
  const { id } = data;
  if (game.questions.length <= 1) return (game.questions[0] = data);
  for (let i = 0; i < game.questions.length; i++) {
    const element = game.questions[i];
    if (element.id === id) {
      game.questions[i] = data;
    }
  }
};

const addQuestionBlock = () => {
  // Finds last element id, creates new block with increased last id
  const { length } = game.questions;
  const { id: lastId } = game.questions[length - 1];
  game.questions.push({ id: lastId + 1 });
};

const removeQuestionBlock = (id) => {
  game.questions = game.questions.filter((q) => q.id !== id);
};
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

.main {
  margin-bottom: 30px;
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

.about {
  display: grid;
  grid-template-columns: 2fr 1fr 1fr;
  grid-template-rows: auto auto;
  gap: 10px;
  margin-bottom: 30px;
}

.questions-list-enter-active,
.questions-list-leave-active {
  transition: all 0.5s ease;
}
.questions-list-enter-from,
.questions-list-leave-to {
  opacity: 0;
  transform: translateX(30px);
}

.about__title,
.questions__title {
  font-family: 'Inter-SemiBold', sans-serif;
  color: var(--font-primary-color);
  font-size: var(--font-tertiary-size);
  line-height: 20px;
}

.questions__title {
  margin-bottom: 20px;
}

.question {
  margin-bottom: 20px;
}

.questions__add {
  display: flex;
  justify-content: center;
}

.footer {
  display: flex;
  gap: 10px;
}
</style>
