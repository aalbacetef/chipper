import './assets/main.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import router from './router'
import { WorkerPeer } from './lib/peer'
import { Event } from './lib/messages'

const app = createApp(App)

app.use(createPinia())
app.use(router)

const url = new URL('@/worker/index.ts', import.meta.url);
const workerPeer = new WorkerPeer(new Worker(url, { type: "module" }));

workerPeer.on(Event.WASMLoaded, () => {
  app.mount("#app");
});
workerPeer.loadWASM("/webui.wasm");


app.provide<WorkerPeer>("workerPeer", workerPeer);

