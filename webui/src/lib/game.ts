


type CanvasContext = CanvasRenderingContext2D | OffscreenCanvasRenderingContext2D;

export type Color = [number, number, number, number];
export type Dims = [number, number];
export type ColorOptions = { set: Color, clear: Color }

// @TODO: run should probably be passed to request animation frame
// @TODO: wrap get display calls in a buffer
// @TODO: flexible scale factor for different kinds of screens (e.g: mobile).
export function render(buf: Uint8Array, ctx: CanvasContext, dims: Dims, colors: ColorOptions): void {
  const [w, h] = dims;
  const n = w * h;
  const copied = GetDisplay(buf);

  if (copied !== n) {
    console.log(`failed to copy bytes, want ${n}, got ${copied}`);
    return;
  }

  const imgData = drawImage(ctx, buf, dims, colors);
  ctx.clearRect(0, 0, w * 10, h * 10);

  const opts: ImageBitmapOptions = {
    resizeWidth: w * 10,
    resizeHeight: h * 10,
    resizeQuality: "pixelated",
  };

  createImageBitmap(imgData, opts)
    .then(bitmap => {
      ctx.drawImage(bitmap, 0, 0);
    });
}


const clearValue = 0;


// drawImage will take the display data and generate an ImageData to draw on the canvas. 
function drawImage(ctx: CanvasContext, buf: Uint8Array, dims: Dims, colors: ColorOptions): ImageData {
  const [w, h] = dims;
  const imgData = ctx.createImageData(w, h);
  const data = imgData.data;
  const n = buf.length;

  for (let k = 0; k < n; k++) {
    const index = k * 4;

    let color = colors.set;
    if (buf[k] === clearValue) {
      color = colors.clear;
    }

    for (let j = 0; j < 4; j++) {
      data[index + j] = color[j];
    }
  }

  return imgData;
}

type KeyMap = {
  [key: string]: number;
}

function kkey(s: string): string {
  const digit = /[0-9]/
  const letter = /[a-zA-Z]/
  if (digit.test(s)) { return 'Digit' + s }
  if (letter.test(s)) { return 'Key' + s }

  return s;
}

function loadKeyMap(): KeyMap {
  const keys = [
    'V', '1', '2', '3', '4',
    'Q', 'W', 'E', 'R',
    'A', 'S', 'D', 'F',
    'Z', 'X', 'C',
  ].map(s => kkey(s));


  const o = {};
  keys.forEach((key, index) => o[key] = index);

  return o;
}

export function mapKeyToHex(s: string): number {
  const keyMap = loadKeyMap();
  if (typeof keyMap[s] === 'undefined') {
    throw new MissingKeyError(s);
  }

  return keyMap[s];
}

class MissingKeyError extends Error {
  constructor(key: string) {
    super(`missing map for key: '${key}`);
  }
}
