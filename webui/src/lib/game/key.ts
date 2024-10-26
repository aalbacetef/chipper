export const KeyMap: Record<string, number> = {
  KeyV: 0,
  Digit1: 1,
  Digit2: 2,
  Digit3: 3,
  Digit4: 4,
  KeyQ: 5,
  KeyW: 6,
  KeyE: 7,
  KeyR: 8,
  KeyA: 9,
  KeyS: 10,
  KeyD: 11,
  KeyF: 12,
  KeyZ: 13,
  KeyX: 14,
  KeyC: 15,
};

export function mapKeyToHex(s: string): number {
  if (typeof KeyMap[s] === 'undefined') {
    throw new MissingKeyError(s);
  }

  return KeyMap[s];
}

export class MissingKeyError extends Error {
  constructor(key: string) {
    super(`missing map for key: '${key}`);
  }
}
