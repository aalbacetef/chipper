import './assets/main.css';

import { createApp } from 'vue';
import { createPinia } from 'pinia';

// @ts-ignore
import App from './App.vue';
import router from './router';
import { WorkerPeer } from './lib/peer';
import { Event } from './lib/messages';
import { AudioState, useAppStore } from './stores/app';

const app = createApp(App);

const worker = new Worker(new URL('@/worker/index.ts', import.meta.url), { type: 'module' });
const workerPeer = new WorkerPeer(worker);

workerPeer.on(Event.WASMLoaded, () => {
  app.use(createPinia());
  app.use(router);

  app.provide<WorkerPeer>('workerPeer', workerPeer);

  app.mount('#app');
});

workerPeer.loadWASM(new URL('webui.wasm', document.baseURI).toString());


// play audio on first interaction
document.addEventListener('DOMContentLoaded', () => {
  document.addEventListener('click', () => {
    const audio = document.querySelector('audio');
    if (audio === null) {
      return;
    }

    setTimeout(() => {
      const appStore = useAppStore();
      appStore.setAudioState(AudioState.Playing);
    }, 150);

    audio.play();
  }, { once: true });
})
