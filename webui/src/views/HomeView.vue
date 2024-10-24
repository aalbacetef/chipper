<script setup lang="ts">
import { inject, ref } from "vue";
import { WorkerPeer } from "@/lib/peer";
import { mapKeyToHex } from '@/lib/game';

import DrawArea from "@/components/DrawArea.vue";
import ColorPicker from "@/components/ColorPicker.vue";

const manifestURL = "/manifest.json"

type RomManifestEntry = {
  name: string;
  path: string;
}

const roms = ref<RomManifestEntry[]>([]);
const loading = ref<boolean>(true);

fetch(manifestURL)
  .then(res => res.json())
  .then((data: Record<string, string>) => {
    roms.value = [];

    Object.keys(data).forEach(name => {
      const row = { name, path: data[name] };
      roms.value.push(row);
    });

    loading.value = false;
  })

const selectedRomIndex = ref<number>(0);

const workerPeer = inject<WorkerPeer>("workerPeer");

function handleLoadROMButton() {
  const rom = roms.value[selectedRomIndex.value];
  workerPeer.loadROM("/" + rom.path);
}

function handleStartButton() {
  workerPeer.startEmu();
}

function handleStopButton() {
  workerPeer.stopEmu();
}

function handleRestartButton() {
  workerPeer.restartEmu();
}

function handleKeyDown(event) {
  event.preventDefault();
  try {
    const key = mapKeyToHex(event.code);
    workerPeer.sendKeyDown(key, event.repeat);
  } catch (err) {
    console.log(err);
  }
}

function handleKeyUp(event) {
  event.preventDefault();
  try {
    const key = mapKeyToHex(event.code);
    workerPeer.sendKeyUp(key, event.repeat);
  } catch (err) {
    console.log(err);
  }
}

function updateColor(name: string, colorHex: string) {

}
</script>

<template>
  <main>
    <div class="loading" v-if="loading">
      loading
    </div>
    <div class="view--wrapper" v-if="!loading">
      <div class="control-panel">
        <div class="rom-loader">

          <select v-model="selectedRomIndex">
            <option v-for="(rom, index) in roms" :key="rom.path" :value="index">{{ rom.name }}</option>
          </select>

          <button @click="handleLoadROMButton">
            Load ROM
          </button>
        </div>

        <div class="color-control">
          <p>pick color: </p>
          <ColorPicker name="set" @changed="updateColor" />
          <ColorPicker name="clear" @changed="updateColor" />
        </div>

        <div class="emu-control">
          <button @click="handleStartButton">
            Start
          </button>

          <button @click="handleStopButton">
            Stop
          </button>

          <button @click="handleRestartButton">
            Restart
          </button>
        </div>

      </div>

      <div class="game-area" tabindex="0" @keydown="handleKeyDown" @keyup="handleKeyUp">
        <DrawArea />
      </div>
    </div>
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
