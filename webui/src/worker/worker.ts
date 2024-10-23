import { render } from '@/lib/game';
import type { GenericMessage, KeyEvent, LoadROM, LoadWASM, StartEmu, TransferOffscreenCanvas, WorkerEvent } from '@/lib/messages';
import { MessageType, Event } from "@/lib/messages";
import { loadWASM, type WASMLoadResult } from '@/lib/wasm';

import * as WASMExec from "@/worker/wasm_exec.js";

export function initialize() {
  registerHandlers();
  notifyStateChange(Event.WorkerLoaded);
}

class WorkerInstance {
  result?: WASMLoadResult;
  canvas?: OffscreenCanvas;


  startRendering(): void {
    console.log('startRendering');
    const ctx = this.canvas!.getContext('2d');
    if (ctx === null) {
      return;
    }
    const [w, h] = [64, 32];
    this.loop(ctx, w, h);
  }

  loop(ctx: OffscreenCanvasRenderingContext2D, w: number, h: number) {
    render(ctx, w, h);
    setTimeout(() => this.loop(ctx, w, h), 16);
  }

  handleMessage(msg: GenericMessage) {
    switch (msg.type) {
      case MessageType.LoadWASM:
        return this.handleLoadWASM(msg as LoadWASM);

      case MessageType.LoadROM:
        return this.handleLoadROM(msg as LoadROM);

      case MessageType.StartEmu:
        return this.handleStartEmu(msg as StartEmu);

      case MessageType.TransferOffscreenCanvas:
        return this.handleTransferOffscreenCanvas(msg as TransferOffscreenCanvas);

      case MessageType.KeyEvent:
        return this.handleKeyEvent(msg as KeyEvent);

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
  console.log("registerHandlers:");
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


