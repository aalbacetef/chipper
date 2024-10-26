import { describe, it } from 'vitest';
import { hexToRGBA, RGBAToHex, type Color } from './color';

describe('color', () => {
  it('should convert hex appropriately', ({ expect }) => {
    const rgba = [0, 0, 0, 255];
    const hex = '#000000';
    const got = hexToRGBA(hex);

    got.forEach((v, k) => expect(v).toBe(rgba[k]));
  });

  it('should convert rgba to hex', ({ expect }) => {
    const rgba: Color = [0, 0, 0, 255];
    const hex = '#000000';

    const got = RGBAToHex(rgba);
    expect(got).toBe(hex);
  });
});
