import { inject, ref } from 'vue';
import { defineStore } from 'pinia';
import { type Color, type ColorOptions } from '@/lib/game';
import { defaultColors } from '@/lib/game';
import { useNotificationStore } from '@/stores/notifications';
import { Status } from '@/lib/status';
import { WorkerPeer } from '@/lib/peer';


const notifications = useNotificationStore()

const workerPeer = inject<WorkerPeer>('workerPeer');

export const useAppStore = defineStore('app', () => {
  const loadedROM = ref<string>('');
  const colors = ref<ColorOptions>(defaultColors);

  function isROMLoaded(): boolean {
    return loadedROM.value === '';
  }

  function setColor(which: 'set' | 'clear', color: Color): void {
    switch (which) {
      case 'set':
        colors.value.set = color;
        return;
      case 'clear':
        colors.value.clear = color;
        return;
      default:
        notifications.push(Status.Error, `invalid prop: '${which}'`);
        return;
    }
  }

  return {
    loadedROM,
    colors,


    isROMLoaded,
    setColor,
  }
});
