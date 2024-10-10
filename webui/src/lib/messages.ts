

export type GenericMessage = {
  type: MessageType;
  data: any;
};

export enum Event {
  WorkerLoaded,
  WASMLoaded,
  ROMLoaded,
  EmuStarted,
  CanvasCreated,
}

export enum MessageType {
  WorkerEvent = "worker-event",
  LoadWASM = "load-wasm",
  LoadROM = "load-rom",
  StartEmu = "start-emu",
  TransferOffscreenCanvas = "transfer-offscreen-canvas",
}

export type WorkerEvent = {
  type: MessageType.WorkerEvent;
  data: {
    state: Event;
  }
}

export type LoadWASM = {
  type: MessageType.LoadWASM;
  data: {
    filename: string;
  };
}

export type LoadROM = {
  type: MessageType.LoadROM;
  data: {
    filename: string;
  };
}

export type StartEmu = {
  type: MessageType.StartEmu;
  data: {};
};

export type TransferOffscreenCanvas = {
  type: MessageType.TransferOffscreenCanvas;
  data: {
    canvas: OffscreenCanvas;
  };
}
