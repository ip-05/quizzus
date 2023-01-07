<template>
  <div class="questions__wrapper">
    <div class="header">
      <div class="header__item title">
        <regular-input v-model="question.name" minimalistic placeholder="Question name" img="svg/icon-question.svg" />
      </div>
      <div class="header__item remove">
        <medium-button v-if="!removable" minimalistic small-icon src="svg/icon-trash.svg" @click="removeQuestion"
          >Remove question</medium-button
        >
      </div>
    </div>
    <form id="form" action="">
      <ul class="questions">
        <li class="questions__field">
          <regular-input
            v-model="question.options[0].name"
            placeholder="Enter an option A"
            img="svg/icon-triangle.svg"
          />
          <input id="a" v-model="correctOption" :value="0" class="qiestion__radio" type="radio" name="answer" />
        </li>
        <li class="questions__field">
          <regular-input v-model="question.options[1].name" placeholder="Enter an option B" img="svg/icon-circle.svg" />
          <input id="b" v-model="correctOption" :value="1" class="qiestion__radio" type="radio" name="answer" />
        </li>
        <li class="questions__field">
          <regular-input v-model="question.options[2].name" placeholder="Enter an option C" img="svg/icon-square.svg" />
          <input id="c" v-model="correctOption" :value="2" class="qiestion__radio" type="radio" name="answer" />
        </li>
        <li class="questions__field">
          <regular-input
            v-model="question.options[3].name"
            placeholder="Enter an option D"
            img="svg/icon-diamond.svg"
          />
          <input id="d" v-model="correctOption" :value="3" class="qiestion__radio" type="radio" name="answer" />
        </li>
      </ul>
    </form>
  </div>
</template>

<script setup>
import { defineProps, ref, watch, computed, onMounted } from 'vue';
import { storeToRefs } from 'pinia';
import { useNewGameStore } from '../stores/new';

const newGameStore = useNewGameStore();

const props = defineProps({
  generatedId: {
    type: Number,
    default: 0,
  },
});

const { questions } = storeToRefs(newGameStore);
const removable = computed(() => questions.value.length <= 1);

const question = computed(() => questions.value.find((q) => q.id === props.generatedId));
const correctOption = ref(null);

onMounted(() => {
  for (let i = 0; i < 4; i++) {
    if (question.value.options[i].correct) {
      correctOption.value = i;
    }
  }
});

// on radio button changes, correct sets to true
watch(correctOption, (answer) => {
  for (let i = 0; i < 4; i++) {
    const checked = answer === i;
    question.value.options[i].correct = checked;
  }
});

const removeQuestion = () => {
  newGameStore.removeQuestion(props.generatedId);
};
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
  background-image: url('/svg/icon-radio-unchecked.svg');
  background-position: center;
}

input[id='a']:checked {
  background-image: url('/svg/icon-radio-red.svg');
}

input[id='b']:checked {
  background-image: url('/svg/icon-radio-blue.svg');
}

input[id='c']:checked {
  background-image: url('/svg/icon-radio-green.svg');
}

input[id='d']:checked {
  background-image: url('/svg/icon-radio-yellow.svg');
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
