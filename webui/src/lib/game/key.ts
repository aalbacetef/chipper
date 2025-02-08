export const KeyMap: Record<KeyList, number> = {
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

export type KeyList = "KeyV" | "Digit1" | "Digit2" | "Digit3" | "Digit4" | "KeyQ" | "KeyW" | "KeyE" | "KeyR" | "KeyA" | "KeyS" | "KeyD" | "KeyF" | "KeyZ" | "KeyX" | "KeyC";

export function mapKeyToHex(s: KeyList): number {
  if (typeof KeyMap[s] === 'undefined') {
    throw new MissingKeyError(s);
  }

  return KeyMap[s];
}

export function mapHexToKey(num: number): KeyList {
  for (const key in KeyMap) {
    const val = KeyMap[key as KeyList];
    if (num === val) {
      return key as KeyList;
    }
  }

  throw new MissingKeyError(`num: ${num}`);
}

export class MissingKeyError extends Error {
  constructor(key: string) {
    super(`missing map for key: '${key}`);
  }
}
