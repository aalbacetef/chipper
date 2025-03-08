<script setup lang="ts">
import { computed } from 'vue';
import { useAppStore } from '@/stores/app';
import { KeyDirection } from '@/lib/messages';
import { type KeyList } from '@/lib/game';

defineOptions({ name: 'keyboard-viewer' });

const store = useAppStore();

const keys = ['1', '2', '3', '4', 'Q', 'W', 'E', 'R', 'A', 'S', 'D', 'F', 'Z', 'X', 'C', 'V'];

function strToKey(s: string): string {
  switch (s) {
    case '1':
    case '2':
    case '3':
    case '4':
      return 'Digit' + s;
    default:
      return 'Key' + s;
  }
}

let rows = [keys.slice(0, 4), keys.slice(4, 8), keys.slice(8, 12), keys.slice(12, 16)];

const pressed = computed(() => {
  let down = [];

  for (const key in store.keyStates) {
    if (store.keyStates[key as KeyList] === KeyDirection.Down) {
      down.push(key);
    }
  }

  return down;
});

function keyIsPressed(key: string): boolean {
  return pressed.value.indexOf(strToKey(key)) > -1;
}
</script>

<template>
  <div class="keyboard-viewer">
    <div class="rows">
      <div class="row" v-for="(row, index) in rows" :key="index">
        <div class="key" v-for="key in row" :key="key" :class="{ pressed: keyIsPressed(key) }">
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

.key.pressed {
  background-color: yellow;
}
</style>
