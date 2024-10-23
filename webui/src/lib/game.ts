


type CanvasContext = CanvasRenderingContext2D | OffscreenCanvasRenderingContext2D;

// @TODO: reuse buffer instead of declaring a new one.
// @TODO: run should probably be passed to request animation frame
// @TODO: wrap get display calls in a buffer
// @TODO: flexible scale factor for different kinds of screens (e.g: mobile).
export function render(ctx: CanvasContext, w: number, h: number): void {
  const n = w * h;
  const buf = new Uint8Array(n);
  const copied = GetDisplay(buf);

  if (copied !== n) {
    console.log(`failed to copy bytes, want ${n}, got ${copied}`);
    return;
  }

  const imgData = drawImage(ctx, buf, w, h);
  ctx.clearRect(0, 0, w * 10, h * 10);

  createImageBitmap(
    imgData,
    {
      resizeWidth: w * 10,
      resizeHeight: h * 10,
      resizeQuality: "pixelated",
    },
  )
    .then(bitmap => {
      ctx.drawImage(bitmap, 0, 0);
    });
}


const setColor = [10, 200, 10, 150]

function drawImage(ctx: CanvasContext, buf: Uint8Array, w: number, h: number): ImageData {
  const imgData = ctx.createImageData(w, h);
  const data = imgData.data;
  const n = buf.length;

  for (let k = 0; k < n; k++) {
    const index = k * 4;
    if (buf[k] === 0) {
      continue;
    }

    for (let j = 0; j < 4; j++) {
      data[index + j] = setColor[j];
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
