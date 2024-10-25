import './assets/main.css';

import { createApp } from 'vue';
import { createPinia } from 'pinia';

import App from './App.vue';
import router from './router';
import { WorkerPeer } from './lib/peer';
import { Event } from './lib/messages';

const app = createApp(App);

app.use(createPinia());
app.use(router);

const worker = new Worker(new URL('@/worker/index.ts', import.meta.url), { type: 'module' });
const workerPeer = new WorkerPeer(worker);

workerPeer.on(Event.WASMLoaded, () => {
  app.mount('#app');
});

workerPeer.loadWASM(new URL('webui.wasm', document.baseURI).toString());

app.provide<WorkerPeer>('workerPeer', workerPeer);
