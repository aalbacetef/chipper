export type WASMLoadResult = {
  go: Go;
  result: WebAssembly.WebAssemblyInstantiatedSource;
};

export async function loadWASM(wasmName: string): Promise<WASMLoadResult> {
  const go = new Go();

  const resp = await fetch(wasmName);
  const buffer = await resp.arrayBuffer();

  const result = await WebAssembly.instantiate(buffer, go.importObject);

  go.run(result.instance);

  return { go, result };
}
