
export type Color = [number, number, number, number];
export type ColorNames = 'set' | 'clear';
export type ColorOptions = { set: Color; clear: Color };

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
  const a = 0xff;

  return [parseInt(r, 16), parseInt(g, 16), parseInt(b, 16), a];
}

export function RGBAToHex(c: Color): string {
  let s = '#';
  s += pad(c[0].toString(16));
  s += pad(c[1].toString(16));
  s += pad(c[2].toString(16));

  return s;
}

function pad(s: string): string {
  if (s.length === 1) {
    s = '0' + s;
  }

  return s;
}

