<script setup lang="ts">
import { onMounted, ref, useTemplateRef } from 'vue';
import { run } from '../lib/game';

const loading = ref<boolean>(true);
const canvas = useTemplateRef<HTMLCanvasElement>("canvas");

onMounted(() => {
  loading.value = false;
  window.LoadROM();
  window.StartEmu();
  const ctx = canvas.value.getContext('2d');
  run(ctx);
});

</script>

<template>
  <div class="draw-area">
    <canvas class="debug" ref="canvas" width="640" height="320">Loading...</canvas>
  </div>
</template>

<style></style>
