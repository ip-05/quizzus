'use strict';

import { useSocketStore } from '../stores/socket';
import { defineNuxtRouteMiddleware } from '#imports';

export default defineNuxtRouteMiddleware((to, from) => {
  if (process.client) {
    if (from.path.startsWith('/game/') || from.path.startsWith('/console/')) {
      const socketStore = useSocketStore();
      socketStore.leaveGame();
    }
  }
});
