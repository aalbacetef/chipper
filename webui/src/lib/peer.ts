import {
  MessageType,
  Event,
  type StopEmu,
  type RestartEmu,
  type SetColors,
  type SetTickPeriod,
} from '@/lib/messages';

import {
  KeyDirection,
  type GenericMessage,
  type LoadROM,
  type LoadWASM,
  type TransferOffscreenCanvas,
  type WorkerEvent,
  type StartEmu,
  type KeyEvent,
} from '@/lib/messages';
import { MissingKeyError, mapHexToKey, type ColorOptions, type KeyList } from './game';

import { useAppStore } from '@/stores/app';

const RunOnce = true;

type StateChangeCB = (state: Event) => void;
type Store = {
  setKeyState: (key: KeyList, dir: KeyDirection) => void;
};

// WorkerPeer provides a set of methods to interact with the Worker from the main client code.
export class WorkerPeer {
  store: Store | null = null;
  worker: Worker;
  callbacks: {
    [key in Event]?: StateChangeCB[];
  } = {};
  oneoffs: {
    [key in Event]?: StateChangeCB[];
  } = {};

  offscreenCanvas?: OffscreenCanvas;
  onscreenCanvas?: HTMLCanvasElement;

  constructor(worker: Worker) {
    this.worker = worker;
    this.worker.addEventListener('message', (msg) => this.handleMessage(msg.data));
  }

  setStore() {
    this.store = useAppStore();
  }

  loadWASM(filename: string): void {
    this.postMessage<LoadWASM>({
      type: MessageType.LoadWASM,
      data: { filename },
    });
  }

  loadROM(filename: string): void {
    this.on(Event.ROMLoaded, () => console.log('rom loaded'), RunOnce);

    this.postMessage<LoadROM>({
      type: MessageType.LoadROM,
      data: {
        filename, // @TODO: implement support for this
      },
    });
  }

  startEmu(): void {
    this.on(Event.EmuStarted, () => console.log('emu started'), RunOnce);

    this.postMessage<StartEmu>({
      type: MessageType.StartEmu,
      data: {},
    });
  }

  stopEmu(): void {
    this.on(Event.EmuStopped, () => console.log('emu stopped'), RunOnce);

    this.postMessage<StopEmu>({
      type: MessageType.StopEmu,
      data: {},
    });
  }

  restartEmu(): void {
    this.on(Event.EmuRestarted, () => console.log('emu restarted'), RunOnce);

    this.postMessage<RestartEmu>({
      type: MessageType.RestartEmu,
      data: {},
    });
  }

  setColors(colors: ColorOptions): void {
    this.on(Event.SetColors, () => console.log('updated colors'), RunOnce);

    this.postMessage<SetColors>({
      type: MessageType.SetColors,
      data: colors,
    });
  }

  setTickPeriod(period: number): void {
    this.on(Event.SetTickPeriod, () => console.log(`set tick period to: ${period}ms`), RunOnce);

    this.postMessage<SetTickPeriod>({
      type: MessageType.SetTickPeriod,
      data: period,
    });
  }

  setOnscreenCanvas(canvas: HTMLCanvasElement): void {
    this.onscreenCanvas = canvas;
  }

  makeOffscreenCanvas(): void {
    this.offscreenCanvas = this.onscreenCanvas!.transferControlToOffscreen();

    const msg: TransferOffscreenCanvas = {
      type: MessageType.TransferOffscreenCanvas,
      data: {
        canvas: this.offscreenCanvas,
      },
    };

    this.worker.postMessage(msg, [this.offscreenCanvas]);
  }

  sendKeyUp(key: number, repeat: boolean): void {
    this.sendKeyEvent(KeyDirection.Up, repeat, key);
  }

  sendKeyDown(key: number, repeat: boolean): void {
    this.sendKeyEvent(KeyDirection.Down, repeat, key);
  }

  sendKeyEvent(direction: KeyDirection, repeat: boolean, key: number): void {
    if (this.store !== null) {
      try {
        this.store.setKeyState(mapHexToKey(key), direction);
      } catch (err) {
        if (!(err instanceof MissingKeyError)) {
          console.error('[WorkerPeer.sendKeyEvent] error calling store.setKeyState: ', err);
        }
      }
    }

    this.postMessage<KeyEvent>({
      type: MessageType.KeyEvent,
      data: {
        key,
        repeat,
        direction,
      },
    });
  }

  postMessage<T>(msg: T): void {
    this.worker.postMessage(msg);
  }

  handleWorkerStateChange(msg: WorkerEvent): void {
    console.log('[handleWorkerStateChange] msg:', msg);
    this.runCallbacks(msg.data.state);
  }

  runCallbacks(state: Event) {
    if (typeof this.callbacks[state] === 'undefined') {
      return;
    }

    this.callbacks[state]!.forEach((cb) => cb(state));
  }

  runOneOffs(state: Event) {
    if (typeof this.oneoffs[state] === 'undefined') {
      return;
    }

    this.oneoffs[state]!.forEach((cb) => cb(state));
    this.oneoffs[state] = [];
  }

  handleMessage(msg: GenericMessage): void {
    switch (msg.type) {
      case MessageType.WorkerEvent:
        return this.handleWorkerStateChange(msg as WorkerEvent);

      default:
        console.log('unknown msg type:', msg.type);
        console.log(msg);
    }
  }

  on(state: Event, cb: StateChangeCB, once: boolean = false) {
    let target = this.callbacks;
    if (once) {
      target = this.oneoffs;
    }

    if (typeof target[state] === 'undefined') {
      target[state] = [cb];
      return;
    }

    target[state]!.push(cb);
  }
}
