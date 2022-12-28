<template>
  <div class="questions__wrapper">
    <div class="header">
      <div class="header__item title">
        <regular-input v-model="questions.name" minimalistic placeholder="Question name" img="svg/icon-question.svg" />
      </div>
      <div class="header__item remove">
        <medium-button
          v-if="!removable"
          minimalistic
          small-icon
          src="svg/icon-trash.svg"
          @click="$emit('remove', questions.id)"
          >Remove question</medium-button
        >
      </div>
    </div>
    <form id="form" action="">
      <ul class="questions">
        <li class="questions__field">
          <regular-input v-model="questions.optionA" placeholder="Enter an option A" img="svg/icon-triangle.svg" />
          <input id="a" v-model="questions.answer" class="qiestion__radio" value="a" type="radio" name="answer" />
        </li>
        <li class="questions__field">
          <regular-input v-model="questions.optionB" placeholder="Enter an option B" img="svg/icon-circle.svg" />
          <input id="b" v-model="questions.answer" class="qiestion__radio" value="b" type="radio" name="answer" />
        </li>
        <li class="questions__field">
          <regular-input v-model="questions.optionC" placeholder="Enter an option C" img="svg/icon-square.svg" />
          <input id="c" v-model="questions.answer" class="qiestion__radio" value="c" type="radio" name="answer" />
        </li>
        <li class="questions__field">
          <regular-input v-model="questions.optionD" placeholder="Enter an option D" img="svg/icon-diamond.svg" />
          <input id="d" v-model="questions.answer" class="qiestion__radio" value="d" type="radio" name="answer" />
        </li>
      </ul>
    </form>
  </div>
</template>

<script setup>
import { watch, reactive, defineEmits, defineProps } from 'vue';

const props = defineProps({
  removable: {
    type: Boolean,
    default: false,
  },
  generatedId: {
    type: Number,
    default: 0,
  },
  question: {
    type: Object,
    default() {
      return {
        id: 0,
        name: null,
        optionA: null,
        optionB: null,
        optionC: null,
        optionD: null,
        answer: null,
      };
    },
  },
});
const emit = defineEmits(['updateQuestions', 'remove']);

const questions = reactive({
  id: props.question.id || 0,
  name: props.question.name || null,
  optionA: props.question.optionA || null,
  optionB: props.question.optionB || null,
  optionC: props.question.optionC || null,
  optionD: props.question.optionD || null,
  answer: props.question.answer || null,
});

// Watching for fields changes and emits "updateQuestions" event with data
watch(
  () => [questions.name, questions.optionA, questions.optionB, questions.optionC, questions.optionD, questions.answer],
  ([name, optionA, optionB, optionC, optionD, answer]) => {
    emit('updateQuestions', { id: questions.id, name, optionA, optionB, optionC, optionD, answer });
  }
);
</script>

<style scoped>
.questions__wrapper {
  outline: 3px solid var(--background-secondary-color);
  outline-offset: -3px;
  padding: 15px;
  border-radius: 15px;
}

.header {
  display: flex;
  justify-content: space-between;
}

.questions {
  display: grid;
  grid-template-columns: 1fr 1fr;
  grid-template-columns: auto auto;
  gap: 10px;
}

.questions__field {
  position: relative;
  list-style: none;
}

.qiestion__radio {
  position: absolute;
  top: 50%;
  right: 10px;
  transform: translateY(-50%);
}

input[type='radio'] {
  -webkit-appearance: none;
  appearance: none;
  background-color: transparent;
  margin: 0;
  width: 20px;
  height: 20px;
  background-image: url('svg/icon-radio-unchecked.svg');
  background-position: center;
}

input[id='a']:checked {
  background-image: url('svg/icon-radio-red.svg');
}

input[id='b']:checked {
  background-image: url('svg/icon-radio-blue.svg');
}

input[id='c']:checked {
  background-image: url('svg/icon-radio-green.svg');
}

input[id='d']:checked {
  background-image: url('svg/icon-radio-yellow.svg');
}

.header__item {
  font-size: var(--font-quaternary-size);
  font-family: 'Inter-SemiBold', sans-serif;
  line-height: 20px;
  color: var(--font-primary-color);
  display: flex;
  align-items: center;
  gap: 5px;
  margin-bottom: 15px;
}

.img {
  width: 16px;
  height: 16px;
}

.img img {
  width: 100%;
  height: 100%;
  object-fit: contain;
}
</style>
