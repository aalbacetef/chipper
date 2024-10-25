import { defaultColors, render, type ColorOptions, type Dims } from '@/lib/game';
import type { GenericMessage, KeyEvent, LoadROM, LoadWASM, RestartEmu, SetColors, StartEmu, StopEmu, TransferOffscreenCanvas, WorkerEvent } from '@/lib/messages';
import { MessageType, Event } from "@/lib/messages";
import { loadWASM, type WASMLoadResult } from '@/lib/wasm';

// ensure wasm glue code for Go is registered in the worker.
import "@/worker/wasm_exec.js";

export function initialize() {
  registerHandlers();
  notifyStateChange(Event.WorkerLoaded);
}

const emulatorConstants = {
  w: 64,
  h: 32,
}


// WorkerInstance implements the core functionality of our Web Worker 
// in which we run our WASM emulator.
class WorkerInstance {
  result?: WASMLoadResult;
  canvas?: OffscreenCanvas;
  colors: ColorOptions = defaultColors;


  startRendering(): void {
    console.log('startRendering');

    const ctx = this.canvas!.getContext('2d');
    if (ctx === null) {
      console.log('could not acquire context');
      return;
    }

    const { w, h } = emulatorConstants;
    const buf: Uint8Array = new Uint8Array(w * h);

    this.loop(buf, ctx, w, h);
  }

  loop(buf: Uint8Array, ctx: OffscreenCanvasRenderingContext2D, w: number, h: number) {
    const colors = this.colors;
    const dims: Dims = [w, h];

    render(buf, ctx, dims, colors);
    requestAnimationFrame(() => this.loop(buf, ctx, w, h));
  }

  handleMessage(msg: GenericMessage) {
    switch (msg.type) {
      case MessageType.LoadWASM:
        return this.handleLoadWASM(msg as LoadWASM);

      case MessageType.LoadROM:
        return this.handleLoadROM(msg as LoadROM);

      case MessageType.StartEmu:
        return this.handleStartEmu(msg as StartEmu);

      case MessageType.StopEmu:
        return this.handleStopEmu(msg as StopEmu);

      case MessageType.RestartEmu:
        return this.handleRestartEmu(msg as RestartEmu);

      case MessageType.TransferOffscreenCanvas:
        return this.handleTransferOffscreenCanvas(msg as TransferOffscreenCanvas);

      case MessageType.KeyEvent:
        return this.handleKeyEvent(msg as KeyEvent);

      case MessageType.SetColors:
        return this.handleSetColors(msg as SetColors);

      default:
        console.log('unhandled message: ', msg);
    }
  }

  handleLoadWASM(msg: LoadWASM): void {
    loadWASM(msg.data.filename)
      .then(result => {
        this.result = result;
        notifyStateChange(Event.WASMLoaded);
      });
  }

  handleLoadROM(msg: LoadROM): void {
    const filename = msg.data.filename;
    fetch(filename)
      .then(r => {
        return r.arrayBuffer()
      }).then(buf => {
        const arr = new Uint8Array(buf);

        LoadROM(arr, arr.byteLength);
        notifyStateChange(Event.ROMLoaded);

      });
  }

  handleStartEmu(msg: StartEmu): void {
    self.StartEmu();
    notifyStateChange(Event.EmuStarted);
    this.startRendering();
  }

  handleStopEmu(msg: StopEmu): void {
    self.StopEmu();
    notifyStateChange(Event.EmuStopped);
  }

  handleRestartEmu(msg: RestartEmu): void {
    self.RestartEmu();
    notifyStateChange(Event.EmuRestarted);
  }

  handleSetColors(msg: SetColors): void {
    this.colors = msg.data;
    notifyStateChange(Event.SetColors);
  }


  handleTransferOffscreenCanvas(msg: TransferOffscreenCanvas): void {
    const canvas = msg.data.canvas;
    this.canvas = canvas;
  }

  handleKeyEvent(msg: KeyEvent): void {
    const { key, repeat, direction } = msg.data;

    SendKeyboardEvent(key, repeat, direction);
  }
}

function registerHandlers() {
  console.log("registering handlers...");
  const worker = new WorkerInstance();

  self.addEventListener(
    "message",
    (event) => {
      worker.handleMessage(event.data as GenericMessage);
    },
  );

  notifyStateChange(Event.WorkerLoaded, 100);
}

function notifyStateChange(state: Event, delay?: number): void {
  const message: WorkerEvent = {
    type: MessageType.WorkerEvent,
    data: {
      state,
    }
  };

  if (typeof delay !== 'undefined') {
    setTimeout(() => self.postMessage(message), delay);
    return;
  }

  self.postMessage(message);
}


