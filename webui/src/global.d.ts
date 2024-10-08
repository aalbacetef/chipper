
declare class Go {
  constructor();
  run(instance: WebAssembly.Instance): void;

  exited: boolean;
  mem: DataView;
  importObject: WebAssembly.Imports;
}

declare function StartEmu(): void;
declare function LoadROM(): void;
declare function GetDisplay(buf: Uint8Array): int;

