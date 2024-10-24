


type CanvasContext = CanvasRenderingContext2D | OffscreenCanvasRenderingContext2D;

export type Color = [number, number, number, number];
export type Dims = [number, number];
export type ColorOptions = { set: Color, clear: Color }

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

export const KeyMap: Record<string, number> = {
  "KeyV": 0,
  "Digit1": 1,
  "Digit2": 2,
  "Digit3": 3,
  "Digit4": 4,
  "KeyQ": 5,
  "KeyW": 6,
  "KeyE": 7,
  "KeyR": 8,
  "KeyA": 9,
  "KeyS": 10,
  "KeyD": 11,
  "KeyF": 12,
  "KeyZ": 13,
  "KeyX": 14,
  "KeyC": 15,
};


export function mapKeyToHex(s: string): number {
  if (typeof KeyMap[s] === 'undefined') {
    throw new MissingKeyError(s);
  }

  return KeyMap[s];
}

class MissingKeyError extends Error {
  constructor(key: string) {
    super(`missing map for key: '${key}`);
  }
}
