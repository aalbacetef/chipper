import type { ColorOptions } from './color';

export type Dims = [number, number];

type CanvasContext = CanvasRenderingContext2D | OffscreenCanvasRenderingContext2D;

// @TODO: flexible scale factor for different kinds of screens (e.g: mobile).
export function render(
  buf: Uint8Array,
  ctx: CanvasContext,
  dims: Dims,
  colors: ColorOptions
): void {
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
    resizeQuality: 'pixelated',
  };

  createImageBitmap(imgData, opts).then((bitmap) => {
    ctx.drawImage(bitmap, 0, 0);
  });
}

const clearValue = 0;

// drawImage will take the display data and generate an ImageData to draw on the canvas.
function drawImage(
  ctx: CanvasContext,
  buf: Uint8Array,
  dims: Dims,
  colors: ColorOptions
): ImageData {
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
