import type { KeyDirection } from "./lib/messages";

declare declare global {
  class Go {
    constructor();
    run(instance: WebAssembly.Instance): void;

    exited: boolean;
    mem: DataView;
    importObject: WebAssembly.Imports;
  }

  function StartEmu(): void;
  function LoadROM(arr: Uint8Array, n: number): void;
  function GetDisplay(buf: Uint8Array): number;
  function SendKeyboardEvent(key: number, repeat: boolean, direction: KeyDirection): void;
}

