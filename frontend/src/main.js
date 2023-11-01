import './assets/main.css';

import { createApp } from 'vue';
import { createPinia } from 'pinia';

import App from './App.vue';
import router from './router';
import VueSelect from 'vue-select';

const app = createApp(App);

app.component('v-select', VueSelect);
app.use(createPinia());
app.use(router);

app.mount('#app');
