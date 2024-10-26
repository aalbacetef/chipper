

export type ROMEntry = {
  name: string;
  path: string;
};

const manifestURL = 'roms/manifest.json';

export function loadROMManifest(): Promise<ROMEntry[]> {
  return fetch(manifestURL)
    .then((res) => res.json())
    .then((data: Record<string, string>) => {
      const roms: ROMEntry[] = [];

      Object.keys(data).forEach((name) => {
        const entry = { name, path: data[name] };
        roms.push(entry);
      });

      return roms;
    });
}

