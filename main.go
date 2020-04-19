package main

import (
	"fmt"
	"os"

	"github.com/0xF0D0/c8emulator/chip8"
	"github.com/0xF0D0/c8emulator/renderer"
)

func main() {
	chip8Renderer := renderer.Initialize()
	emulator := chip8.Initialize()
	args := os.Args

	if len(args) < 2 {
		fmt.Println("Usage: c8emulator chip8Application")
		return
	}
	emulator.LoadGame(args[1])
	emulator.BindKeyboardDown(chip8Renderer.KeyboardDown)
	emulator.BindKeyboardUp(chip8Renderer.KeyboardUp)

	chip8Renderer.EmulateCycle = emulator.EmulateCycle
	chip8Renderer.BindRenderInput(emulator.GfxChannel())
	chip8Renderer.RunMainLoop()

}
