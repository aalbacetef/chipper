<script setup lang="ts">
import { ref } from 'vue';
import { defaultColors, hexToRGBA, RGBAToHex } from '@/lib/game';
import type { ColorNames } from '@/lib/game';
import { useAppStore } from '@/stores/app';

defineOptions({ name: 'color-picker' });

const appStore = useAppStore();

type Props = {
  name: ColorNames;
  display: string;
};

const props = defineProps<Props>();

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
  appStore.setColor(props.name, hexToRGBA(target.value));
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
