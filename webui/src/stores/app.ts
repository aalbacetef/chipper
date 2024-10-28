import { inject, ref } from 'vue';
import { defineStore } from 'pinia';
import { type Color, type ColorNames, type ColorOptions } from '@/lib/game';
import { defaultColors } from '@/lib/game';
import { useNotificationStore } from '@/stores/notifications';
import { Status } from '@/lib/status';
import { WorkerPeer } from '@/lib/peer';

export type Buttons = 'start' | 'stop' | 'restart';

export enum AudioState {
  Playing,
  NotPlaying,
  Paused,
}

export const defaultTickPeriod = 2 // 2 milliseconds 

export const useAppStore = defineStore('app', () => {
  // state
  const tickPeriodMilliseconds = ref<number>(defaultTickPeriod);
  const audioState = ref<AudioState>(AudioState.NotPlaying);
  const loadedROM = ref<string>('');
  const colors = ref<ColorOptions>(defaultColors);

  // helpers
  const notifications = useNotificationStore();
  const workerPeer = inject<WorkerPeer>('workerPeer');
  if (typeof workerPeer === 'undefined') {
    throw new Error('could not load workerPeer');
  }

  function isROMLoaded(): boolean {
    return loadedROM.value === '';
  }

  function audioIsPlaying(): boolean {
    return audioState.value === AudioState.Playing;
  }

  function setAudioState(s: AudioState): void {
    switch (s) {
      default:
        notifications.push(Status.Error, `unknown audio status: ${s}`);
        return;
      case AudioState.Playing, AudioState.Paused, AudioState.NotPlaying:
        audioState.value = s;
        return;
    }
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

  function buttonClicked(which: Buttons): void {
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

  function setTickPeriod(period: number): void {
    if (period <= 0) {
      throw new Error(`period must be > 0, got ${period}`);
    }

    tickPeriodMilliseconds.value = period;
    workerPeer!.setTickPeriod(period);
  }

  return {
    // state props 
    audioState,
    loadedROM,
    colors,
    tickPeriodMilliseconds,

    // methods
    audioIsPlaying,
    setAudioState,
    buttonClicked,
    isROMLoaded,
    setColor,
    setTickPeriod,
  };
});
