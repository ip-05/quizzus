import { defineStore } from 'pinia';
import { watch } from 'vue';
import { useGameStore } from './game';
import { useAuthStore } from './auth';
import { useCookie, useRoute, useRuntimeConfig } from '#imports';

export const useSocketStore = defineStore('SocketStore', () => {
  const config = useRuntimeConfig();
  // const tokenCookie = useCookie('token');

  const authStore = useAuthStore();
  const gameStore = useGameStore();
  // const route = useRoute();
  // watch(
  //   route,
  //   (to) => {
  //     console.log(to);
  //     // showExtra.value = !!to.meta?.showExtra;??
  //   },
  //   { flush: 'pre', immediate: true, deep: true }
  // );

  const fns = {
    JOINED_GAME: gameStore.joinedGame,
    USER_JOINED: gameStore.userJoined,
    USER_LEFT: gameStore.userLeft,
    GAME_STARTING: gameStore.gameStarting,
    GAME_IN_PROGRESS: gameStore.gameInProgress,
    ROUND_IN_PROGRESS: gameStore.roundInprogress,
    ROUND_FINISHED: gameStore.roundFinished,
    ANSWER_QUESTION: gameStore.answerQuestion,
    GAME_DELETED: gameStore.hostLeft,
    NOT_OWNER: gameStore.gameForbidden,
    USER_ANSWERED: gameStore.userAnswered,
    GAME_FINISHED: gameStore.gameFinished,
    ROUND_WAITING: () => {},
    ANSWER_ACCEPTED: () => {},
  };

  const socket = new WebSocket(`${config.public.socketUrl}?token=${localStorage.getItem('token')}`);

  socket.addEventListener('open', () => {
    console.log('Socket connected', gameStore.inviteCode);
    setInterval(() => {
      socket.send(JSON.stringify({ message: 'PING' }));
    }, 15000);
  });

  socket.addEventListener('message', (m) => {
    const { message, error, data } = JSON.parse(m.data);
    console.log('message', message, 'error', error, 'data', data);
    // console.log(typeof fns[message], message, fns);
    if (fns[message]) {
      fns[message](data);
    }
  });

  socket.addEventListener('close', () => {
    console.log('Socket disconnected');
  });

  function send(d) {
    const data = JSON.stringify(d);
    if (socket.readyState === WebSocket.OPEN) {
      console.log('send in readystate open: ', data);
      socket.send(data);
      return;
    }

    socket.addEventListener(
      'open',
      () => {
        console.log('sent in listener in "send": ', data);
        socket.send(data);
      },
      { once: true }
    );
  }

  function joinGame(inviteCode) {
    send({
      message: 'JOIN_GAME',
      data: {
        gameId: inviteCode,
      },
    });
  }

  function leaveGame() {
    send({
      message: 'LEAVE_GAME',
    });
  }

  return { send, joinGame, leaveGame };
});
