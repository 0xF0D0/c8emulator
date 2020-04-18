package main

import (
	"time"

	"github.com/0xF0D0/c8emulator/renderer"
)

func main() {
	chip8Renderer := renderer.Initialize()

	go func() {
		//time.Sleep(2*time.Second)
		b := make([]byte, 64*32)
		for i := 0; i < 64*32; i++ {
			b[i] = byte(i % 2)
		}
		ch := make(chan []byte)
		chip8Renderer.BindRenderInput(ch)
		time.Sleep(2 * time.Second)
		ch <- b
	}()

	chip8Renderer.RunMainLoop()

}
