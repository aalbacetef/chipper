<script setup lang="ts">
import { RouterLink, RouterView } from 'vue-router';
import { useNotificationStore } from '@/stores/notifications';
import NotificationItem from '@/components/notification-item.vue';

const notifications = useNotificationStore();
</script>

<template>
  <header>
    <nav>
      <RouterLink to="/">Home</RouterLink>
    </nav>

    <div class="banner">
      <div class="title">chipper</div>
      <div class="sub-title">A CHIP-8 emulator</div>
    </div>
  </header>

  <RouterView />

  <div class="notifications">
    <notification-item
      v-for="item in notifications.items"
      @close="notifications.close(item.id)"
      :key="item.id"
      :state="item.status"
      :text="item.text"
    ></notification-item>
  </div>
</template>

<style scoped>
.title {
  font-size: 2.4em;
}

.sub-title {
  font-size: 1.6em;
}

nav {
  display: none;
}

.notifications {
  position: absolute;
  z-index: 9999;
  padding: 15px;
  right: 15px;
  bottom: 15px;
  color: white;
}
</style>
