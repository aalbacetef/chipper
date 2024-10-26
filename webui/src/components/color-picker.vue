<script setup lang="ts">
import { ref } from 'vue';
import { defaultColors, RGBAToHex } from '@/lib/game';

defineOptions({ name: 'color-picker' });

type Props = {
  name: string;
  display: string;
};

const props = defineProps<Props>();
const emit = defineEmits<{
  update: [name: string, colorHex: string];
}>();

const chosenColor = ref<string>('#000000');

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

function handleInput(event: Event): void {
  const target = event.target as HTMLInputElement;
  emit('update', props.name, target.value);
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
