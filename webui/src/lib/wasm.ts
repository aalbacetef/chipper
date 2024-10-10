
export type WASMLoadResult = {
  go: Go;
  result: WebAssembly.WebAssemblyInstantiatedSource;
}

export async function loadWASM(wasmName: string): Promise<WASMLoadResult> {
  console.log('wasmName:', wasmName);
  const go = new Go();

  const result = await WebAssembly.instantiateStreaming(
    fetch(wasmName, { headers: { "Content-Type": "application/wasm" } }),
    go.importObject,
  )
  console.log('result:', result);

  go.run(result.instance);


  return { go, result };
}
