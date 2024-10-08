import './assets/main.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import router from './router'
import { initialize } from './lib/wasm';

const app = createApp(App)

app.use(createPinia())
app.use(router)

initialize("webui.wasm")
  .then(r => {
    app.provide("wasmLoad", r);
    app.mount('#app');
  });
