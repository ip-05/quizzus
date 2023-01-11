import { defineStore } from 'pinia';
import { ref } from 'vue';
import { useGameStore } from './game';
import { useCookie } from '#imports';

export const useSocketStore = defineStore('SocketStore', () => {
  const tokenCookie = useCookie('token');

  const gameStore = useGameStore();

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

  const socket = new WebSocket(`ws://localhost:3001/ws?token=${tokenCookie.value}`);

  socket.addEventListener('open', () => {
    console.log('Socket connected', gameStore.inviteCode);
    // joinGame(gameStore.inviteCode);
  });

  socket.addEventListener('message', (m) => {
    const { message, error, data } = JSON.parse(m.data);
    console.log('message', message, 'error', error, 'data', data);
    fns[message](data);
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

  return { send, joinGame };
});
