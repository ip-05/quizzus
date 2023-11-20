import { createRouter, createWebHistory } from 'vue-router';
import DiscoverView from '../views/DiscoverView.vue';
import DashboardView from '../views/DashboardView.vue';
import LibraryView from '../views/LibraryView.vue';
import StarredView from '../views/StarredView.vue';
import RecentView from '../views/RecentView.vue';
import CreateView from '../views/CreateView.vue';
import EditView from '../views/EditView.vue';
import GameView from '../views/GameView.vue';
import ConsoleView from '../views/ConsoleView.vue';
import PlayView from '../views/PlayView.vue';
import SigninView from '../views/SigninView.vue';
import NotFoundView from '../views/NotFoundView.vue';
import DevComponentsView from '../views/DevComponentsView.vue';

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
    },
    {
      path: '/dashboard',
      name: 'dashboard',
      component: DashboardView,
    },
    {
      path: '/library',
      name: 'library',
      component: LibraryView,
    },
    {
      path: '/starred',
      name: 'starred',
      component: StarredView,
    },
    {
      path: '/recent',
      name: 'recent',
      component: RecentView,
    },
    {
      path: '/create',
      name: 'create',
      component: CreateView,
    },
    {
      path: '/edit/:id',
      name: 'edit',
      component: EditView,
    },
    {
      path: '/game/:id',
      name: 'game',
      component: GameView,
    },
    {
      path: '/console/:id',
      name: 'console',
      component: ConsoleView,
    },
    {
      path: '/play/:id',
      name: 'play',
      component: PlayView,
    },
    {
      path: '/signin',
      name: 'signin',
      component: SigninView,
    },
    {
      path: '/dev/components',
      name: 'Components Test',
      component: DevComponentsView,
    },
    { path: '/:pathMatch(.*)*', name: 'NotFound', component: NotFoundView },
  ],
});

export default router;
