<script setup lang="ts">
import { Status } from '@/lib/status';

defineOptions({ name: 'notification-item' });

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
  <div class="notification" :class="props.state">
    <div class="close" @click="handleClose">×</div>
    <p>{{ props.text }}</p>
  </div>
</template>

<style scoped>
.notification {
  padding: 10px 15px;
}

.notification.success {
  background-color: #33ff33;
  background-color: #3d8d3d;
}

.notification.error {
  background-color: #ff3333;
  background-color: #b94f4f;
}

.notification.warn {
  background-color: #aaaa11;
}

.close {
  text-align: right;
  cursor: pointer;
  height: fit-content;
  position: absolute;
  right: 30px;
  top: 18px;
}
</style>
