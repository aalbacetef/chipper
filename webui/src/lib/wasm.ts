
export type WASMLoadResult = {
  go: Go;
  result: WebAssembly.WebAssemblyInstantiatedSource;
}

export async function loadWASM(wasmName: string): Promise<WASMLoadResult> {
  const go = new Go();

  const result = await WebAssembly.instantiateStreaming(
    fetch(wasmName, { headers: { "Content-Type": "application/wasm" } }),
    go.importObject,
  )

  go.run(result.instance);


  return { go, result };
}
