import { run } from '@/lib/game';
import type { GenericMessage, LoadROM, LoadWASM, StartEmu, TransferOffscreenCanvas, WorkerEvent } from '@/lib/messages';
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


  startRendering() {

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
    // @TODO: implement support for loading various ROMs
    self.LoadROM();

    notifyStateChange(Event.ROMLoaded);
  }

  handleStartEmu(msg: StartEmu): void {
    self.StartEmu();

    notifyStateChange(Event.EmuStarted);
  }

  handleTransferOffscreenCanvas(msg: TransferOffscreenCanvas): void {
    const canvas = msg.data.canvas;
    this.canvas = canvas;
    const ctx = canvas.getContext('2d');
    if (ctx === null) {
      return;
    }

    const [w, h] = [64, 32];
    run(ctx, w, h, self.GetDisplay);
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


