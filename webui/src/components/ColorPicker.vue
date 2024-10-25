<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { type Color, defaultColors, RGBAToHex } from '@/lib/game';

type Props = {
  name: string;
  display: string;
};

const props = defineProps<Props>();
const emit = defineEmits<{
  update: [name: string, colorHex: string];
}>();

const chosenColor = ref<Color>('#000000');

onMounted(() => {
  let v = chosenColor.value;
  switch (props.name) {
    case 'set':
      v = RGBAToHex(defaultColors.set);
      break;
    case 'clear':
      v = RGBAToHex(defaultColors.clear);
      break;
  }
  chosenColor.value = v;
});

function handleInput(event) {
  emit('update', [props.name, event.target.value]);
}
</script>

<template>
  <div class="color-picker">
    <span class="name">{{ props.display }}: </span>
    <input type="color" :value="chosenColor" @input="handleInput" />
  </div>
</template>

<style scoped>
.name {
  display: inline-block;
  width: 100px;
}
</style>
