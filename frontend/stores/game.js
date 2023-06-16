import axios from 'axios';
import { defineStore } from 'pinia';
import { ref, reactive, computed, watch } from 'vue';
import { useErrorStore } from './error';
import { useAuthStore } from './auth';
import { useSocketStore } from './socket';
import { useRuntimeConfig, navigateTo } from '#imports';

/**
 * game-wait
 * game-wait-admin
 * game-in
 * game-over
 * game-over-admin
 * game-correct
 * game-wrong
 * game-graph-admin
 */

export const useGameStore = defineStore('GameStore', () => {
  const config = useRuntimeConfig();
  const errorStore = useErrorStore();

  const authStore = useAuthStore();

  const socketStore = useSocketStore();

  // Game Mode
  const state = ref('game-wait');
  const active = ref(false);

  // Game Info
  const gameId = ref(0);
  const inviteCode = ref('');
  const owner = ref(null);
  const isOwner = computed(() => authStore.user.id === owner.value);

  const topic = ref('' || 'Game topic');
  const totalPoints = ref(0);
  const points = ref(null || 10);
  const roundTime = ref(null || 20);
  const timer = ref(null); // timer for dynamic island for all users synchronously
  const participants = ref([]);
  const participantsNumber = computed(() => participants.value.length);

  const gameStarted = ref(false);
  const countdown = ref(null);
  const canContinue = ref(true); // can admin start game or give another question (depens on game state)

  const id = ref(63457378356); // Local questions id for Vue to render properly
  const questions = ref([]);
  // appendQuestion();
  const questionsNumber = computed({
    get() {
      return questions.value.length;
    },
    set(count) {
      questions.value = new Array(count);
    },
  });

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
            Authorization: 'Bearer ' + authStore.token,
          },
        }
      );
      fillGameFields(data);
    } catch (error) {
      const message = error.response.data.error;
      errorStore.message = message;
    }
  }

  async function getGame(id) {
    try {
      const { data } = await axios.get(`/games/${id}`, {
        baseURL: config.public.apiUrl,
        headers: {
          Authorization: 'Bearer ' + authStore.token,
        }
      });
      // for regular user to join
      if (data.message === 'Game found' && data.topic) {
        inviteCode.value = query.invite_code;
        return 'player';
      }
      fillGameFields(data);
      return 'admin';
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
            Authorization: 'Bearer ' + authStore.token,
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
    owner.value = data.ownerId;
    console.log(questions.value);
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
    answer: null,
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

  const selectedOption = ref(null);
  const submittedAnswersByUsers = ref({});
  /**
   * userId: idx
   * userId2: idx2
   */

  const submittedAnswers = computed(() => {
    const answers = [
      { id: 0, answers: 0 },
      { id: 1, answers: 0 },
      { id: 2, answers: 0 },
      { id: 3, answers: 0 },
    ];

    for (const userId of Object.keys(submittedAnswersByUsers.value)) {
      const option = submittedAnswersByUsers.value[userId];
      answers[option].answers += 1;
    }

    return answers;
  });

  const graphAnswers = computed(() => {
    const answers = submittedAnswers.value.map(({ answers }) => answers);
    const max = Math.max(...answers);
    const min = 10; // default minimum if every number is 0
    const percents = answers.map((n, i, arr) => (n / max) * 100 || min);
    return percents;
  });

  const leaderboardShown = ref(false);
  const leaderboard = reactive([]);
  const leaderboardSorted = computed(() =>
    leaderboard.sort((a, b) => b.points - a.points || a.participant.localeCompare(b.participant))
  );
  const top3 = computed(() => {
    const { length } = leaderboardSorted.value;
    if (!length) return false;
    return leaderboardSorted.value.slice(0, 3);
  });
  const winner = computed(() => {
    const { length } = leaderboardSorted.value;
    return length ? leaderboardSorted.value[0] : false;
  });

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
    state.value = 'game-over';
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

  /**
   * WebSockets event
   * ----------------
   */

  // user sends event when joins the game
  function joinedGame(data) {
    console.log('test', data);
    // for admin
    console.log('isOwner: ', isOwner.value);
    if (isOwner.value) {
      setGameWaitAdmin();
      return;
    }
    active.value = true;
    // for regular players
    setGameWait();
    points.value = data.points;
    roundTime.value = data.roundTime;
    topic.value = data.topic;
    questionsNumber.value = data.questionCount;
    const { owner, members } = data;
    for (const member of Object.keys(members)) {
      if (member !== owner.id) {
        const { id, name, profilePicture } = members[member];
        participants.value.push({ id, name, profilePicture });
      }
    }
  }

  // recieved when another player joins game
  function userJoined(data) {
    const { id, name, profilePicture } = data;
    participants.value.push({ id, name, profilePicture });
  }

  // admin sends event
  function startGame() {
    active.value = true;
    gameStarted.value = true;
    canContinue.value = false;
    socketStore.send({
      message: 'START_GAME',
    });
  }

  // player sends when leaves
  function userLeft(data) {
    participants.value = participants.value.filter(({ profilePicture }) => data.profilePicture !== profilePicture);
  }

  // admin starts game and countdown starts for every user
  function gameStarting(data) {
    countdown.value = data;
  }

  // event right after countdown ending
  function gameInProgress() {
    countdown.value = null;
  }

  function roundInprogress(data) {
    if (data) {
      currentQuestion.id = data.question.id;
      currentQuestion.name = data.question.name;
      currentQuestion.options = data.question.options;
      currentQuestion.answer = null;
      timer.value = data.timer;
    }
    // for admin displays graph of answers
    if (isOwner.value) {
      setGameGraphAdmin();
      return;
    }
    // for players displays currect question
    setGameIn();
  }

  function userAnswered(data) {
    submittedAnswersByUsers.value[data.user] = data.option;
  }

  // player gets info about their aswers
  function roundFinished(data) {
    const answer = data.options.filter(({ correct }) => correct)[0];
    currentQuestion.answer = answer.name;
    currentQuestion.options = data.options;
    correctAnswer.value = answer;
    selectedOption.value = null;
    totalPoints.value = data.leaderboard[authStore.user.id];
    submittedAnswersByUsers.value = {};
    // displaying default wait screen for admin
    if (isOwner.value) {
      setGameWaitAdmin();
      canContinue.value = true;
      return;
    }
    // displaying correct screen if player is correct
    if (data.correct) {
      setGameCorrect();
      return;
    }
    // displaying correct screen if player is wrong
    setGameWrong();
  }

  // admin gives another round
  function nextQuestion() {
    socketStore.send({
      message: 'NEXT_ROUND',
    });
    canContinue.value = false;
  }

  // player submits answer
  watch(selectedOption, (s) => {
    answerQuestion(s);
  });

  function answerQuestion(data) {
    socketStore.send({
      message: 'ANSWER_QUESTION',
      data: {
        option: data,
      },
    });
  }

  function gameFinished(data) {
    gameStarted.value = false;
    canContinue.value = false;
    Object.assign(
      leaderboard,
      participants.value.map(({ id, name }) => ({ id, participant: name, points: data[id] }))
    );
    if (isOwner.value) {
      setGameOverAdmin();
      return;
    }
    setGameOverClient();
  }

  // admin of the game left - all users go home
  function hostLeft() {
    active.value = false;
    navigateTo('/');
  }

  // if users tries to get in before admin
  function gameForbidden() {
    active.value = false;
    navigateTo('/');
  }

  return {
    active,
    state,
    topic,
    totalPoints,
    inviteCode,
    countdown,
    points,
    roundTime,
    timer,
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
    selectedOption,
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
    isOwner,
    //
    gameStarted,
    canContinue,
    joinedGame,
    userJoined,
    userLeft,
    startGame,
    gameStarting,
    gameInProgress,
    roundInprogress,
    roundFinished,
    answerQuestion,
    userAnswered,
    hostLeft,
    gameForbidden,
    nextQuestion,
    gameFinished,
  };
});
