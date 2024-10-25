<script setup lang="ts">
import { inject, onMounted, useTemplateRef } from 'vue';
import { workerPeer } from '@/lib/peer';

const canvas = useTemplateRef<HTMLCanvasElement>('canvas');
const workerPeer = inject<WorkerPeer>('workerPeer');

onMounted(() => {
  workerPeer.setOnscreenCanvas(canvas.value);
  workerPeer.makeOffscreenCanvas();
});
</script>

<template>
  <div class="draw-area">
    <canvas ref="canvas" width="640" height="320">Loading...</canvas>
  </div>
</template>

<style scoped>
canvas {
  background-color: var(--bg-color);
}

.draw-area {
  width: fit-content;
  border-radius: 5px;
  border: 5px solid var(--neon-purple);
}
</style>
