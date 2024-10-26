<script setup lang="ts">
import { Status } from '@/lib/status';

defineOptions({ name: "notification-item" });

const props = defineProps<{
  text: string;
  state: Status;
}>();

const emit = defineEmits<{
  close: [];
}>();

function handleClose() {
  emit('close');
}
</script>

<template>
  <div class="notification--wrapper" :class="props.state">
    <div @click.prevent="handleClose" class="close">x</div>
    <p>{{ props.text }}</p>
  </div>
</template>

<style scoped>
.notification--wrapper {
  position: absolute;
  z-index: 9999;
  padding: 15px;
  right: 15px;
  bottom: 15px;
}

.notification--wrapper.success {
  background-color: #33FF33;
}

.notification--wrapper.error {
  background-color: #FF3333;
}

.notification--wrapper.warn {
  background-color: #AAAA11;
}
</style>
