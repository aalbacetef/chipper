<script setup lang="ts">
import { inject } from "vue";
import { WorkerPeer } from "@/lib/peer";
import { Event } from "@/lib/messages";
import DrawArea from "@/components/DrawArea.vue";

const workerPeer = inject<WorkerPeer>("workerPeer");
workerPeer.on(Event.EmuStarted, () => {
  workerPeer.syncCanvases();
});

function handleLoadROMButton() {
  workerPeer.loadROM();
}

function handleStartButton() {
  workerPeer.startEmu();
}

</script>

<template>
  <main>
    <div class="control-panel">
      <button @click="handleLoadROMButton">
        Load ROM
      </button>

      <button @click="handleStartButton">
        Start
      </button>
    </div>

    <DrawArea />
  </main>
</template>
