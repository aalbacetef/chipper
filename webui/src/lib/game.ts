


type CanvasContext = CanvasRenderingContext2D | OffscreenCanvasRenderingContext2D;

export type Dims = [number, number];
export type Color = [number, number, number, number];
export type ColorOptions = { set: Color, clear: Color }

export type ROMEntry = {
  name: string;
  path: string;
}

const manifestURL = 'roms/manifest.json';

export function loadROMManifest(): Promise<ROMEntry[]> {
  return fetch(manifestURL)
    .then(res => res.json())
    .then((data: Record<string, string>) => {
      const roms: ROMEntry[] = [];

      Object.keys(data).forEach(name => {
        const entry = { name, path: data[name] };
        roms.push(entry);
      });

      return roms;
    })
}

export const defaultColors: ColorOptions = {
  set: [10, 200, 10, 150],
  clear: [0, 0, 0, 255],
};

export function hexToRGBA(s: string): Color {
  console.log('hexToRGBA: ', s);

  const wantLen = 7;
  if (s.length !== wantLen) {
    throw new Error(`invalid hex value: ${s}`);
  }

  const r = s.substring(1, 3);
  const g = s.substring(3, 5);
  const b = s.substring(5);
  const a = 0xFF;

  return [
    parseInt(r, 16),
    parseInt(g, 16),
    parseInt(b, 16),
    a,
  ];
}

export function RGBAToHex(c: Color): string {
  let s = "#";
  s += pad(c[0].toString(16));
  s += pad(c[1].toString(16))
  s += pad(c[2].toString(16));

  return s;
}

function pad(s: string): string {
  if (s.length === 1) {
    s = '0' + s;
  }

  return s;
}

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
