
// const w = 64;
// const h = 32;
const delay = 180;


type CanvasContext = CanvasRenderingContext2D | OffscreenCanvasRenderingContext2D;

// @TODO: reuse buffer instead of declaring a new one.
// @TODO: run should probably be passed to request animation frame
// @TODO: wrap get display calls in a buffer
export function run(ctx: CanvasContext, w: number, h: number, GetDisplay): void {
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
      setTimeout(() => run(ctx, w, h, GetDisplay), delay);
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
