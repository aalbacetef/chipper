<script setup lang="ts">
import { RouterLink, RouterView } from 'vue-router';
</script>

<template>
  <header>
    <nav>
      <RouterLink to="/">Home</RouterLink>
    </nav>

    <div class="banner">
      <div class="title">chipper</div>
      <div class="sub-title">A CHIP-8 emulator</div>
    </div>
  </header>

  <RouterView />
</template>

<style scoped>
.title {
  font-size: 2.4em;
}

.sub-title {
  font-size: 1.6em;
}

nav {
  display: none;
}
</style>
