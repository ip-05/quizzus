'use strict';

import { defineNuxtRouteMiddleware, useCookie, navigateTo } from '#imports';

export default defineNuxtRouteMiddleware((to) => {
  const { path } = to;
  if (path === '/' || path === '/auth/google') return;

  // only authorized users
  // const tokenCookie = useCookie('token');
  if (!localStorage.getItem('token')) {
    return navigateTo('/');
  }
});
