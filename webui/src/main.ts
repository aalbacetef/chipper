import './assets/main.css';

import { createApp } from 'vue';
import { createPinia } from 'pinia';

import App from './App.vue';
import router from './router';
import { WorkerPeer } from './lib/peer';
import { Event } from './lib/messages';

const app = createApp(App);


const worker = new Worker(new URL('@/worker/index.ts', import.meta.url), { type: 'module' });
const workerPeer = new WorkerPeer(worker);

workerPeer.on(Event.WASMLoaded, () => {
  app.use(createPinia());
  app.use(router);
  app.provide<WorkerPeer>('workerPeer', workerPeer);
  app.mount('#app');
});

const baseURL = window.location.origin + window.location.pathname;
workerPeer.loadWASM(new URL('webui.wasm', baseURL).toString());

