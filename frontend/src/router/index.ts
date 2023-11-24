import { createRouter, createWebHistory } from 'vue-router';
import DiscoverView from '@/views/DiscoverView.vue';
import DashboardView from '../views/DashboardView.vue';
import LibraryView from '../views/LibraryView.vue';
import StarredView from '../views/StarredView.vue';
import RecentView from '../views/RecentView.vue';
import CreateView from '../views/CreateView.vue';
import EditView from '../views/EditView.vue';
import GameView from '../views/GameView.vue';
import ConsoleView from '../views/ConsoleView.vue';
import PlayView from '../views/PlayView.vue';
import LoginView from '../views/LoginView.vue';
import JoinView from '../views/JoinView.vue';
import NotFoundView from '../views/NotFoundView.vue';
import DevComponentsView from '../views/DevComponentsView.vue';
import { useAuthStore } from '@/stores/auth';

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      redirect: '/discover',
    },
    {
      path: '/discover',
      name: 'discover',
      component: DiscoverView,
      meta: { requiresAuth: false },
    },
    {
      path: '/dashboard',
      name: 'dashboard',
      component: DashboardView,
      meta: { requiresAuth: true },
    },
    {
      path: '/library',
      name: 'library',
      component: LibraryView,
      meta: { requiresAuth: true },
    },
    {
      path: '/starred',
      name: 'starred',
      component: StarredView,
      meta: { requiresAuth: true },
    },
    {
      path: '/recent',
      name: 'recent',
      component: RecentView,
      meta: { requiresAuth: true },
    },
    {
      path: '/create',
      name: 'create',
      component: CreateView,
      meta: { requiresAuth: true },
    },
    {
      path: '/login',
      name: 'login',
      component: LoginView,
      meta: { requiresAuth: false },
    },
    {
      path: '/edit/:id',
      name: 'edit',
      component: EditView,
      meta: { requiresAuth: true },
    },
    {
      path: '/game/:id',
      name: 'game',
      component: GameView,
      meta: { requiresAuth: true },
    },
    {
      path: '/join',
      name: 'join',
      component: JoinView,
      meta: { requiresAuth: false },
    },
    {
      path: '/play/:id',
      name: 'play',
      component: PlayView,
      meta: { requiresAuth: false },
    },
    {
      path: '/console/:id',
      name: 'console',
      component: ConsoleView,
      meta: { requiresAuth: true },
    },
    {
      path: '/dev/components',
      name: 'Components Test',
      component: DevComponentsView,
      meta: { requiresAuth: false },
    },
    { path: '/:pathMatch(.*)*', name: 'NotFound', component: NotFoundView },
  ],
});

router.beforeEach((to, from, next) => {
  const { isAuthenticated } = useAuthStore();
  if (to.meta.requiresAuth && !isAuthenticated) next({ name: 'login' });
  else next();
});

export default router;
