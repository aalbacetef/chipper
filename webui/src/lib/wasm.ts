
type WASMLoadResult = {
  go: Go;
  result: WebAssembly.WebAssemblyInstantiatedSource;
}

export async function initialize(wasmName: string): Promise<WASMLoadResult> {
  const go = new Go();

  console.log('initialize');
  const result = await WebAssembly.instantiateStreaming(
    fetch(wasmName, { headers: { "Content-Type": "application/wasm" } }),
    go.importObject,
  )
  console.log('result:', result);

  go.run(result.instance);


  return { go, result };
}
