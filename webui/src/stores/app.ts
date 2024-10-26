import { inject, ref } from 'vue';
import { defineStore } from 'pinia';
import { type Color, type ColorNames, type ColorOptions } from '@/lib/game';
import { defaultColors } from '@/lib/game';
import { useNotificationStore } from '@/stores/notifications';
import { Status } from '@/lib/status';
import { WorkerPeer } from '@/lib/peer';


export type Buttons =
  | "start"
  | "stop"
  | "restart";

const notifications = useNotificationStore()

const workerPeer = inject<WorkerPeer>('workerPeer');
if (typeof workerPeer === 'undefined') {
  throw new Error("could not load workerPeer");
}

export const useAppStore = defineStore('app', () => {
  const loadedROM = ref<string>('');
  const colors = ref<ColorOptions>(defaultColors);

  function isROMLoaded(): boolean {
    return loadedROM.value === '';
  }

  function setColor(which: ColorNames, color: Color): void {
    switch (which) {
      case 'set':
        colors.value.set = color;
        break;
      case 'clear':
        colors.value.clear = color;
        break;
      default:
        notifications.push(Status.Error, `invalid prop: '${which}'`);
        return;
    }

    workerPeer!.setColors(colors.value);
  }

  function buttonClicked(which: Buttons) {
    switch (which) {
      case 'start':
        workerPeer!.startEmu();
        break;
      case 'stop':
        workerPeer!.stopEmu();
        break;
      case 'restart':
        workerPeer!.restartEmu();
        break;
      default:
        console.log(`unknown button clicked: '${which}'`);
        return;
    }
  }

  return {
    loadedROM,
    colors,


    buttonClicked,
    isROMLoaded,
    setColor,
  }
});
