<script setup lang="ts">
import { inject, ref } from "vue";
import { WorkerPeer } from "@/lib/peer";
import { mapKeyToHex } from '@/lib/game';
import DrawArea from "@/components/DrawArea.vue";

const roms = [
  { name: "IBM Logo", path: "/old-roms/ibm-logo.ch8" },
  { name: "Lunar Landar", path: "/old-roms/lunar-landar.ch8" },
  { name: "Maze", path: "/old-roms/maze-alt-david-winter-199x.ch8" },
  { name: "Particle Demo", path: "/old-roms/particle-demo-zero-2008.ch8" },
  { name: "Pong", path: "/old-roms/pong.ch8" },
  { name: "Trip8", path: "/old-roms/trip8-demo.ch8" },
  { name: "Zero", path: "/old-roms/zero-demo-2007.ch8" },
  { name: "Keypad Test", path: "/old-roms/6-keypad.ch8" },
  { name: "Space Invaders", path: "/old-roms/Space Invaders [David Winter].ch8" },
  { name: "Airplane", path: "/old-roms/Airplane.ch8" },

  { name: "15puzzle", path: "/roms/15puzzle.ch8" },
  { name: "blinky", path: "/roms/blinky.ch8" },
  { name: "blitz", path: "/roms/blitz.ch8" },
  { name: "breakout", path: "/roms/breakout.ch8" },
  { name: "brix", path: "/roms/brix.ch8" },
  { name: "connect4", path: "/roms/connect4.ch8" },
  { name: "guess", path: "/roms/guess.ch8" },
  { name: "hidden", path: "/roms/hidden.ch8" },
  { name: "invaders", path: "/roms/invaders.ch8" },
  { name: "kaleid", path: "/roms/kaleid.ch8" },
  { name: "maze", path: "/roms/maze.ch8" },
  { name: "merlin", path: "/roms/merlin.ch8" },
  { name: "missile", path: "/roms/missile.ch8" },
  { name: "pong", path: "/roms/pong.ch8" },
  { name: "pong2", path: "/roms/pong2.ch8" },
  { name: "puzzle", path: "/roms/puzzle.ch8" },
  { name: "squash", path: "/roms/squash.ch8" },
  { name: "syzygy", path: "/roms/syzygy.ch8" },
  { name: "tank", path: "/roms/tank.ch8" },
  { name: "tetris", path: "/roms/tetris.ch8" },
  { name: "tictac", path: "/roms/tictac.ch8" },
  { name: "ufo", path: "/roms/ufo.ch8" },
  { name: "vbrix", path: "/roms/vbrix.ch8" },
  { name: "vers", path: "/roms/vers.ch8" },
  { name: "wall", path: "/roms/wall.ch8" },
  { name: "wipeoff", path: "/roms/wipeoff.ch8" },

]

const selectedRomIndex = ref<number>(0);

const workerPeer = inject<WorkerPeer>("workerPeer");

function handleLoadROMButton() {
  console.log('selectedRomIndex:', selectedRomIndex.value);
  const rom = roms[selectedRomIndex.value];
  workerPeer.loadROM(rom.path);
}

function handleStartButton() {
  workerPeer.startEmu();
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
</script>

<template>
  <main>
    <div class="control-panel">
      <div class="rom-loader">

        <select v-model="selectedRomIndex">
          <option v-for="(rom, index) in roms" :key="rom.path" :value="index">{{ rom.name }}</option>
        </select>

        <button @click="handleLoadROMButton">
          Load ROM
        </button>
      </div>

      <button @click="handleStartButton">
        Start
      </button>
    </div>

    <div class="game-area" tabindex="0" @keydown="handleKeyDown" @keyup="handleKeyUp">
      <DrawArea />
    </div>
  </main>
</template>

<style>
main {
  width: 100%;
  height: 100%;
  position: relative;
  z-index: 999999999;
}

.game-area {
  width: 100%;
  height: 100%;
}
</style>
