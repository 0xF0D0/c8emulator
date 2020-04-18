package main

import (
	"fmt"
	"time"

	"github.com/0xF0D0/c8emulator/renderer"
)

func main() {
	chip8Renderer := renderer.Initialize()

	go func() {
		for {
			select {
			case k := <-chip8Renderer.KeyboardDown:
				fmt.Println("key down", k)
			case k := <-chip8Renderer.KeyboardUp:
				fmt.Println("key up", k)
			}
		}
	}()

	go func() {
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
