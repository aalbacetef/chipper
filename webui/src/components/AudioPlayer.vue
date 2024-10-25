<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import { type AudioManifest, type Song } from '@/lib/music';

type Props = {
  manifest: AudioManifest;
};

const props = defineProps<Props>();
const manifest = props.manifest;

const currentSongIndex = ref<number>(0);
const audioElem = ref<HTMLAudioElement>();

onMounted(() => {
  const album = manifest.album;
  if (album.songs.length === 0) {
    throw new Error('empty album');
  }
});

function playSong() {
  audioElem.value?.play();
}

function playNext() {
  pauseAudio();
  const n = manifest.album.songs.length;
  currentSongIndex.value = (currentSongIndex.value + 1) % n;
  audioElem.value!.load();
  playSong();
}

function playPrev() {
  let index = currentSongIndex.value - 1;
  const n = manifest.album.songs.length;
  if (index < 0) {
    index = n - 1;
  }

  currentSongIndex.value = index;

  audioElem.value!.load();
  playSong();
}

function pauseAudio() {
  audioElem.value?.pause();
}

const currentSong = computed<Song>(() => {
  const index = currentSongIndex.value;
  const n = manifest.album.songs.length;
  if (index >= n || index < 0) {
    return { path: '', position: 0, name: '' };
  }

  return manifest.album.songs[index];
});
</script>

<template>
  <div class="audio-player--wrapper">
    <audio autoplay ref="audioElem">
      <source :src="currentSong.path" type="audio/mp3" />
    </audio>
    <div class="album">
      <div class="panel">
        <div class="song-title">{{ currentSong.position }} - {{ currentSong.name }}</div>
        <div class="controls">
          <button @click="playPrev">⏮</button>
          <button @click="playSong">▶</button>
          <button @click="pauseAudio">⏸</button>
          <button @click="playNext">⏭</button>
        </div>
      </div>

      <div class="cover"><img :src="manifest.album.cover" /></div>
    </div>

    <div class="music-note">{{ manifest.album.attribution }}</div>
  </div>
</template>

<style scoped>
.audio-player--wrapper {
  margin-top: 5px;
  align-items: end;
  display: flex;
  flex-direction: column;
}

.album {
  display: flex;
  flex-direction: row;
  height: fit-content;
}

.panel {
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}

.song-title {
  padding-top: 5px;
  text-align: center;
}

.cover {
  display: flex;
}

.album img {
  width: 64px;
  height: auto;
  padding: 0;
  margin: 0;
}
</style>
