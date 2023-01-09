import axios from 'axios';
import { defineStore } from 'pinia';
import { ref, reactive, computed } from 'vue';
import { useErrorStore } from './error';
import { useRuntimeConfig, useCookie } from '#imports';

/**
 * game-wait
 * game-wait-admin
 * game-in
 * game-over-client
 * game-over-admin
 * game-correct
 * game-wrong
 * game-graph-admin
 */

export const useGameStore = defineStore('GameStore', () => {
  const config = useRuntimeConfig();
  const errorStore = useErrorStore();
  const tokenCookie = useCookie('token');

  // Game Mode
  const state = ref('game-graph-admin');

  // Game Info
  const gameId = ref(0);
  const inviteCode = ref('');

  const topic = ref('');
  const totalPoints = ref(0);
  const points = ref(null);
  const roundTime = ref(null);
  const participants = ref([]);
  const participantsNumber = computed(() => participants.value.length);

  const id = ref(0); // Local questions id for Vue to render properly
  const questions = ref([]);
  appendQuestion();
  const questionsNumber = computed(() => questions.value.length);

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
    // If topic or points or roundTime is not filled
    if (!topic.value.length || !points.value || !roundTime.value) return false;

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
          roundTime: parseInt(roundTime.value),
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
      fillGameFields(data);
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
      fillGameFields(data);
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
          roundTime: parseInt(roundTime.value),
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
      fillGameFields(data);
    } catch (error) {
      const message = error.response.data.error;
      errorStore.message = message;
    }
  }

  function fillGameFields(data) {
    gameId.value = data.id;
    inviteCode.value = data.inviteCode;
    topic.value = data.topic;
    points.value = data.points;
    roundTime.value = data.roundTime;
    questions.value = data.questions;
  }

  // User exits and game resets
  function resetGame() {
    topic.value = '';
    points.value = null;
    roundTime.value = null;
    questions.value = [];
    appendQuestion();
  }

  // Current question in game
  const currentQuestion = reactive({
    name: 'Question name',
    id: 0,
    options: [
      {
        name: 'Option A',
        correct: false,
        id: 0,
      },
      {
        name: 'Option B',
        correct: false,
        id: 1,
      },
      {
        name: 'Option C',
        correct: false,
        id: 2,
      },
      {
        name: 'Option D',
        correct: false,
        id: 3,
      },
    ],
  });
  // After answering question
  const correctAnswer = ref({
    name: 'Option C',
    id: 2,
  });

  const submittedAnswers = ref([
    { id: 0, answers: 8 },
    { id: 1, answers: 6 },
    { id: 2, answers: 15 },
    { id: 3, answers: 0 },
  ]);

  const graphAnswers = computed(() => {
    const answers = submittedAnswers.value.map(({ answers }) => answers);
    const max = Math.max(...answers);
    const min = 10; // default minimum if every number is 0
    const percents = answers.map((n, i, arr) => (n / max) * 100 || min);
    return percents;
  });

  const leaderboardShown = ref(false);
  const leaderboard = reactive([
    {
      points: 175,
      participant: 'Joe Biden',
    },
    {
      points: 150,
      participant: 'Sussy baka',
    },
    {
      points: 200,
      participant: 'Amogus',
    },
    {
      points: 30,
      participant: 'zaza',
    },
  ]);
  const leaderboardSorted = computed(() => leaderboard.sort((a, b) => b.points - a.points));
  const top3 = computed(() => leaderboardSorted.value.slice(0, 3));
  const winner = computed(() => leaderboardSorted.value[0]);

  // Hint in the bottom of the game display
  const hint = computed(() => {
    const hints = {
      'game-wait': 'Waiting game to start',
      'game-wait-admin': 'Waiting game to start',
      'game-in': 'One correct answer',
      'game-over-client': null,
      'game-over-admin': null,
      'game-correct': 'You answer is correct',
      'game-wrong': 'Oops, you are mistaken',
      'game-graph-admin': currentQuestion.name,
    };
    return hints[state.value];
  });

  // Game mode setters
  function setGameWait() {
    state.value = 'game-wait';
  }
  function setGameWaitAdmin() {
    state.value = 'game-wait-admin';
  }
  function setGameIn() {
    state.value = 'game-in';
  }
  function setGameOverClient() {
    state.value = 'game-over-client';
  }
  function setGameOverAdmin() {
    state.value = 'game-over-admin';
  }
  function setGameCorrect() {
    state.value = 'game-correct';
  }
  function setGameWrong() {
    state.value = 'game-wrong';
  }
  function setGameGraphAdmin() {
    state.value = 'game-graph-admin';
  }

  return {
    state,
    topic,
    totalPoints,
    inviteCode,
    points,
    roundTime,
    questions,
    questionsNumber,
    removeQuestion,
    appendQuestion,
    participants,
    participantsNumber,
    currentQuestion,
    correctAnswer,
    leaderboardShown,
    leaderboard,
    leaderboardSorted,
    top3,
    submittedAnswers,
    graphAnswers,
    winner,
    hint,
    nextable,
    getGame,
    postGame,
    updateGame,
    resetGame,
    setGameWait,
    setGameWaitAdmin,
    setGameIn,
    setGameOverClient,
    setGameOverAdmin,
    setGameCorrect,
    setGameWrong,
    setGameGraphAdmin,
  };
});
