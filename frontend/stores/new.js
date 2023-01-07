import axios from 'axios';
import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import { useErrorStore } from './error';
import { useRuntimeConfig, useCookie } from '#imports';

export const useNewGameStore = defineStore('NewGameStore', () => {
  const config = useRuntimeConfig();
  const errorStore = useErrorStore();
  const tokenCookie = useCookie('token');

  const gameId = ref(0);
  const inviteCode = ref('');

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

  async function postGame() {
    try {
      const { data } = await axios.post(
        '/games',
        {
          topic: topic.value,
          roundTime: parseInt(time.value),
          points: parseInt(points.value),
          questions: questions.value,
        },
        {
          baseURL: config.public.apiUrl,
          headers: {
            Authorization: 'Bearer ' + tokenCookie.value,
          },
        }
      );
      gameId.value = data.id;
      inviteCode.value = data.inviteCode;
      topic.value = data.topic;
      points.value = data.points;
      time.value = data.roundTime;
      questions.value = data.questions;
    } catch (error) {
      const message = error.response.data.error;
      errorStore.message = message;
    }
  }

  async function getGame(query) {
    try {
      const { data } = await axios.get(`/games`, {
        baseURL: config.public.apiUrl,
        headers: {
          Authorization: 'Bearer ' + tokenCookie.value,
        },
        params: query,
      });
      gameId.value = data.id;
      inviteCode.value = data.inviteCode;
      topic.value = data.topic;
      points.value = data.points;
      time.value = data.roundTime;
      questions.value = data.questions;
    } catch (error) {
      const message = error.response.data.error;
      errorStore.message = message;
    }
  }

  async function updateGame(query) {
    try {
      const { data } = await axios.patch(
        '/games',
        {
          topic: topic.value,
          roundTime: parseInt(time.value),
          points: parseInt(points.value),
          questions: questions.value,
        },
        {
          baseURL: config.public.apiUrl,
          headers: {
            Authorization: 'Bearer ' + tokenCookie.value,
          },
          params: query,
        }
      );
      gameId.value = data.id;
      inviteCode.value = data.inviteCode;
      topic.value = data.topic;
      points.value = data.points;
      time.value = data.roundTime;
      questions.value = data.questions;
    } catch (error) {
      const message = error.response.data.error;
      errorStore.message = message;
    }
  }

  // User exits and game resets
  function resetGame() {
    topic.value = '';
    points.value = null;
    time.value = null;
    questions.value = [];
    appendQuestion();
  }

  return {
    gameId,
    topic,
    inviteCode,
    time,
    points,
    questions,
    postGame,
    nextable,
    resetGame,
    appendQuestion,
    removeQuestion,
    getGame,
    updateGame,
  };
});
