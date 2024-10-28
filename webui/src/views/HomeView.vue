<script setup lang="ts">
import { inject, ref } from 'vue';
import { WorkerPeer } from '@/lib/peer';
import { loadROMManifest, type ROMEntry, mapKeyToHex } from '@/lib/game';
import { loadAudioManifest, type AudioManifest } from '@/lib/music';

import AudioPlayer from '@/components/audio-player.vue';
import DrawArea from '@/components/draw-area.vue';
import ColorPicker from '@/components/color-picker.vue';
import { useAppStore, type Buttons } from '@/stores/app';

const workerPeer = inject<WorkerPeer>('workerPeer');
const appStore = useAppStore();

const roms = ref<ROMEntry[]>([]);
const loading = ref<boolean>(true);
const audioManifest = ref<AudioManifest>();
const selectedRomIndex = ref<number>(0);
const tickPeriod = ref<number>(appStore.tickPeriodMilliseconds);

loadAudioManifest()
  .then((m) => (audioManifest.value = m))
  .catch((err) => console.error('failed to load audio manifest: ', err));

loadROMManifest()
  .then((data) => {
    roms.value = data;
    loading.value = false;
  })
  .catch((err) => console.error('failed to load rom manifest: ', err));

function handleLoadROMButton() {
  const rom = roms.value[selectedRomIndex.value];
  const romURL = new URL(rom.path, document.baseURI);
  workerPeer!.loadROM(romURL.toString());
}

function handleButton(which: Buttons): void {
  appStore.buttonClicked(which);
}

function handleKeyDown(event: KeyboardEvent) {
  try {
    const key = mapKeyToHex(event.code);
    workerPeer!.sendKeyDown(key, event.repeat);
  } catch (err) {
    console.log(err);
  }
}

function handleKeyUp(event: KeyboardEvent) {
  try {
    const key = mapKeyToHex(event.code);
    workerPeer!.sendKeyUp(key, event.repeat);
  } catch (err) {
    console.log(err);
  }
}

function updateTickPeriod() {
  appStore.setTickPeriod(tickPeriod.value);
}
</script>

<template>
  <main>
    <div class="loading" v-if="loading">loading...</div>
    <div class="view--wrapper" v-if="!loading">
      <div class="control-panel">
        <div class="rom-loader">
          <select v-model="selectedRomIndex">
            <option v-for="(rom, index) in roms" :key="rom.path" :value="index">
              {{ rom.name }}
            </option>
          </select>

          <button @click="handleLoadROMButton">Load ROM</button>
        </div>

        <div class="emulator-settings">
          <div class="color-control">
            <p>pick color:</p>
            <color-picker name="set" display="foreground" />
            <color-picker name="clear" display="background" />
          </div>
          <div class="advanced">
            <label>
              <p>Tick period (in milliseconds)</p>
              <input
                type="number"
                v-model="tickPeriod"
                @change="updateTickPeriod"
                min="1"
                step="1"
              />
            </label>
          </div>
        </div>

        <div class="state-control">
          <button @click="() => handleButton('start')">Start</button>
          <button @click="() => handleButton('stop')">Stop</button>
          <button @click="() => handleButton('restart')">Restart</button>
        </div>
      </div>

      <div
        class="game-area"
        tabindex="0"
        @keydown.prevent="handleKeyDown"
        @keyup.prevent="handleKeyUp"
      >
        <draw-area />
      </div>
    </div>

    <audio-player
      :manifest="audioManifest"
      v-if="audioManifest !== null && typeof audioManifest !== 'undefined'"
    />
  </main>
</template>

<style scoped>
main {
  width: 100%;
  position: relative;
  margin-top: 15px;
}

.control-panel {
  margin-bottom: 10px;
}

.emulator-settings {
  display: flex;
  flex-direction: row;
}

.state-control {
  margin-top: 10px;
  display: flex;
  flex-direction: row;
}

.state-control button {
  margin-right: 5px;
}

.color-control {
  margin: 10px 0;
}

.color-control p {
  padding: 0;
  margin: 0;
}

.advanced {
  margin-left: 10px;
}

.game-area {
  width: 100%;
  height: 100%;
}

.game-area:focus-visible {
  outline: none;
}
</style>
