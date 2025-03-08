
# chipper

Chipper is a CHIP-8 emulator written in Go. 

In this repo you'll also find a Vue app which embeds the emulator via WASM, allowing you to play CHIP-8 games, directly from the browser!

## Roadmap 

This is mostly for myself, just a set of TO-DOs for improving the code now that chipper hit MVP.

- [ ] flesh out README
- [ ] finish writing tests for the entire instruction set
- [x] improve the Web UI with keymap information
- [ ] improve the Web UI with keymap configuration
- [x] move state management to Pinia 
- [ ] filter out some ROMs that won't show by default in a production build
