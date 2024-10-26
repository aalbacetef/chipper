import { NewNotification, type Notification } from '@/lib/notification';
import type { Status } from "@/lib/status";
import { uuid } from '@/lib/uuid';
import { defineStore } from "pinia";
import { ref } from "vue";

export const useNotificationStore = defineStore('notifications', () => {
  const items = ref<Notification[]>([]);

  function close(id: string): void {
    items.value = items.value.filter(notification => notification.id !== id);
  }

  function push(status: Status, text: string, duration?: number): void {
    const id = uuid()
    const notification = NewNotification(id, status, text, duration);
    items.value.push(notification);

    setTimeout(() => close(id), duration);
  }

  return { items, close, push };
});
