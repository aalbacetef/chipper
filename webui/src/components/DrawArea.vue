<script setup lang="ts">
import { inject, onMounted, ref, useTemplateRef } from 'vue';

const loading = ref<boolean>(true);
const canvas = useTemplateRef<HTMLCanvasElement>("canvas");


onMounted(() => {
  loading.value = false;

  const workerPeer = inject<WorkerPeer>("workerPeer");
  workerPeer.setOnscreenCanvas(canvas.value);
  workerPeer.makeOffscreenCanvas();
});

</script>

<template>
  <div class="draw-area">
    <canvas class="debug" ref="canvas" width="640" height="320">Loading...</canvas>
  </div>
</template>

<style></style>
