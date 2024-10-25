export type Song = {
  position: number;
  name: string;
  path: string;
};

export type Album = {
  attribution: string;
  artist: string;
  name: string;
  songs: Song[];
  cover: string;
};

export type AudioManifest = {
  album: Album;
};

const manifestURLPath = 'audio/manifest.json';

export function loadAudioManifest(): Promise<AudioManifest> {
  return fetch(manifestURLPath)
    .then((res) => res.json())
    .catch((err) => console.error('could not load manifest: ', err));
}
