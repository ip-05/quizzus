import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import { useCookie } from '#imports';

export const useNewGameStore = defineStore('NewGameStore', () => {
  const tokenCookie = useCookie('token');

  const topic = ref('');
  const points = ref(null);
  const time = ref(null);

  const id = ref(0);
  const questions = ref([]);
  appendQuestion();

  function appendQuestion() {
    id.value += 1;
    questions.value.push({
      id: id.value,
      name: '',
      options: [
        {
          name: '',
          correct: false,
        },
        {
          name: '',
          correct: false,
        },
        {
          name: '',
          correct: false,
        },
        {
          name: '',
          correct: false,
        },
      ],
    });
  }

  function removeQuestion(id) {
    questions.value = questions.value.filter((q) => q.id !== id);
  }

  // Whether all necessary fields are filled
  const nextable = computed(() => {
    // If topic or points or time is not filled
    if (!topic.value.length || !points.value || !time.value) return false;

    // If there is no questions
    if (!questions.value.length) return false;

    const names = [];
    const checks = [];
    // Collects all text field and checks
    for (const { name, options } of questions.value) {
      names.push(name);
      for (const { name, correct } of options) {
        names.push(name);
        checks.push(correct);
      }
    }

    const questionsChecked = questions.value.length === checks.filter((c) => c === true).length;
    const textsFilled = names.every((n) => n.length > 0);
    return questionsChecked && textsFilled;
  });

  function postGame() {}

  // User exits and game resets
  function resetGame() {
    topic.value = '';
    points.value = null;
    time.value = null;
    questions.value = [];
    appendQuestion();
  }

  return { topic, time, points, questions, postGame, nextable, resetGame, appendQuestion, removeQuestion };
});
