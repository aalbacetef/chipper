

export type GenericMessage = {
  type: MessageType;
  data: any;
};

export enum Event {
  WorkerLoaded = "worker-loaded",
  WASMLoaded = "wasm-loaded",
  ROMLoaded = "rom-loaded",
  EmuStarted = "emu-started",
  CanvasCreated = "canvas-created",
}

export enum MessageType {
  WorkerEvent = "worker-event",
  LoadWASM = "load-wasm",
  LoadROM = "load-rom",
  StartEmu = "start-emu",
  TransferOffscreenCanvas = "transfer-offscreen-canvas",
  KeyEvent = "key-event",
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

export enum KeyDirection {
  Up,
  Down,
}

export type KeyEvent = {
  type: MessageType.KeyEvent;
  data: {
    repeat: boolean;
    key: number;
    direction: KeyDirection;
  }
}
