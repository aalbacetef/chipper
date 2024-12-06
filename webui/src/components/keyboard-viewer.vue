<script setup lang="ts">
import { KeyMap } from '@/lib/game';

defineOptions({ name: 'keyboard-viewer' });


console.log('keys:', Object.keys(KeyMap));

const keys = [
  '1', '2', '3', '4',
  'Q', 'W', 'E', 'R',
  'A', 'S', 'D', 'F',
  'Z', 'X', 'C', 'V',
];

function strToKey(s: string): string {
  switch (s) {
    case '1':
    case '2':
    case '3':
    case '4':
      return "Digit" + s;
    default:
      return "Key" + s;
  }
}

let rows = [
  keys.slice(0, 4),
  keys.slice(4, 8),
  keys.slice(8, 12),
  keys.slice(12, 16),
];
</script>

<template>
  <div class="keyboard-viewer">
    <div class="rows">
      <div class="row" v-for="(row, index) in rows" :key="index">
        <div class="key" v-for="key in row" :key="key">
          {{ key }}
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.keyboard-viewer {
  padding: 20px;
}

.rows {
  display: flex;
  flex-direction: column;
}

.row {
  display: flex;
  flex-direction: row;
}

.key {
  margin-right: 5px;
  padding: 5px 10px;
  outline: 1px solid var(--fg-color);
}
</style>
