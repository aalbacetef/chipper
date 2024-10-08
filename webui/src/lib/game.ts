
const w = 64;
const h = 32;
const delay = 180;

export function run(ctx: CanvasRenderingContext2D): void {
  const n = w * h;
  const buf = new Uint8Array(n);
  const copied = window.GetDisplay(buf);
  // ctx.setTransform(1, 0, 0, 1, 0, 0);

  if (copied !== n) {
    console.log(`failed to copy bytes, want ${n}, got ${copied}`);
    return;
  }

  const imgData = drawImage(ctx, buf);
  ctx.clearRect(0, 0, w * 10, h * 10);
  createImageBitmap(imgData, { resizeWidth: w * 10, resizeHeight: h * 10, resizeQuality: "pixelated" })
    .then(bitmap => {
      ctx.drawImage(bitmap, 0, 0);
      setTimeout(() => run(ctx), delay);
    });
}


const setColor = [10, 200, 10, 150]

function drawImage(ctx: CanvasRenderingContext2D, buf: Uint8Array): ImageData {
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
