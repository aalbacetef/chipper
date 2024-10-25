<script setup lang="ts">
import { inject, ref } from 'vue';
import { WorkerPeer } from '@/lib/peer';
import { loadROMManifest, type ROMEntry, mapKeyToHex, hexToRGBA, defaultColors } from '@/lib/game';
import { loadAudioManifest, type AudioManifest } from '@/lib/music';

import AudioPlayer from '@/components/AudioPlayer.vue';
import DrawArea from '@/components/DrawArea.vue';
import ColorPicker from '@/components/ColorPicker.vue';

const roms = ref<ROMEntry[]>([]);
const loading = ref<boolean>(true);
const audioManifest = ref<AudioManifest>();

let colorSet = defaultColors.set;
let colorClear = defaultColors.clear;

loadAudioManifest()
  .then((m) => (audioManifest.value = m))
  .catch((err) => console.error('failed to load audio manifest: ', err));

loadROMManifest()
  .then((data) => {
    roms.value = data;
    loading.value = false;
  })
  .catch((err) => console.error('failed to load rom manifest: ', err));

const selectedRomIndex = ref<number>(0);

const workerPeer = inject<WorkerPeer>('workerPeer');

function handleLoadROMButton() {
  const rom = roms.value[selectedRomIndex.value];
  const romURL = new URL(rom.path, document.baseURI);
  workerPeer!.loadROM(romURL.toString());
}

function handleStartButton() {
  workerPeer!.startEmu();
}

function handleStopButton() {
  workerPeer!.stopEmu();
}

function handleRestartButton() {
  workerPeer!.restartEmu();
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

function updateColor(args: [string, string]): void {
  const [name, colorHex] = args;

  const color = hexToRGBA(colorHex);
  switch (name) {
    case 'set':
      colorSet = color;
      break;
    case 'clear':
      colorClear = color;
      break;
  }

  workerPeer!.setColors({ set: colorSet, clear: colorClear });
}
</script>

<template>
  <main>
    <div class="loading" v-if="loading">loading</div>
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

        <div class="color-control">
          <p>pick color:</p>
          <ColorPicker name="set" display="foreground" @update="updateColor" />
          <ColorPicker name="clear" display="background" @update="updateColor" />
        </div>

        <div class="emu-control">
          <button @click="handleStartButton">Start</button>

          <button @click="handleStopButton">Stop</button>

          <button @click="handleRestartButton">Restart</button>
        </div>
      </div>

      <div
        class="game-area"
        tabindex="0"
        @keydown.prevent="handleKeyDown"
        @keyup.prevent="handleKeyUp"
      >
        <DrawArea />
      </div>
    </div>

    <AudioPlayer
      :manifest="audioManifest"
      v-if="audioManifest !== null && typeof audioManifest !== 'undefined'"
    />
  </main>
</template>

<style scoped>
main {
  width: 100%;
  height: 100%;
  position: relative;
  margin-top: 15px;
}

.control-panel {
  margin-bottom: 10px;
}

.emu-control {
  margin-top: 10px;
  display: flex;
  flex-direction: row;
}

.emu-control button {
  margin-right: 5px;
}

.color-control {
  margin: 10px 0;
}

.color-control p {
  padding: 0;
  margin: 0;
}

.game-area {
  width: 100%;
  height: 100%;
}

.game-area:focus-visible {
  outline: none;
}
</style>
